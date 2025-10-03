package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

func TestShareWalletUseCase(t *testing.T) {
	publisher := events.InMemoryEventPublisher{}
	repos := repository.NewInMemory(publisher)
	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		JwtSecret: "example",
		Repos:     repos,
	})

	user1 := createMember("member1", "Gabriel", "Lopes", "gabriel@example.com")
	user2 := createMember("member2", "Matheus", "Lopes", "matheus@example.com")
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

	if len(updatedWallet.Members) != 2 {
		t.Errorf("Expected wallet to have 2 users, got %v", len(wallet.Members))
	}
}
