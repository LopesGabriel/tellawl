package usecases

import (
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type registerTransactionUseCase struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
}

type RegisterTransactionUseCaseInput struct {
	TransactionRegisteredByUserId string
	WalletId                      string
	Amount                        int
	CategoryId                    string
	Offset                        int
	TransactionType               string
	Description                   string
}

func NewRegisterTransactionUseCase(userRepository repository.UserRepository, walletRepository repository.WalletRepository) *registerTransactionUseCase {
	return &registerTransactionUseCase{
		userRepository:   userRepository,
		walletRepository: walletRepository,
	}
}

func (usecase *registerTransactionUseCase) Execute(input RegisterTransactionUseCaseInput) (*models.Transaction, error) {
	if input.Offset == 0 {
		input.Offset = 100
	}

	if models.TransactionType(input.TransactionType) != models.TransactionTypeDeposit && models.TransactionType(input.TransactionType) != models.TransactionTypeWithdraw {
		return nil, ErrInvalidTransactionType
	}

	user, err := usecase.userRepository.FindByID(input.TransactionRegisteredByUserId)
	if err != nil {
		return nil, err
	}

	wallet, err := usecase.walletRepository.FindById(input.WalletId)
	if err != nil {
		return nil, err
	}

	transaction, err := wallet.RegisterNewTransaction(
		models.Monetary{Value: input.Amount, Offset: input.Offset},
		*user,
		models.TransactionType(input.TransactionType),
		input.CategoryId,
		input.Description,
	)
	if err != nil {
		if err.Error() == "user is not allowed to register transactions" {
			return nil, ErrInsufficientPermissions
		}

		return nil, err
	}

	err = usecase.walletRepository.Save(wallet)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
