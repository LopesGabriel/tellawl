package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/lopesgabriel/tellawl/packages/tracing"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/ports"
	repohttp "github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PostgreSQLWalletRepository struct {
	db         *sql.DB
	publisher  ports.EventPublisher
	memberRepo *repohttp.HTTPMemberRepository
	tracer     trace.Tracer
}

func NewPostgreSQLWalletRepository(db *sql.DB, publisher ports.EventPublisher, memberRepo *repohttp.HTTPMemberRepository) *PostgreSQLWalletRepository {
	return &PostgreSQLWalletRepository{
		db:         db,
		publisher:  publisher,
		memberRepo: memberRepo,
		tracer:     tracing.GetTracer("github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database/postgresql/PostgreSQLWalletRepository"),
	}
}

func (r *PostgreSQLWalletRepository) FindById(ctx context.Context, id string) (*models.Wallet, error) {
	ctx, span := r.tracer.Start(ctx, "FindByID")
	defer span.End()

	query := `SELECT id, creator_id, name, balance_value, balance_offset, created_at, updated_at 
			  FROM wallets WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

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
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if updatedAt.Valid {
		wallet.UpdatedAt = &updatedAt.Time
	}

	// Load users
	members, err := r.loadWalletMembers(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	wallet.Members = members

	// Load transactions
	transactions, err := r.loadWalletTransactions(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	wallet.Transactions = transactions

	span.SetStatus(codes.Ok, "Wallet found")
	return &wallet, nil
}

func (r *PostgreSQLWalletRepository) FindByUserId(ctx context.Context, userId string) ([]models.Wallet, error) {
	ctx, span := r.tracer.Start(ctx, "FindByUserId")
	defer span.End()

	query := `SELECT DISTINCT w.id, w.creator_id, w.name, w.balance_value, w.balance_offset, w.created_at, w.updated_at 
			  FROM wallets w 
			  JOIN wallet_users wu ON w.id = wu.wallet_id 
			  WHERE wu.member_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
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
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return nil, err
		}

		if updatedAt.Valid {
			wallet.UpdatedAt = &updatedAt.Time
		}

		wallets = append(wallets, wallet)
	}

	span.SetStatus(codes.Ok, "Wallets found")
	return wallets, nil
}

func (r *PostgreSQLWalletRepository) Save(ctx context.Context, wallet *models.Wallet) error {
	ctx, span := r.tracer.Start(ctx, "Save")
	defer span.End()

	tx, err := r.db.Begin()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
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
	_, err = tx.ExecContext(ctx, query,
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
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	// Sync wallet users
	err = r.syncWalletMembers(ctx, tx, wallet.Id, wallet.Members)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	// Save transactions
	for _, transaction := range wallet.Transactions {
		_, err = tx.ExecContext(ctx, `INSERT INTO transactions (id, wallet_id, amount_value, amount_offset, created_by, type, description, created_at)
						  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
						  ON CONFLICT (id) DO NOTHING`,
			transaction.Id,
			wallet.Id,
			transaction.Amount.Value,
			transaction.Amount.Offset,
			transaction.CreatedBy.Id,
			string(transaction.Type),
			transaction.Description,
			transaction.CreatedAt,
		)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	if err := r.publisher.Publish(ctx, wallet.Events()); err != nil {
		slog.Error("error publishing events", slog.String("error", err.Error()))
	}
	wallet.ClearEvents()

	span.SetStatus(codes.Ok, "Wallet saved")
	return err
}

func (r *PostgreSQLWalletRepository) loadWalletMembers(ctx context.Context, walletId string) ([]models.Member, error) {
	ctx, span := r.tracer.Start(ctx, "loadWalletMembers", trace.WithAttributes(
		attribute.String("wallet.id", walletId),
	))
	defer span.End()

	query := `SELECT wu.member_id
			  FROM wallet_users wu
			  WHERE wu.wallet_id = $1`

	rows, err := r.db.QueryContext(ctx, query, walletId)
	if err != nil {
		span.SetStatus(codes.Error, "query failed")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var users []models.Member

	for rows.Next() {
		var user models.Member

		err := rows.Scan(
			&user.Id,
		)
		if err != nil {
			span.SetStatus(codes.Error, "query failed")
			span.RecordError(err)
			return nil, err
		}

		span.AddEvent(fmt.Sprintf("Retrieving data for user %s", user.Id), trace.WithAttributes(
			attribute.String("user.id", user.Id),
		))

		userData, err := r.memberRepo.FindByID(ctx, user.Id)
		if err != nil {
			span.SetStatus(codes.Error, "failed to get member data")
			span.RecordError(err)
			return nil, err
		}

		user.FirstName = userData.FirstName
		user.LastName = userData.LastName
		user.Email = userData.Email
		user.CreatedAt = userData.CreatedAt
		if userData.UpdatedAt != nil {
			user.UpdatedAt = userData.UpdatedAt
		}

		span.AddEvent(fmt.Sprintf("User data retrieved for user %s", user.Id), trace.WithAttributes(
			attribute.String("user.id", user.Id),
			attribute.String("user.email", user.Email),
		))

		users = append(users, user)
	}

	return users, nil
}

func (r *PostgreSQLWalletRepository) loadWalletTransactions(ctx context.Context, walletId string) ([]models.Transaction, error) {
	ctx, span := r.tracer.Start(ctx, "loadWalletTransactions", trace.WithAttributes(
		attribute.String("wallet.id", walletId),
	))
	defer span.End()

	query := `SELECT t.id, t.amount_value, t.amount_offset, t.type, t.description, t.created_at, t.created_by
			  FROM transactions t
			  WHERE t.wallet_id = $1
			  ORDER BY t.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, walletId)
	if err != nil {
		span.SetStatus(codes.Error, "query failed")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction
		var user models.Member

		err := rows.Scan(
			&transaction.Id,
			&transaction.Amount.Value,
			&transaction.Amount.Offset,
			&transaction.Type,
			&transaction.Description,
			&transaction.CreatedAt,
			&user.Id,
		)
		if err != nil {
			span.SetStatus(codes.Error, "scan failed")
			span.RecordError(err)
			return nil, err
		}

		span.AddEvent(fmt.Sprintf("Retrieving data for user %s", user.Id), trace.WithAttributes(
			attribute.String("user.id", user.Id),
		))

		userData, err := r.memberRepo.FindByID(ctx, user.Id)
		if err != nil {
			span.SetStatus(codes.Error, "failed to retrieve user data")
			span.RecordError(err)
			return nil, err
		}

		user.FirstName = userData.FirstName
		user.LastName = userData.LastName
		user.Email = userData.Email
		user.CreatedAt = userData.CreatedAt
		if userData.UpdatedAt != nil {
			user.UpdatedAt = userData.UpdatedAt
		}

		span.AddEvent(fmt.Sprintf("User data retrieved for user %s", user.Id), trace.WithAttributes(
			attribute.String("user.id", user.Id),
			attribute.String("user.email", user.Email),
		))

		transaction.CreatedBy = user
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *PostgreSQLWalletRepository) syncWalletMembers(ctx context.Context, tx *sql.Tx, walletId string, users []models.Member) error {
	ctx, span := r.tracer.Start(ctx, "syncWalletMembers", trace.WithAttributes(
		attribute.String("wallet.id", walletId),
	))
	defer span.End()
	// Get current user IDs
	rows, err := tx.QueryContext(ctx, "SELECT member_id FROM wallet_users WHERE wallet_id = $1", walletId)
	if err != nil {
		span.SetStatus(codes.Error, "failed to query current users")
		span.RecordError(err)
		return err
	}
	defer rows.Close()

	currentUsers := make(map[string]bool)
	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			span.SetStatus(codes.Error, "failed to scan current users")
			span.RecordError(err)
			return err
		}
		currentUsers[userId] = true
	}

	// Build new user set
	newUsers := make(map[string]bool)
	for _, user := range users {
		newUsers[user.Id] = true
	}

	// Remove users not in new set
	for userId := range currentUsers {
		if !newUsers[userId] {
			_, err = tx.ExecContext(ctx, "DELETE FROM wallet_users WHERE wallet_id = $1 AND member_id = $2", walletId, userId)
			if err != nil {
				span.SetStatus(codes.Error, "failed to remove old users")
				span.RecordError(err)
				return err
			}
		}
	}

	// Add new users
	for userId := range newUsers {
		if !currentUsers[userId] {
			_, err = tx.ExecContext(ctx, "INSERT INTO wallet_users (wallet_id, member_id, assigned_at) VALUES ($1, $2, $3)", walletId, userId, time.Now())
			if err != nil {
				span.SetStatus(codes.Error, "failed to add new users")
				span.RecordError(err)
				return err
			}
		}
	}

	return nil
}
