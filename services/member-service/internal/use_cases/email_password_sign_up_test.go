package usecases

import (
	"testing"

	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	inmemory "github.com/lopesgabriel/tellawl/services/member-service/internal/infra/events/in_memory"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestEmailPasswordSignUpUseCase(t *testing.T) {
	t.Run("Sign up member with email password", func(t *testing.T) {
		publisher := inmemory.InitInMemoryEventPublisher()
		repo := repository.NewInMemory(publisher)
		usecases := InitUseCases(InitUseCasesArgs{
			JwtSecret: "t3st-S3cret",
			Repos:     repo,
			Tracer:    noop.NewTracerProvider().Tracer("test"),
		})

		firstName := "John"
		lastName := "Doe"
		email := "john.doe@example.com"
		password := "password123"

		member, err := usecases.EmailPasswordSignUp(t.Context(), EmailPasswordSignUpUseCaseInput{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			Password:  password,
		})
		if err != nil {
			t.Errorf("Error signing up member: %v", err)
		}
		if member == nil {
			t.Error("Expected member to be created")
			return
		}
		if member.FirstName != firstName {
			t.Errorf("Expected first name to be %s, got %s", firstName, member.FirstName)
		}
		if member.LastName != lastName {
			t.Errorf("Expected last name to be %s, got %s", lastName, member.LastName)
		}
		if member.Email != email {
			t.Errorf("Expected email to be %s, got %s", email, member.Email)
		}
		if member.HashedPassword == password {
			t.Errorf("Expected password to be hashed")
		}
		if !member.ValidatePassword(password) {
			t.Errorf("Expected password to be valid")
		}
	})

	t.Run("Should not sign up when an account with that email already exists", func(t *testing.T) {
		publisher := inmemory.InitInMemoryEventPublisher()
		repo := repository.NewInMemory(publisher)
		usecases := InitUseCases(InitUseCasesArgs{
			JwtSecret: "t3st-S3cret",
			Repos:     repo,
			Tracer:    noop.NewTracerProvider().Tracer("test"),
		})

		repo.Members.Save(t.Context(), &models.Member{Email: "john.doe@example.com"})

		_, err := usecases.EmailPasswordSignUp(t.Context(), EmailPasswordSignUpUseCaseInput{
			Email:     "john.doe@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "XXXXXXXXXXX",
		})
		if err != ErrMemberAlreadyExists {
			t.Errorf("Expected error to be %v, got %v", ErrMemberAlreadyExists, err)
		}
	})
}
