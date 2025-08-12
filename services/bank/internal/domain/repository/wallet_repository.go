package repository

import "github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"

type WalletRepository interface {
	Save(wallet *models.Wallet) error
}
