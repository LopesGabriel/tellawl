package usecases_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/publisher"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	lognoop "go.opentelemetry.io/otel/log/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

func TestWalletCreation(t *testing.T) {
	appLogger, err := logger.Init(t.Context(), logger.InitLoggerArgs{
		LoggerProvider: lognoop.NewLoggerProvider(),
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Shutdown(t.Context())

	publisher := publisher.NewInMemoryEventPublisher(appLogger)
	repos := database.NewInMemory(publisher)
	memberRepo := database.NewInMemoryMemberRepository(publisher)
	repos.Member = memberRepo

	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		Repos:  repos,
		Tracer: tracenoop.NewTracerProvider().Tracer("test"),
		Logger: appLogger,
	})

	userId := uuid.NewString()
	user := models.Member{
		Id:        userId,
		FirstName: "Gabriel",
		LastName:  "Lopes",
		Email:     "example@example.com",
		CreatedAt: time.Now(),
	}
	memberRepo.Items = append(memberRepo.Items, user)

	wallet, err := useCases.CreateWallet(t.Context(), usecases.CreateWalletUseCaseInput{
		CreatorID: userId,
		Name:      "My Wallet",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if wallet.Name != "My Wallet" {
		t.Errorf("Expected wallet name to be 'My Wallet', got %v", wallet.Name)
	}

	persistedWallet, err := repos.Wallet.FindById(t.Context(), wallet.Id)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if persistedWallet.Name != "My Wallet" {
		t.Errorf("Expected persisted wallet name to be 'My Wallet', got %v", persistedWallet.Name)
	}
}
