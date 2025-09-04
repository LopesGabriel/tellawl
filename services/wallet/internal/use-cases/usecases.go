package usecases

import (
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type UseCase struct {
	jwtSecret string
	repos     *repository.Repositories
}

type NewUseCasesArgs struct {
	JwtSecret string
	Repos     *repository.Repositories
}

func NewUseCases(args NewUseCasesArgs) *UseCase {
	return &UseCase{
		jwtSecret: args.JwtSecret,
		repos:     args.Repos,
	}
}
