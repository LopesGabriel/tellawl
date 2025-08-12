package repository

import "github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"

type UserRepository interface {
	FindByID(id string) (*models.User, error)
}
