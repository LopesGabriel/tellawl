package usecases

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type CreateWalletUseCase struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
}

type CreateWalletUseCaseInput struct {
	CreatorID string
	Name      string
}

func NewCreateWalletUseCase(userRepository repository.UserRepository, walletRepository repository.WalletRepository) *CreateWalletUseCase {
	return &CreateWalletUseCase{
		userRepository:   userRepository,
		walletRepository: walletRepository,
	}
}

func (usecase *CreateWalletUseCase) Execute(input CreateWalletUseCaseInput) (*models.Wallet, error) {
	if input.CreatorID == "" {
		return nil, MissingRequiredFieldsError("CreatorID")
	}

	if input.Name == "" {
		return nil, MissingRequiredFieldsError("Name")
	}

	creator, err := usecase.userRepository.FindByID(input.CreatorID)
	if err != nil {
		return nil, errors.Join(errors.New("could not find creator id"), err)
	}

	wallet := models.CreateNewWallet(input.Name, creator)

	if err := usecase.walletRepository.Save(wallet); err != nil {
		return nil, errors.Join(errors.New("could not persist the wallet"), err)
	}

	return wallet, nil
}
