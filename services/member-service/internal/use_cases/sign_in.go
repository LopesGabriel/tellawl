package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database"
	"go.opentelemetry.io/otel/codes"
)

type SignInUseCaseInput struct {
	Email    string
	Password string
}

type SignInUseCaseOutput struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (uc *UseCases) SignIn(ctx context.Context, input SignInUseCaseInput) (*SignInUseCaseOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "SignInUseCase")
	defer span.End()

	member, err := uc.repos.Members.FindByEmail(ctx, input.Email)
	if err != nil {
		if err == database.ErrNotFound {
			span.RecordError(err)
			span.SetStatus(codes.Error, fmt.Sprintf("member not found: %s", input.Email))
			return nil, ErrInvalidCredentials
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("failed to find member by email: %s", input.Email))
		return nil, err
	}

	if !member.ValidatePassword(input.Password) {
		span.SetStatus(codes.Error, "invalid password")
		return nil, ErrInvalidCredentials
	}

	token, err := uc.jwtService.GenerateToken(member.Id, 24*time.Hour)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to generate token")
		return nil, err
	}

	refreshToken, err := uc.jwtService.GenerateToken(member.Id, 7*24*time.Hour)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to generate refresh token")
		return nil, err
	}

	return &SignInUseCaseOutput{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    int64((24 * time.Hour).Seconds()),
	}, nil
}
