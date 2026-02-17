package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/publisher"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	lognoop "go.opentelemetry.io/otel/log/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func TestShareWalletUseCase(t *testing.T) {
	appLogger, err := logger.Init(t.Context(), logger.InitLoggerArgs{
		LoggerProvider: lognoop.NewLoggerProvider(),
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Shutdown(t.Context())

	eventPublisher := publisher.NewInMemoryEventPublisher(appLogger)
	memberRepo := database.NewInMemoryMemberRepository(eventPublisher)
	repos := database.NewInMemory(eventPublisher)
	repos.Member = memberRepo
	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		Repos:  repos,
		Tracer: tracenoop.NewTracerProvider().Tracer("test"),
		Logger: appLogger,
	})

	user1 := createMember("member1", "Gabriel", "Lopes", "gabriel@example.com")
	user2 := createMember("member2", "Matheus", "Lopes", "matheus@example.com")
	memberRepo.Items = append(memberRepo.Items, *user1)
	memberRepo.Items = append(memberRepo.Items, *user2)

	wallet := models.CreateNewWallet("Test wallet", user1)
	repos.Wallet.Save(t.Context(), wallet)

	updatedWallet, err := useCases.ShareWallet(t.Context(), usecases.ShareWalletUseCaseInput{
		WalletCreatorId: user1.Id,
		WalletId:        wallet.Id,
		SharedUserEmail: user2.Email,
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(updatedWallet.Members) != 2 {
		t.Errorf("Expected wallet to have 2 users, got %v", len(updatedWallet.Members))
	}
}
