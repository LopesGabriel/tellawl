package usecases_test

import (
	"testing"

	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/application/use-cases"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
)

func TestUserCreation(t *testing.T) {
	publisher := events.InMemoryEventPublisher{}
	userRepository := database.NewInMemoryUserRepository(publisher)
	sut := usecases.NewCreateUserUseCase(userRepository)

	user, err := sut.Execute(usecases.CreateUserUseCaseInput{
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

	if len(userRepository.Items) < 1 {
		t.Errorf("Expected at least one user in the repository, got %v", len(userRepository.Items))
	}

	if !user.ValidatePassword("S4mpl3P4ssW0rd") {
		t.Errorf("Expected password to be valid, got invalid")
	}
}
