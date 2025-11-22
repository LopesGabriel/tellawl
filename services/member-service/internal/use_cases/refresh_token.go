package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/codes"
)

type RefreshTokenUseCaseInput struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenUseCaseOutput struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (uc *UseCases) RefreshToken(ctx context.Context, input RefreshTokenUseCaseInput) (*RefreshTokenUseCaseOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "RefreshTokenUseCase")
	defer span.End()

	memberID, err := uc.jwtService.ValidateToken(input.RefreshToken)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate token")
		return nil, ErrInvalidCredentials
	}

	member, err := uc.repos.Members.FindByID(ctx, memberID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to find member by ID")
		return nil, ErrMemberNotFound
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

	message := fmt.Sprintf("Tokens refreshed for member %s", member.FirstName)
	uc.logger.Info(ctx, message, slog.String("id", member.Id))

	return &RefreshTokenUseCaseOutput{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    int64((24 * time.Hour).Seconds()),
	}, nil
}
