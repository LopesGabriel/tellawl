package usecases

import (
	"fmt"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/jwt"
	"go.opentelemetry.io/otel/trace"
)

type UseCases struct {
	jwtService *jwt.JWTService
	repos      *repository.Repositories
	tracer     trace.Tracer
	logger     *logger.AppLogger
}

type InitUseCasesArgs struct {
	JwtSecret string
	Repos     *repository.Repositories
	Tracer    trace.Tracer
	Logger    *logger.AppLogger
}

func InitUseCases(args InitUseCasesArgs) *UseCases {
	return &UseCases{
		jwtService: jwt.NewJWTService(args.JwtSecret),
		repos:      args.Repos,
		tracer:     args.Tracer,
		logger:     args.Logger,
	}
}

var (
	ErrMemberNotFound      = fmt.Errorf("member not found")
	ErrMemberAlreadyExists = fmt.Errorf("member already exists")
	ErrInvalidCredentials  = fmt.Errorf("invalid credentials")
)
