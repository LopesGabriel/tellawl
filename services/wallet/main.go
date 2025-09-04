package main

import (
	"fmt"
	"log/slog"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

const VERSION = "1.0.0"
const JWT_SECRET = "ex4mpl3-Secr3t"
const PORT = 8080

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	publisher := events.InMemoryEventPublisher{}

	slog.Info("Starting Wallet Service", slog.String("version", VERSION))
	repos := repository.NewInMemory(publisher)
	useCases := usecases.NewUseCases(usecases.NewUseCasesArgs{
		JwtSecret: JWT_SECRET,
		Repos:     repos,
	})
	apiHandler := controllers.NewAPIHandler(useCases, VERSION)

	slog.Info(fmt.Sprintf("API Server listenig on port %d", PORT))
	apiHandler.Listen(PORT)
}
