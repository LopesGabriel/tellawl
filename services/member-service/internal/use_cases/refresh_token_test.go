package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/publisher"
	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	noopl "go.opentelemetry.io/otel/log/noop"
	noopt "go.opentelemetry.io/otel/trace/noop"
)

func TestRefreshTokenUseCase(t *testing.T) {
	t.Run("Successfully generates new token from refresh token", func(t *testing.T) {
		publisher := publisher.InitInMemoryEventPublisher()
		memberRepo := database.InitInMemoryMemberRepository(publisher)
		logger, err := logger.Init(t.Context(), logger.InitLoggerArgs{
			LoggerProvider:   noopl.NewLoggerProvider(),
			ServiceNamespace: "test",
		})
		if err != nil {
			t.Fatal(err)
		}

		ucs := usecases.InitUseCases(usecases.InitUseCasesArgs{
			JwtSecret: "t3st-S3cret",
			Repos:     &repository.Repositories{Members: memberRepo},
			Tracer:    noopt.NewTracerProvider().Tracer("test"),
			Logger:    logger,
		})

		_, err = ucs.EmailPasswordSignUp(t.Context(), usecases.EmailPasswordSignUpUseCaseInput{
			Email:     "john.doe@example.com",
			Password:  "XXXX123XXXTest",
			FirstName: "John",
			LastName:  "Doe",
		})
		if err != nil {
			t.Fatal(err)
		}

		signInResult, err := ucs.SignIn(t.Context(), usecases.SignInUseCaseInput{
			Email:    "john.doe@example.com",
			Password: "XXXX123XXXTest",
		})
		if err != nil {
			t.Fatal(err)
		}

		result, err := ucs.RefreshToken(t.Context(), usecases.RefreshTokenUseCaseInput{
			RefreshToken: signInResult.RefreshToken,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(result.Token) == 0 {
			t.Fatal("Token is empty")
		}
		if len(result.RefreshToken) == 0 {
			t.Fatal("Refresh token is empty")
		}
		if result.ExpiresIn == 0 {
			t.Fatal("ExpiresIn is empty")
		}
	})

	t.Run("should fail when sending invalid refresh token", func(t *testing.T) {
		publisher := publisher.InitInMemoryEventPublisher()
		memberRepo := database.InitInMemoryMemberRepository(publisher)
		logger, err := logger.Init(t.Context(), logger.InitLoggerArgs{
			LoggerProvider:   noopl.NewLoggerProvider(),
			ServiceNamespace: "test",
		})
		if err != nil {
			t.Fatal(err)
		}

		ucs := usecases.InitUseCases(usecases.InitUseCasesArgs{
			JwtSecret: "t3st-S3cret",
			Repos:     &repository.Repositories{Members: memberRepo},
			Tracer:    noopt.NewTracerProvider().Tracer("test"),
			Logger:    logger,
		})

		_, err = ucs.RefreshToken(t.Context(), usecases.RefreshTokenUseCaseInput{
			RefreshToken: "Invalid token",
		})
		if err == nil {
			t.Fatal("Expected error")
		}
	})
}
