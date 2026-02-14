package usecases

import (
	"context"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
)

type AuthenticateUserUseCaseInput struct {
	Token string
}

func (usecase *UseCase) AuthenticateUser(ctx context.Context, input AuthenticateUserUseCaseInput) (*models.Member, error) {
	if input.Token == "" {
		return nil, errx.ErrInvalidInput
	}

	member, err := usecase.repos.Member.ValidateToken(ctx, input.Token)

	return member, err
}
