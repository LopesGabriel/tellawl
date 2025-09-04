package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

func TestShareWalletUseCase(t *testing.T) {
	publisher := events.InMemoryEventPublisher{}
	repos := repository.NewInMemory(publisher)
	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		JwtSecret: "example",
		Repos:     repos,
	})

	user1, _ := models.CreateNewUser("Gabriel", "Lopes", "gabriel@example.com", "pw1")
	user2, _ := models.CreateNewUser("Matheus", "Lopes", "matheus@example.com", "pw2")
	repos.User.Save(t.Context(), user1)
	repos.User.Save(t.Context(), user2)

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

	if len(updatedWallet.Users) != 2 {
		t.Errorf("Expected wallet to have 2 users, got %v", len(wallet.Users))
	}
}
