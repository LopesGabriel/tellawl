package category

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

type deleteCategoryUseCase struct {
	walletRepository repository.WalletRepository
}

type DeleteCategoryUseCaseInput struct {
	WalletId   string
	CategoryId string
}

func NewDeleteCategoryUseCase(walletRepository repository.WalletRepository) *deleteCategoryUseCase {
	return &deleteCategoryUseCase{
		walletRepository: walletRepository,
	}
}

func (usecase *deleteCategoryUseCase) Execute(input DeleteCategoryUseCaseInput) error {
	if input.WalletId == "" {
		return usecases.MissingRequiredFieldsError("WalletId")
	}

	if input.CategoryId == "" {
		return usecases.MissingRequiredFieldsError("CategoryId")
	}

	wallet, err := usecase.walletRepository.FindById(input.WalletId)
	if err != nil {
		return errors.Join(errors.New("wallet not found"), err)
	}

	if err := wallet.DeleteCategory(input.CategoryId); err != nil {
		return err
	}

	if err := usecase.walletRepository.Save(wallet); err != nil {
		return errors.Join(errors.New("could not persist the wallet"), err)
	}

	return nil
}
