package usecases

import (
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type shareWalletUseCase struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
}

type ShareWalletUseCaseInput struct {
	WalletCreatorId string
	WalletId        string
	SharedUserEmail string
}

func NewShareWalletUseCase(userRepository repository.UserRepository, walletRepository repository.WalletRepository) *shareWalletUseCase {
	return &shareWalletUseCase{
		userRepository:   userRepository,
		walletRepository: walletRepository,
	}
}

func (usecase *shareWalletUseCase) Execute(input ShareWalletUseCaseInput) (*models.Wallet, error) {
	creatorUser, err := usecase.userRepository.FindByID(input.WalletCreatorId)
	if err != nil {
		return nil, err
	}

	wallet, err := usecase.walletRepository.FindById(input.WalletId)
	if err != nil {
		return nil, err
	}

	if wallet.CreatorId != creatorUser.Id {
		return nil, ErrInsufficientPermissions
	}

	sharedUser, err := usecase.userRepository.FindByEmail(input.SharedUserEmail)
	if err != nil {
		return nil, err
	}

	wallet.AddUser(sharedUser)
	err = usecase.walletRepository.Save(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
