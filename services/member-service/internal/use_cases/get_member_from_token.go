package usecases

import (
	"context"
	"log/slog"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
	"go.opentelemetry.io/otel/codes"
)

type GetMemberFromTokenUseCaseInput struct {
	Token string
}

func (uc *UseCases) GetMemberFromToken(ctx context.Context, input GetMemberFromTokenUseCaseInput) (*models.Member, error) {
	ctx, span := uc.tracer.Start(ctx, "GetMemberFromTokenUseCase")
	defer span.End()

	memberId, err := uc.jwtService.ValidateToken(input.Token)
	if err != nil {
		uc.logger.Error(ctx, "failed to validate token", slog.String("error", err.Error()))
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate token")
		return nil, ErrInvalidCredentials
	}

	member, err := uc.repos.Members.FindByID(ctx, memberId)
	if err != nil {
		uc.logger.Error(ctx, "failed to find member", slog.String("error", err.Error()))
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to find member")
		return nil, ErrMemberNotFound
	}

	return member, nil
}
