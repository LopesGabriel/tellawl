package usecases

import (
	"context"
	"errors"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type CreateUserUseCaseInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (usecase *UseCase) CreateUser(ctx context.Context, input CreateUserUseCaseInput) (*models.User, error) {
	if input.FirstName == "" {
		return nil, errx.MissingRequiredFieldsError("FirstName")
	}

	if input.LastName == "" {
		return nil, errx.MissingRequiredFieldsError("LastName")
	}

	if input.Email == "" {
		return nil, errx.MissingRequiredFieldsError("Email")
	}

	if input.Password == "" {
		return nil, errx.MissingRequiredFieldsError("Password")
	}

	existingUser, err := usecase.repos.User.FindByEmail(ctx, input.Email)
	if err != nil {
		if err != errx.ErrNotFound {
			return nil, err
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
