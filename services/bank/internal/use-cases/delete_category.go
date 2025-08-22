package usecases

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type DeleteCategoryUseCase struct {
	walletRepository repository.WalletRepository
}

type DeleteCategoryUseCaseInput struct {
	WalletId   string
	CategoryId string
}

func NewDeleteCategoryUseCase(walletRepository repository.WalletRepository) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{
		walletRepository: walletRepository,
	}
}

func (usecase *DeleteCategoryUseCase) Execute(input DeleteCategoryUseCaseInput) error {
	if input.WalletId == "" {
		return MissingRequiredFieldsError("WalletId")
	}

	if input.CategoryId == "" {
		return MissingRequiredFieldsError("CategoryId")
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