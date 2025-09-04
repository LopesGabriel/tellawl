package usecases

import (
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

type ShareWalletUseCaseInput struct {
	WalletCreatorId string
	WalletId        string
	SharedUserEmail string
}

func (usecase *UseCase) ShareWallet(input ShareWalletUseCaseInput) (*models.Wallet, error) {
	creatorUser, err := usecase.repos.User.FindByID(input.WalletCreatorId)
	if err != nil {
		return nil, err
	}

	wallet, err := usecase.repos.Wallet.FindById(input.WalletId)
	if err != nil {
		return nil, err
	}

	if wallet.CreatorId != creatorUser.Id {
		return nil, ErrInsufficientPermissions
	}

	sharedUser, err := usecase.repos.User.FindByEmail(input.SharedUserEmail)
	if err != nil {
		return nil, err
	}

	wallet.AddUser(sharedUser)
	err = usecase.repos.Wallet.Save(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
