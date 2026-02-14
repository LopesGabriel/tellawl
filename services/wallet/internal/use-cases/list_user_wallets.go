package usecases

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type ListUserWalletsUseCaseInput struct {
	UserId string
}

func (usecase *UseCase) ListUserWallets(ctx context.Context, input ListUserWalletsUseCaseInput) ([]models.Wallet, error) {
	user, err := usecase.repos.Member.FindByID(ctx, input.UserId)
	if err != nil {
		return nil, errx.ErrInvalidCreatorID
	}

	userWallets, err := usecase.repos.Wallet.FindByUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	return userWallets, nil
}
