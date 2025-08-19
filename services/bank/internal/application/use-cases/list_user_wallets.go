package usecases

import (
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type listUserWalletsUseCase struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
}

type ListUserWalletsUseCaseInput struct {
	UserId string
}

func NewListUserWalletsUseCase(userRepository repository.UserRepository, walletRepository repository.WalletRepository) *listUserWalletsUseCase {
	return &listUserWalletsUseCase{
		userRepository:   userRepository,
		walletRepository: walletRepository,
	}
}

func (usecase *listUserWalletsUseCase) Execute(input ListUserWalletsUseCaseInput) ([]*models.Wallet, error) {
	user, err := usecase.userRepository.FindByID(input.UserId)
	if err != nil {
		return nil, ErrInvalidCreatorID
	}

	userWallets, err := usecase.walletRepository.FindByUserId(user.Id)
	if err != nil {
		return nil, err
	}

	return userWallets, nil
}
