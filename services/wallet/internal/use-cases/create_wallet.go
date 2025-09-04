package usecases

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

type CreateWalletUseCaseInput struct {
	CreatorID string
	Name      string
}

func (usecase *UseCase) CreateWallet(input CreateWalletUseCaseInput) (*models.Wallet, error) {
	if input.CreatorID == "" {
		return nil, MissingRequiredFieldsError("CreatorID")
	}

	if input.Name == "" {
		return nil, MissingRequiredFieldsError("Name")
	}

	creator, err := usecase.repos.User.FindByID(input.CreatorID)
	if err != nil {
		return nil, errors.Join(errors.New("could not find creator id"), err)
	}

	wallet := models.CreateNewWallet(input.Name, creator)

	if err := usecase.repos.Wallet.Save(wallet); err != nil {
		return nil, errors.Join(errors.New("could not persist the wallet"), err)
	}

	return wallet, nil
}
