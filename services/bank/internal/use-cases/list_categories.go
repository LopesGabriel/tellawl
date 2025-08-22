package usecases

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type ListCategoriesUseCase struct {
	walletRepository repository.WalletRepository
}

type ListCategoriesUseCaseInput struct {
	WalletId string
}

func NewListCategoriesUseCase(walletRepository repository.WalletRepository) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		walletRepository: walletRepository,
	}
}

func (usecase *ListCategoriesUseCase) Execute(input ListCategoriesUseCaseInput) ([]models.Category, error) {
	if input.WalletId == "" {
		return nil, MissingRequiredFieldsError("WalletId")
	}

	wallet, err := usecase.walletRepository.FindById(input.WalletId)
	if err != nil {
		return nil, errors.Join(errors.New("wallet not found"), err)
	}

	return wallet.GetCategories(), nil
}