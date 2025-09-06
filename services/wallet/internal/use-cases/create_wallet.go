package usecases

import (
	"context"
	"errors"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type CreateWalletUseCaseInput struct {
	CreatorID string
	Name      string
}

func (usecase *UseCase) CreateWallet(ctx context.Context, input CreateWalletUseCaseInput) (*models.Wallet, error) {
	if input.CreatorID == "" {
		return nil, MissingRequiredFieldsError("CreatorID")
	}

	if input.Name == "" {
		return nil, MissingRequiredFieldsError("Name")
	}

	creator, err := usecase.repos.User.FindByID(ctx, input.CreatorID)
	if err != nil {
		return nil, errors.Join(errors.New("could not find creator id"), err)
	}

	wallet := models.CreateNewWallet(input.Name, creator)

	if err := usecase.repos.Wallet.Save(ctx, wallet); err != nil {
		return nil, errors.Join(errors.New("could not persist the wallet"), err)
	}

	return wallet, nil
}
