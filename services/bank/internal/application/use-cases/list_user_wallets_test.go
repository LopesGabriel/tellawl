package usecases

import (
	"testing"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
)

func TestListUserWalletsUseCase(t *testing.T) {
	var sut listUserWalletsUseCase
	var userRepository *database.InMemoryUserRepository
	var walletRepository *database.InMemoryWalletRepository
	eventPublisher := &events.InMemoryEventPublisher{}

	t.Run("should list user wallets", func(t *testing.T) {
		userRepository = database.NewInMemoryUserRepository(eventPublisher)
		walletRepository = database.NewInMemoryWalletRepository(eventPublisher)

		sut = listUserWalletsUseCase{
			userRepository:   userRepository,
			walletRepository: walletRepository,
		}

		user, _ := models.CreateNewUser("Gabriel", "Lopes", "gabriel@example.com", "pw1")
		user2, _ := models.CreateNewUser("Matheus", "Lopes", "matheus@example.com", "pw2")

		userRepository.Save(user)
		userRepository.Save(user2)

		wallet := models.CreateNewWallet("Test wallet", user)
		wallet.AddUser(user2)
		walletRepository.Save(wallet)

		wallet2 := models.CreateNewWallet("Test wallet 2", user2)
		walletRepository.Save(wallet2)

		wallets, err := sut.Execute(ListUserWalletsUseCaseInput{
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
