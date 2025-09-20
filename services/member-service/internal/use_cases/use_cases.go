package usecases

import (
	"fmt"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	"go.opentelemetry.io/otel/trace"
)

type useCase struct {
	jwtSecret string
	repos     *repository.Repositories
	tracer    trace.Tracer
}

type InitUseCasesArgs struct {
	JwtSecret string
	Repos     *repository.Repositories
	Tracer    trace.Tracer
}

func InitUseCases(args InitUseCasesArgs) *useCase {
	return &useCase{
		jwtSecret: args.JwtSecret,
		repos:     args.Repos,
		tracer:    args.Tracer,
	}
}

var (
	ErrMemberNotFound      = fmt.Errorf("member not found")
	ErrMemberAlreadyExists = fmt.Errorf("member already exists")
)
