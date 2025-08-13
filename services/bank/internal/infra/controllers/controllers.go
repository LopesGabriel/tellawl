package controllers

import "github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"

type controller struct {
	version          string
	walletRepository repository.WalletRepository
	userRepository   repository.UserRepository
}

type NewControllerParams struct {
	WalletRepository repository.WalletRepository
	UserRepository   repository.UserRepository
}

func NewController(params NewControllerParams) *controller {
	return &controller{
		version:          "v1",
		walletRepository: params.WalletRepository,
		userRepository:   params.UserRepository,
	}
}
