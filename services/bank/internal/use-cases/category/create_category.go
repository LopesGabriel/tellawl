package category

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

type createCategoryUseCase struct {
	walletRepository repository.WalletRepository
}

type CreateCategoryUseCaseInput struct {
	WalletId string
	Name     string
}

func NewCreateCategoryUseCase(walletRepository repository.WalletRepository) *createCategoryUseCase {
	return &createCategoryUseCase{
		walletRepository: walletRepository,
	}
}

func (usecase *createCategoryUseCase) Execute(input CreateCategoryUseCaseInput) (*models.Category, error) {
	if input.WalletId == "" {
		return nil, usecases.MissingRequiredFieldsError("WalletId")
	}

	if input.Name == "" {
		return nil, usecases.MissingRequiredFieldsError("Name")
	}

	wallet, err := usecase.walletRepository.FindById(input.WalletId)
	if err != nil {
		return nil, errors.Join(errors.New("wallet not found"), err)
	}

	category, err := wallet.AddCustomCategory(input.Name)
	if err != nil {
		return nil, err
	}

	if err := usecase.walletRepository.Save(wallet); err != nil {
		return nil, errors.Join(errors.New("could not persist the wallet"), err)
	}

	return category, nil
}
