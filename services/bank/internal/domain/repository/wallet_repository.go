package repository

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

var ErrWalletNotFound = errors.New("wallet not found")

type WalletRepository interface {
	FindById(id string) (*models.Wallet, error)
	Save(wallet *models.Wallet) error
}
