package usecases_test

import (
	"testing"
	"time"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/publisher"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	lognoop "go.opentelemetry.io/otel/log/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func TestListUserWalletsUseCase(t *testing.T) {
	appLogger, err := logger.Init(t.Context(), logger.InitLoggerArgs{
		LoggerProvider: lognoop.NewLoggerProvider(),
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Shutdown(t.Context())

	eventPublisher := publisher.NewInMemoryEventPublisher(appLogger)

	t.Run("should list user wallets", func(t *testing.T) {
		repos := database.NewInMemory(eventPublisher)
		memberRepo := database.NewInMemoryMemberRepository(eventPublisher)
		repos.Member = memberRepo
		useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
			Repos:  repos,
			Tracer: tracenoop.NewTracerProvider().Tracer("test"),
			Logger: appLogger,
		})

		user := createMember("user1", "Gabriel", "Lopes", "gabriel@example.com")
		user2 := createMember("user2", "Matheus", "Lopes", "matheus@example.com")

		memberRepo.Items = append(memberRepo.Items, *user)
		memberRepo.Items = append(memberRepo.Items, *user2)

		wallet := models.CreateNewWallet("Test wallet", user)
		wallet.AddUser(user2)
		repos.Wallet.Save(t.Context(), wallet)

		wallet2 := models.CreateNewWallet("Test wallet 2", user2)
		repos.Wallet.Save(t.Context(), wallet2)

		wallets, err := useCases.ListUserWallets(t.Context(), usecases.ListUserWalletsUseCaseInput{
			UserId: user2.Id,
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(wallets) != 2 {
			t.Errorf("Expected 2 wallets, got %v", len(wallets))
		}
	})
}

func createMember(id, firstName, lastName, email string) *models.Member {
	return &models.Member{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: time.Now(),
	}
}
