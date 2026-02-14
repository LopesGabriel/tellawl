package usecases

import (
	"context"
	"errors"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"go.opentelemetry.io/otel/codes"
)

type CreateWalletUseCaseInput struct {
	CreatorID string
	Creator   *models.Member
	Name      string
}

func (usecase *UseCase) CreateWallet(ctx context.Context, input CreateWalletUseCaseInput) (*models.Wallet, error) {
	ctx, span := usecase.tracer.Start(ctx, "CreateWallet")
	defer span.End()

	if input.CreatorID == "" {
		span.AddEvent("missing required field: CreatorID")
		return nil, errx.MissingRequiredFieldsError("CreatorID")
	}

	if input.Name == "" {
		span.AddEvent("missing required field: Name")
		return nil, errx.MissingRequiredFieldsError("Name")
	}

	var creator *models.Member
	if input.Creator != nil {
		creator = input.Creator
	} else {
		member, err := usecase.repos.Member.FindByID(ctx, input.CreatorID)
		if err != nil {
			span.SetStatus(codes.Error, "could not find creator member")
			span.RecordError(err)
			return nil, errors.Join(errors.New("could not find creator id"), err)
		}
		creator = member
	}

	wallet := models.CreateNewWallet(input.Name, creator)

	if err := usecase.repos.Wallet.Save(ctx, wallet); err != nil {
		span.SetStatus(codes.Error, "could not persist the wallet")
		span.RecordError(err)
		return nil, errors.Join(errors.New("could not persist the wallet"), err)
	}

	span.SetStatus(codes.Ok, "wallet created")
	return wallet, nil
}
