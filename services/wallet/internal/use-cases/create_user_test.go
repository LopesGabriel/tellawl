package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
	usecase "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

func TestUserCreation(t *testing.T) {
	publisher := events.InMemoryEventPublisher{}
	repos := repository.NewInMemory(publisher)
	usecases := usecase.NewUseCases(usecase.NewUseCasesArgs{
		JwtSecret: "example",
		Repos:     repos,
	})

	user, err := usecases.CreateUser(usecase.CreateUserUseCaseInput{
		FirstName: "Gabriel",
		LastName:  "Lopes",
		Email:     "lopesgabriel@example.com",
		Password:  "S4mpl3P4ssW0rd",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.FirstName != "Gabriel" {
		t.Errorf("Expected first name to be 'Gabriel', got %v", user.FirstName)
	}

	if user.LastName != "Lopes" {
		t.Errorf("Expected last name to be 'Lopes', got %v", user.LastName)
	}

	if user.Email != "lopesgabriel@example.com" {
		t.Errorf("Expected email to be 'lopesgabriel@example.com', got %v", user.Email)
	}

	if user.HashedPassword == "S4mpl3P4ssW0rd" {
		t.Errorf("Expected password to be hashed, got %v", user.HashedPassword)
	}

	if !user.ValidatePassword("S4mpl3P4ssW0rd") {
		t.Errorf("Expected password to be valid, got invalid")
	}
}
