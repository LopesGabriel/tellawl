package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

func TestRegisterTransactionUseCase(t *testing.T) {
	var repos *repository.Repositories
	eventPublisher := &events.InMemoryEventPublisher{}
	useCaseArgs := usecases.NewUseCasesArgs{
		JwtSecret: "example",
		Repos:     repos,
	}

	t.Run("should register a transaction", func(t *testing.T) {
		repos = repository.NewInMemory(eventPublisher)
		useCases := usecases.NewUseCases(useCaseArgs)

		user, _ := models.CreateNewUser("Matheus", "Lopes", "matheus@example.com", "pw2")
		repos.User.Save(user)

		wallet := models.CreateNewWallet("Test wallet", user)
		repos.Wallet.Save(wallet)

		transaction, err := useCases.RegisterTransaction(usecases.RegisterTransactionUseCaseInput{
			TransactionRegisteredByUserId: user.Id,
			WalletId:                      wallet.Id,
			Amount:                        1000000,
			TransactionType:               "deposit",
			Description:                   "Test salary deposit",
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if transaction.Amount.Value != 1000000 {
			t.Errorf("Expected transaction amount to be 1000000, got %v", transaction.Amount.Value)
		}
		if transaction.Amount.Offset != 100 {
			t.Errorf("Expected transaction offset to be 100, got %v", transaction.Amount.Offset)
		}

		wallet, err = repos.Wallet.FindById(wallet.Id)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if wallet.Balance.Value != 1000000 {
			t.Errorf("Expected wallet balance to be 1000000, got %v", wallet.Balance.Value)
		}

		if len(wallet.Transactions) != 1 {
			t.Errorf("Expected wallet to have 1 transaction, got %v", len(wallet.Transactions))
		}
	})

	t.Run("should register a transaction with custom offset", func(t *testing.T) {
		repos = repository.NewInMemory(eventPublisher)
		useCases := usecases.NewUseCases(useCaseArgs)

		user, _ := models.CreateNewUser("Matheus", "Lopes", "matheus@example.com", "pw2")
		repos.User.Save(user)

		wallet := models.CreateNewWallet("Test wallet", user)
		repos.Wallet.Save(wallet)

		transaction, err := useCases.RegisterTransaction(usecases.RegisterTransactionUseCaseInput{
			TransactionRegisteredByUserId: user.Id,
			WalletId:                      wallet.Id,
			Amount:                        3000000,
			Offset:                        1000,
			TransactionType:               "deposit",
			Description:                   "Custom offset deposit",
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if transaction.Amount.Value != 3000000 {
			t.Errorf("Expected transaction amount to be 3000000, got %v", transaction.Amount.Value)
		}
		if transaction.Amount.Offset != 1000 {
			t.Errorf("Expected transaction offset to be 1000, got %v", transaction.Amount.Offset)
		}

		wallet, err = repos.Wallet.FindById(wallet.Id)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if wallet.Balance.Value != 300000 {
			t.Errorf("Expected wallet balance to be 300000, got %v", wallet.Balance.Value)
		}

		if len(wallet.Transactions) != 1 {
			t.Errorf("Expected wallet to have 1 transaction, got %v", len(wallet.Transactions))
		}
	})

	t.Run("User with no access should not register a transaction", func(t *testing.T) {
		repos = repository.NewInMemory(eventPublisher)
		useCases := usecases.NewUseCases(useCaseArgs)

		user, _ := models.CreateNewUser("Gabriel", "Lopes", "gabriel@example.com", "pw1")
		user2, _ := models.CreateNewUser("Matheus", "Lopes", "matheus@example.com", "pw2")
		repos.User.Save(user)
		repos.User.Save(user2)

		wallet := models.CreateNewWallet("Test wallet", user)
		repos.Wallet.Save(wallet)

		_, err := useCases.RegisterTransaction(usecases.RegisterTransactionUseCaseInput{
			TransactionRegisteredByUserId: user2.Id,
			WalletId:                      wallet.Id,
			Amount:                        10000,
			Offset:                        1000,
			TransactionType:               "deposit",
			Description:                   "Unauthorized attempt",
		})

		if err != usecases.ErrInsufficientPermissions {
			t.Errorf("Expected insufficient permissions error, got %v", err)
		}

		wallet, err = repos.Wallet.FindById(wallet.Id)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if wallet.Balance.Value != 0 {
			t.Errorf("Expected wallet balance to be 0, got %v", wallet.Balance.Value)
		}

		if len(wallet.Transactions) != 0 {
			t.Errorf("Expected wallet to have 0 transaction, got %v", len(wallet.Transactions))
		}
	})

	t.Run("User with access should be able to register a transaction", func(t *testing.T) {
		repos = repository.NewInMemory(eventPublisher)
		useCases := usecases.NewUseCases(useCaseArgs)

		user, _ := models.CreateNewUser("Gabriel", "Lopes", "gabriel@example.com", "pw1")
		user2, _ := models.CreateNewUser("Matheus", "Lopes", "matheus@example.com", "pw2")
		repos.User.Save(user)
		repos.User.Save(user2)

		wallet := models.CreateNewWallet("Test wallet", user)
		repos.Wallet.Save(wallet)
		wallet.AddUser(user2)
		repos.Wallet.Save(wallet)

		_, err := useCases.RegisterTransaction(usecases.RegisterTransactionUseCaseInput{
			TransactionRegisteredByUserId: user2.Id,
			WalletId:                      wallet.Id,
			Amount:                        350000,
			TransactionType:               "deposit",
			Description:                   "Authorized user deposit",
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		wallet, err = repos.Wallet.FindById(wallet.Id)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if wallet.Balance.Value != 350000 {
			t.Errorf("Expected wallet balance to be 350000, got %v", wallet.Balance.Value)
		}

		if len(wallet.Transactions) != 1 {
			t.Errorf("Expected wallet to have 1 transaction, got %v", len(wallet.Transactions))
		}
	})
}
