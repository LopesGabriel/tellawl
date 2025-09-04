package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/ports"
)

type PostgreSQLWalletRepository struct {
	db        *sql.DB
	publisher ports.EventPublisher
}

func NewPostgreSQLWalletRepository(db *sql.DB, publisher ports.EventPublisher) *PostgreSQLWalletRepository {
	return &PostgreSQLWalletRepository{
		db:        db,
		publisher: publisher,
	}
}

func (r *PostgreSQLWalletRepository) FindById(ctx context.Context, id string) (*models.Wallet, error) {
	query := `SELECT id, creator_id, name, balance_value, balance_offset, created_at, updated_at 
			  FROM wallets WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var wallet models.Wallet
	var updatedAt sql.NullTime

	err := row.Scan(
		&wallet.Id,
		&wallet.CreatorId,
		&wallet.Name,
		&wallet.Balance.Value,
		&wallet.Balance.Offset,
		&wallet.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	if updatedAt.Valid {
		wallet.UpdatedAt = &updatedAt.Time
	}

	// Load users
	users, err := r.loadWalletUsers(id)
	if err != nil {
		return nil, err
	}
	wallet.Users = users

	// Load transactions
	transactions, err := r.loadWalletTransactions(id)
	if err != nil {
		return nil, err
	}
	wallet.Transactions = transactions

	return &wallet, nil
}

func (r *PostgreSQLWalletRepository) FindByUserId(ctx context.Context, userId string) ([]models.Wallet, error) {
	query := `SELECT DISTINCT w.id, w.creator_id, w.name, w.balance_value, w.balance_offset, w.created_at, w.updated_at 
			  FROM wallets w 
			  JOIN wallet_users wu ON w.id = wu.wallet_id 
			  WHERE wu.user_id = $1`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []models.Wallet

	for rows.Next() {
		var wallet models.Wallet
		var updatedAt sql.NullTime

		err := rows.Scan(
			&wallet.Id,
			&wallet.CreatorId,
			&wallet.Name,
			&wallet.Balance.Value,
			&wallet.Balance.Offset,
			&wallet.CreatedAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			wallet.UpdatedAt = &updatedAt.Time
		}

		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func (r *PostgreSQLWalletRepository) Save(ctx context.Context, wallet *models.Wallet) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save wallet
	query := `INSERT INTO wallets (id, creator_id, name, balance_value, balance_offset, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  ON CONFLICT (id) DO UPDATE SET
			  name = EXCLUDED.name,
			  balance_value = EXCLUDED.balance_value,
			  balance_offset = EXCLUDED.balance_offset,
			  updated_at = $8`

	now := time.Now()
	_, err = tx.Exec(query,
		wallet.Id,
		wallet.CreatorId,
		wallet.Name,
		wallet.Balance.Value,
		wallet.Balance.Offset,
		wallet.CreatedAt,
		wallet.UpdatedAt,
		now,
	)

	if err != nil {
		return err
	}

	// Save wallet users
	_, err = tx.Exec("DELETE FROM wallet_users WHERE wallet_id = $1", wallet.Id)
	if err != nil {
		return err
	}

	for _, user := range wallet.Users {
		_, err = tx.Exec("INSERT INTO wallet_users (wallet_id, user_id) VALUES ($1, $2)", wallet.Id, user.Id)
		if err != nil {
			return err
		}
	}

	// Save transactions
	for _, transaction := range wallet.Transactions {
		createdByJSON, _ := json.Marshal(transaction.CreatedBy)
		_, err = tx.Exec(`INSERT INTO transactions (id, wallet_id, amount_value, amount_offset, created_by, type, description, created_at)
						  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
						  ON CONFLICT (id) DO NOTHING`,
			transaction.Id,
			wallet.Id,
			transaction.Amount.Value,
			transaction.Amount.Offset,
			createdByJSON,
			string(transaction.Type),
			transaction.Description,
			transaction.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgreSQLWalletRepository) loadWalletUsers(walletId string) ([]models.User, error) {
	query := `SELECT u.id, u.first_name, u.last_name, u.email, u.hashed_password, u.created_at, u.updated_at
			  FROM users u
			  JOIN wallet_users wu ON u.id = wu.user_id
			  WHERE wu.wallet_id = $1`

	rows, err := r.db.Query(query, walletId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var updatedAt sql.NullTime

		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.HashedPassword,
			&user.CreatedAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			user.UpdatedAt = &updatedAt.Time
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *PostgreSQLWalletRepository) loadWalletTransactions(walletId string) ([]models.Transaction, error) {
	query := `SELECT id, amount_value, amount_offset, created_by, type, description, created_at
			  FROM transactions WHERE wallet_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, walletId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction
		var createdByJSON []byte

		err := rows.Scan(
			&transaction.Id,
			&transaction.Amount.Value,
			&transaction.Amount.Offset,
			&createdByJSON,
			&transaction.Type,
			&transaction.Description,
			&transaction.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(createdByJSON, &transaction.CreatedBy)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
