package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	inmemory "github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database/in_memory"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	lognoop "go.opentelemetry.io/otel/log/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func TestRegisterTransactionUseCase(t *testing.T) {
	appLogger, err := logger.Init(t.Context(), logger.InitLoggerArgs{
		LoggerProvider: lognoop.NewLoggerProvider(),
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Shutdown(t.Context())

	eventPublisher := events.NewInMemoryEventPublisher(appLogger)
	repos := repository.NewInMemory(eventPublisher)
	memberRepo := inmemory.NewInMemoryMemberRepository(eventPublisher)
	repos.Member = memberRepo

	t.Run("should register a transaction", func(t *testing.T) {
		useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
			Repos:  repos,
			Tracer: tracenoop.NewTracerProvider().Tracer("test"),
			Logger: appLogger,
		})

		user := createMember("member1", "Matheus", "Lopes", "matheus@example.com")
		memberRepo.Items = append(memberRepo.Items, *user)

		wallet := models.CreateNewWallet("Test wallet", user)
		repos.Wallet.Save(t.Context(), wallet)

		transaction, err := useCases.RegisterTransaction(t.Context(), usecases.RegisterTransactionUseCaseInput{
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

		wallet, err = repos.Wallet.FindById(t.Context(), wallet.Id)
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
		repos := repository.NewInMemory(eventPublisher)
		memberRepo := inmemory.NewInMemoryMemberRepository(eventPublisher)
		repos.Member = memberRepo
		useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
			Repos:  repos,
			Tracer: tracenoop.NewTracerProvider().Tracer("test"),
			Logger: appLogger,
		})

		user := createMember("member1", "Matheus", "Lopes", "matheus@example.com")
		memberRepo.Items = append(memberRepo.Items, *user)

		wallet := models.CreateNewWallet("Test wallet", user)
		repos.Wallet.Save(t.Context(), wallet)

		transaction, err := useCases.RegisterTransaction(t.Context(), usecases.RegisterTransactionUseCaseInput{
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

		wallet, err = repos.Wallet.FindById(t.Context(), wallet.Id)
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
		repos := repository.NewInMemory(eventPublisher)
		memberRepo := inmemory.NewInMemoryMemberRepository(eventPublisher)
		repos.Member = memberRepo
		useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
			Repos:  repos,
			Tracer: tracenoop.NewTracerProvider().Tracer("test"),
			Logger: appLogger,
		})

		user := createMember("member1", "Gabriel", "Lopes", "gabriel@example.com")
		user2 := createMember("member2", "Matheus", "Lopes", "matheus@example.com")
		memberRepo.Items = append(memberRepo.Items, *user)
		memberRepo.Items = append(memberRepo.Items, *user2)

		wallet := models.CreateNewWallet("Test wallet", user)
		repos.Wallet.Save(t.Context(), wallet)

		_, err := useCases.RegisterTransaction(t.Context(), usecases.RegisterTransactionUseCaseInput{
			TransactionRegisteredByUserId: user2.Id,
			WalletId:                      wallet.Id,
			Amount:                        10000,
			Offset:                        1000,
			TransactionType:               "deposit",
			Description:                   "Unauthorized attempt",
		})

		if err != errx.ErrInsufficientPermissions {
			t.Errorf("Expected insufficient permissions error, got %v", err)
		}

		wallet, err = repos.Wallet.FindById(t.Context(), wallet.Id)
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
		repos := repository.NewInMemory(eventPublisher)
		memberRepo := inmemory.NewInMemoryMemberRepository(eventPublisher)
		repos.Member = memberRepo
		useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
			Repos:  repos,
			Tracer: tracenoop.NewTracerProvider().Tracer("test"),
			Logger: appLogger,
		})

		user := createMember("member1", "Gabriel", "Lopes", "gabriel@example.com")
		user2 := createMember("member2", "Matheus", "Lopes", "matheus@example.com")
		memberRepo.Items = append(memberRepo.Items, *user)
		memberRepo.Items = append(memberRepo.Items, *user2)

		wallet := models.CreateNewWallet("Test wallet", user)
		repos.Wallet.Save(t.Context(), wallet)
		wallet.AddUser(user2)
		repos.Wallet.Save(t.Context(), wallet)

		_, err := useCases.RegisterTransaction(t.Context(), usecases.RegisterTransactionUseCaseInput{
			TransactionRegisteredByUserId: user2.Id,
			WalletId:                      wallet.Id,
			Amount:                        350000,
			TransactionType:               "deposit",
			Description:                   "Authorized user deposit",
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		wallet, err = repos.Wallet.FindById(t.Context(), wallet.Id)
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
