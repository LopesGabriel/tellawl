package category

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

type listCategoriesUseCase struct {
	walletRepository repository.WalletRepository
}

type ListCategoriesUseCaseInput struct {
	WalletId string
}

func NewListCategoriesUseCase(walletRepository repository.WalletRepository) *listCategoriesUseCase {
	return &listCategoriesUseCase{
		walletRepository: walletRepository,
	}
}

func (usecase *listCategoriesUseCase) Execute(input ListCategoriesUseCaseInput) ([]models.Category, error) {
	if input.WalletId == "" {
		return nil, usecases.MissingRequiredFieldsError("WalletId")
	}

	wallet, err := usecase.walletRepository.FindById(input.WalletId)
	if err != nil {
		return nil, errors.Join(errors.New("wallet not found"), err)
	}

	return wallet.GetCategories(), nil
}
