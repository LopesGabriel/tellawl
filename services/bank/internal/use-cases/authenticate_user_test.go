package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

func TestAuthenticateUser(t *testing.T) {
	publisher := events.InMemoryEventPublisher{}
	userRepository := database.NewInMemoryUserRepository(publisher)
	sut := usecases.NewAuthenticateUserUseCase(userRepository, "TestS3cret")

	password := "S4mpl3P4ssW0rd"
	user, _ := models.CreateNewUser("Gabriel", "Lopes", "example@example.com", password)
	userRepository.Save(user)

	token, err := sut.Execute(usecases.AuthenticateUserUseCaseInput{
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
