package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

type CreateUserUseCaseInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (usecase *UseCase) CreateUser(ctx context.Context, input CreateUserUseCaseInput) (*models.User, error) {
	if input.FirstName == "" {
		return nil, MissingRequiredFieldsError("FirstName")
	}

	if input.LastName == "" {
		return nil, MissingRequiredFieldsError("LastName")
	}

	if input.Email == "" {
		return nil, MissingRequiredFieldsError("Email")
	}

	if input.Password == "" {
		return nil, MissingRequiredFieldsError("Password")
	}

	existingUser, err := usecase.repos.User.FindByEmail(ctx, input.Email)
	if err != nil {
		if !strings.Contains(err.Error(), "user not found") {
			return nil, errors.Join(errors.New("could not validate existing user"), err)
		}
	}

	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	user, err := models.CreateNewUser(input.FirstName, input.LastName, input.Email, input.Password)
	if err != nil {
		return nil, errors.Join(errors.New("could not create user"), err)
	}

	if err := usecase.repos.User.Save(ctx, user); err != nil {
		return nil, errors.Join(errors.New("could not persist the user"), err)
	}

	return user, nil
}
