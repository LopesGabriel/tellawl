package usecases

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type authenticateUserUseCase struct {
	userRepository repository.UserRepository
	jwtSecret      string
}

func NewAuthenticateUserUseCase(userRepository repository.UserRepository, jwtSecret string) *authenticateUserUseCase {
	return &authenticateUserUseCase{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

type AuthenticateUserUseCaseInput struct {
	Email    string
	Password string
}

func (usecase *authenticateUserUseCase) Execute(input AuthenticateUserUseCaseInput) (string, error) {
	if input.Email == "" || input.Password == "" {
		return "", ErrInvalidInput
	}

	user, err := usecase.userRepository.FindByEmail(input.Email)
	if err != nil {
		return "", err
	}

	if !user.ValidatePassword(input.Password) {
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"iss": "com.tellaw.bank",
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(usecase.jwtSecret))
}
