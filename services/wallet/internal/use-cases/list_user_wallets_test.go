package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

func TestListUserWalletsUseCase(t *testing.T) {
	var repos *repository.Repositories
	eventPublisher := &events.InMemoryEventPublisher{}

	t.Run("should list user wallets", func(t *testing.T) {
		repos = repository.NewInMemory(eventPublisher)
		useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
			JwtSecret: "example",
			Repos:     repos,
		})

		user, _ := models.CreateNewUser("Gabriel", "Lopes", "gabriel@example.com", "pw1")
		user2, _ := models.CreateNewUser("Matheus", "Lopes", "matheus@example.com", "pw2")

		repos.User.Save(t.Context(), user)
		repos.User.Save(t.Context(), user2)

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
