package usecases

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
)

type AuthenticateUserUseCaseInput struct {
	Token string
}

func (usecase *UseCase) AuthenticateUser(ctx context.Context, input AuthenticateUserUseCaseInput) (string, error) {
	if input.Token == "" {
		return "", errx.ErrInvalidInput
	}

	return "", nil
}
