package category

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

type updateCategoryUseCase struct {
	walletRepository repository.WalletRepository
}

type UpdateCategoryUseCaseInput struct {
	WalletId   string
	CategoryId string
	Name       string
}

func NewUpdateCategoryUseCase(walletRepository repository.WalletRepository) *updateCategoryUseCase {
	return &updateCategoryUseCase{
		walletRepository: walletRepository,
	}
}

func (usecase *updateCategoryUseCase) Execute(input UpdateCategoryUseCaseInput) (*models.Category, error) {
	if input.WalletId == "" {
		return nil, usecases.MissingRequiredFieldsError("WalletId")
	}

	if input.CategoryId == "" {
		return nil, usecases.MissingRequiredFieldsError("CategoryId")
	}

	if input.Name == "" {
		return nil, usecases.MissingRequiredFieldsError("Name")
	}

	wallet, err := usecase.walletRepository.FindById(input.WalletId)
	if err != nil {
		return nil, errors.Join(errors.New("wallet not found"), err)
	}

	category, err := wallet.UpdateCategory(input.CategoryId, input.Name)
	if err != nil {
		return nil, err
	}

	if err := usecase.walletRepository.Save(wallet); err != nil {
		return nil, errors.Join(errors.New("could not persist the wallet"), err)
	}

	return category, nil
}
