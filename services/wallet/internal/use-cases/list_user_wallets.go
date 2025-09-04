package usecases

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

type ListUserWalletsUseCaseInput struct {
	UserId string
}

func (usecase *UseCase) ListUserWallets(ctx context.Context, input ListUserWalletsUseCaseInput) ([]models.Wallet, error) {
	user, err := usecase.repos.User.FindByID(ctx, input.UserId)
	if err != nil {
		return nil, ErrInvalidCreatorID
	}

	userWallets, err := usecase.repos.Wallet.FindByUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return userWallets, nil
}
