package usecases

import (
	"context"
	"fmt"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database"
	"go.opentelemetry.io/otel/codes"
)

type EmailPasswordSignUpUseCaseInput struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

func (uc *UseCases) EmailPasswordSignUp(ctx context.Context, input EmailPasswordSignUpUseCaseInput) (*models.Member, error) {
	ctx, span := uc.tracer.Start(ctx, "EmailPasswordSignUp")
	defer span.End()

	member, err := uc.repos.Members.FindByEmail(ctx, input.Email)
	if err != nil {
		if err != database.ErrNotFound {
			span.RecordError(err)
			span.SetStatus(codes.Error, fmt.Sprintf("failed to find member by email: %s", input.Email))
			return nil, err
		}
	}

	if member != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("member with email %s already exists", input.Email))
		return nil, ErrMemberAlreadyExists
	}

	member, err = models.CreateNewMember(input.FirstName, input.LastName, input.Email, input.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create new member")
		return nil, err
	}

	err = uc.repos.Members.Save(ctx, member)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to persist member")
		return nil, err
	}

	span.SetStatus(codes.Ok, "member created successfully")
	return member, nil
}
