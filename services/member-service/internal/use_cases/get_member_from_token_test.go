package usecases_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/repository"
	inmemoryevt "github.com/lopesgabriel/tellawl/services/member-service/internal/infra/events/in_memory"
	usecases "github.com/lopesgabriel/tellawl/services/member-service/internal/use_cases"
	noopl "go.opentelemetry.io/otel/log/noop"
	noopt "go.opentelemetry.io/otel/trace/noop"
)

func TestGetMemberFromTokenUseCase(t *testing.T) {
	t.Run("successfully retrieves member from token", func(t *testing.T) {
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

		signInResult, err := ucs.SignIn(t.Context(), usecases.SignInUseCaseInput{
			Email:    "john.doe@example.com",
			Password: "XXXX123XXXTest",
		})
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(signInResult.Token)

		member, err := ucs.GetMemberFromToken(t.Context(), usecases.GetMemberFromTokenUseCaseInput{
			Token: signInResult.Token,
		})
		if err != nil {
			t.Fatal(err)
		}

		if member.Email != "john.doe@example.com" {
			t.Fatal("member email is not correct")
		}
		if member.FirstName != "John" {
			t.Fatal("member first name is not correct")
		}
		if member.LastName != "Doe" {
			t.Fatal("member last name is not correct")
		}
	})

	t.Run("should throw invalid credentials error", func(t *testing.T) {
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

		_, err = ucs.GetMemberFromToken(t.Context(), usecases.GetMemberFromTokenUseCaseInput{
			Token: "invalid credentials",
		})
		if err == nil {
			t.Fatal("expected error")
		}

		if !errors.Is(err, usecases.ErrInvalidCredentials) {
			t.Fatal("expected invalid credentials error")
		}
	})

	t.Run("should throw member not found error", func(t *testing.T) {
		publisher := inmemoryevt.InitInMemoryEventPublisher()
		repo := repository.NewInMemory(publisher)
		repo2 := repository.NewInMemory(publisher)
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

		signInResult, err := ucs.SignIn(t.Context(), usecases.SignInUseCaseInput{
			Email:    "john.doe@example.com",
			Password: "XXXX123XXXTest",
		})
		if err != nil {
			t.Fatal(err)
		}

		ucs = usecases.InitUseCases(usecases.InitUseCasesArgs{
			JwtSecret: "t3st-S3cret",
			Repos:     repo2,
			Tracer:    noopt.NewTracerProvider().Tracer("test"),
			Logger:    logger,
		})

		_, err = ucs.GetMemberFromToken(t.Context(), usecases.GetMemberFromTokenUseCaseInput{
			Token: signInResult.Token,
		})
		if err == nil {
			t.Fatal("expected error")
		}

		if !errors.Is(err, usecases.ErrMemberNotFound) {
			t.Fatal("expected member not found error")
		}
	})
}
