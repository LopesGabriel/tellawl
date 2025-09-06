package usecases_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

func TestWalletCreation(t *testing.T) {
	publisher := events.InMemoryEventPublisher{}
	repos := repository.NewInMemory(publisher)
	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		JwtSecret: "examle",
		Repos:     repos,
	})

	userId := uuid.NewString()
	user := models.User{
		Id:        userId,
		FirstName: "Gabriel",
		LastName:  "Lopes",
		Email:     "example@example.com",
		CreatedAt: time.Now(),
	}
	repos.User.Save(t.Context(), &user)

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
