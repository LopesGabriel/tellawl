package usecases

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type RegisterTransactionUseCaseInput struct {
	TransactionRegisteredByUserId string
	TransactionRegisteredByUser   *models.Member
	WalletId                      string
	Amount                        int
	Offset                        int
	TransactionType               string
	Description                   string
}

func (usecase *UseCase) RegisterTransaction(ctx context.Context, input RegisterTransactionUseCaseInput) (*models.Transaction, error) {
	if input.Offset == 0 {
		input.Offset = 100
	}

	if models.TransactionType(input.TransactionType) != models.TransactionTypeDeposit && models.TransactionType(input.TransactionType) != models.TransactionTypeWithdraw {
		return nil, errx.ErrInvalidTransactionType
	}

	var user *models.Member
	if input.TransactionRegisteredByUser != nil {
		user = input.TransactionRegisteredByUser
	} else {
		member, err := usecase.repos.Member.FindByID(ctx, input.TransactionRegisteredByUserId)
		if err != nil {
			return nil, err
		}
		user = member
	}

	wallet, err := usecase.repos.Wallet.FindById(ctx, input.WalletId)
	if err != nil {
		return nil, err
	}

	transaction, err := wallet.RegisterNewTransaction(
		models.Monetary{Value: input.Amount, Offset: input.Offset},
		*user,
		models.TransactionType(input.TransactionType),
		input.Description,
	)
	if err != nil {
		if err.Error() == "user is not allowed to register transactions" {
			return nil, errx.ErrInsufficientPermissions
		}

		return nil, err
	}

	err = usecase.repos.Wallet.Save(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
