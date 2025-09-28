package usecases

import (
	"fmt"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	"go.opentelemetry.io/otel/trace"
)

type UseCases struct {
	jwtSecret string
	repos     *repository.Repositories
	tracer    trace.Tracer
}

type InitUseCasesArgs struct {
	JwtSecret string
	Repos     *repository.Repositories
	Tracer    trace.Tracer
}

func InitUseCases(args InitUseCasesArgs) *UseCases {
	return &UseCases{
		jwtSecret: args.JwtSecret,
		repos:     args.Repos,
		tracer:    args.Tracer,
	}
}

var (
	ErrMemberNotFound      = fmt.Errorf("member not found")
	ErrMemberAlreadyExists = fmt.Errorf("member already exists")
)
