package usecases_test

import (
	"testing"
	"time"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
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

		user := createMember("user1", "Gabriel", "Lopes", "gabriel@example.com")
		user2 := createMember("user2", "Matheus", "Lopes", "matheus@example.com")

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

func createMember(id, firstName, lastName, email string) *models.Member {
	return &models.Member{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: time.Now(),
	}
}
