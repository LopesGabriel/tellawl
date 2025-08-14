package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	router := setupHttpServer()
	http.ListenAndServe(":8080", router)
}

func setupHttpServer() *mux.Router {
	router := mux.NewRouter()

	publisher := events.InMemoryEventPublisher{}
	userRepository := database.NewInMemoryUserRepository(publisher)
	walletRepository := database.NewInMemoryWalletRepository(publisher)

	healthHandler := controllers.NewHealthHttpHandler()
	createUserHandler := controllers.NewCreateUserHttpHandler(userRepository)
	signInHander := controllers.NewSignInHttpHandler(userRepository)
	createWalletHandler := controllers.NewCreateWalletHttpHandler(
		userRepository,
		walletRepository,
	)

	router.Handle("/health", healthHandler).Methods("GET")
	router.Handle("/sign-up", createUserHandler).Methods("POST")
	router.Handle("/sign-in", signInHander).Methods("POST")

	// Authenticated routes
	router.Handle("/wallets", controllers.JWTAuthMiddleware(createWalletHandler))

	return router
}
