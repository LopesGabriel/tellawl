package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/events"
	usecase "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

func TestAuthenticateUser(t *testing.T) {
	publisher := events.InMemoryEventPublisher{}
	repos := repository.NewInMemory(publisher)
	usecases := usecase.NewUseCases(usecase.NewUseCasesArgs{
		JwtSecret: "example",
		Repos:     repos,
	})

	password := "S4mpl3P4ssW0rd"
	user, _ := models.CreateNewUser("Gabriel", "Lopes", "example@example.com", password)
	repos.User.Save(t.Context(), user)

	token, err := usecases.AuthenticateUser(t.Context(), usecase.AuthenticateUserUseCaseInput{
		Email:    "example@example.com",
		Password: password,
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected token to be defined")
	}
}
