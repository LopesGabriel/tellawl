package usecases_test

import (
	"testing"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	inmemoryevt "github.com/lopesgabriel/tellawl/services/member-service/internal/infra/events/in_memory"
	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	noopl "go.opentelemetry.io/otel/log/noop"
	noopt "go.opentelemetry.io/otel/trace/noop"
)

func TestSignInUseCase(t *testing.T) {
	t.Run("Success Sign In with email and password", func(t *testing.T) {
		publisher := inmemoryevt.InitInMemoryEventPublisher()
		repo := repository.NewInMemory(publisher)
		logger, err := logger.Init(t.Context(), logger.InitLoggerArgs{
			LoggerProvider:   noopl.NewLoggerProvider(),
			ServiceNamespace: "test",
		})
		if err != nil {
			t.Fatal(err)
		}

		ucs := usecases.InitUseCases(usecases.InitUseCasesArgs{
			JwtSecret: "t3st-S3cret",
			Repos:     repo,
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

		result, err := ucs.SignIn(t.Context(), usecases.SignInUseCaseInput{
			Email:    "john.doe@example.com",
			Password: "XXXX123XXXTest",
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(result.Token) == 0 {
			t.Fatal("Token is empty")
		}

		if len(result.RefreshToken) == 0 {
			t.Fatal("RefreshToken is empty")
		}

		if result.ExpiresIn == 0 {
			t.Fatal("ExpiresIn is empty")
		}
	})
}
