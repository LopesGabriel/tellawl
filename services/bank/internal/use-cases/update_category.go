package usecases

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type UpdateCategoryUseCase struct {
	walletRepository repository.WalletRepository
}

type UpdateCategoryUseCaseInput struct {
	WalletId   string
	CategoryId string
	Name       string
}

func NewUpdateCategoryUseCase(walletRepository repository.WalletRepository) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{
		walletRepository: walletRepository,
	}
}

func (usecase *UpdateCategoryUseCase) Execute(input UpdateCategoryUseCaseInput) (*models.Category, error) {
	if input.WalletId == "" {
		return nil, MissingRequiredFieldsError("WalletId")
	}

	if input.CategoryId == "" {
		return nil, MissingRequiredFieldsError("CategoryId")
	}

	if input.Name == "" {
		return nil, MissingRequiredFieldsError("Name")
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