package category

import (
	"fmt"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type getCategoryUseCase struct {
	walletRepository repository.WalletRepository
}

type GetCategoryUseCaseInput struct {
	WalletId   string
	CategoryId string
}

func NewGetCategoryUseCase(walletRepository repository.WalletRepository) *getCategoryUseCase {
	return &getCategoryUseCase{
		walletRepository: walletRepository,
	}
}

func (usecase *getCategoryUseCase) Execute(input GetCategoryUseCaseInput) (*models.Category, error) {
	wallet, err := usecase.walletRepository.FindById(input.WalletId)
	if err != nil {
		return nil, err
	}

	var category models.Category
	categories := wallet.GetCategories()
	for _, c := range categories {
		if c.Id == input.CategoryId {
			category = c
			break
		}
	}

	if category.Id == "" {
		return nil, fmt.Errorf("category %s not found", input.CategoryId)
	}

	return &category, nil
}
