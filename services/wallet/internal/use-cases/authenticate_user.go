package usecases

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthenticateUserUseCaseInput struct {
	Email    string
	Password string
}

func (usecase *UseCase) AuthenticateUser(ctx context.Context, input AuthenticateUserUseCaseInput) (string, error) {
	if input.Email == "" || input.Password == "" {
		return "", ErrInvalidInput
	}

	user, err := usecase.repos.User.FindByEmail(ctx, input.Email)
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

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
