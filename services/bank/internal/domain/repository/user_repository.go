package repository

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Save(user *models.User) error
}
