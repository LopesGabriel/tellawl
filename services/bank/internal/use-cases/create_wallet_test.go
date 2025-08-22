package usecases_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

func TestWalletCreation(t *testing.T) {
	publisher := events.InMemoryEventPublisher{}
	userRepository := database.NewInMemoryUserRepository(publisher)
	walletRepository := database.NewInMemoryWalletRepository(publisher)
	sut := usecases.NewCreateWalletUseCase(userRepository, walletRepository)

	userId := uuid.NewString()
	user := models.User{
		Id:        userId,
		FirstName: "Gabriel",
		LastName:  "Lopes",
		Email:     "example@example.com",
		CreatedAt: time.Now(),
	}
	userRepository.Items = append(userRepository.Items, user)

	wallet, err := sut.Execute(usecases.CreateWalletUseCaseInput{
		CreatorID: userId,
		Name:      "My Wallet",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if wallet.Name != "My Wallet" {
		t.Errorf("Expected wallet name to be 'My Wallet', got %v", wallet.Name)
	}

	if len(walletRepository.Items) < 1 {
		t.Errorf("Expected at least one wallet, got %v", len(walletRepository.Items))
	}
}
