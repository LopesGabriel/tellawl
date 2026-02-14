package usecases

import (
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	"go.opentelemetry.io/otel/trace"
)

type UseCase struct {
	repos  *repository.Repositories
	tracer trace.Tracer
	logger *logger.AppLogger
}

type NewUseCasesArgs struct {
	Logger *logger.AppLogger
	Repos  *repository.Repositories
	Tracer trace.Tracer
}

func NewUseCases(args NewUseCasesArgs) *UseCase {
	return &UseCase{
		repos:  args.Repos,
		tracer: args.Tracer,
		logger: args.Logger,
	}
}
