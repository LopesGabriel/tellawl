package usecases

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type ShareWalletUseCaseInput struct {
	WalletCreatorId string
	WalletCreator   *models.Member
	WalletId        string
	SharedUserEmail string
}

func (usecase *UseCase) ShareWallet(ctx context.Context, input ShareWalletUseCaseInput) (*models.Wallet, error) {
	var creator *models.Member

	if input.WalletCreator != nil {
		creator = input.WalletCreator
	} else {
		creatorUser, err := usecase.repos.Member.FindByID(ctx, input.WalletCreatorId)
		if err != nil {
			return nil, err
		}
		creator = creatorUser
	}

	wallet, err := usecase.repos.Wallet.FindById(ctx, input.WalletId)
	if err != nil {
		return nil, err
	}

	if wallet.CreatorId != creator.Id {
		return nil, errx.ErrInsufficientPermissions
	}

	sharedUser, err := usecase.repos.Member.FindByEmail(ctx, input.SharedUserEmail)
	if err != nil {
		return nil, err
	}

	wallet.AddUser(sharedUser)
	err = usecase.repos.Wallet.Save(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
