package usecases

import (
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

type ListUserWalletsUseCaseInput struct {
	UserId string
}

func (usecase *UseCase) ListUserWallets(input ListUserWalletsUseCaseInput) ([]models.Wallet, error) {
	user, err := usecase.repos.User.FindByID(input.UserId)
	if err != nil {
		return nil, ErrInvalidCreatorID
	}

	userWallets, err := usecase.repos.Wallet.FindByUserId(user.Id)
	if err != nil {
		return nil, err
	}

	return userWallets, nil
}
