package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/database"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/events"
)

const VERSION = "1.0.0"

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

	healthHandler := controllers.NewHealthHttpHandler(VERSION)
	signUpHandler := controllers.NewSignUpHttpHandler(userRepository, VERSION)
	signInHander := controllers.NewSignInHttpHandler(userRepository, VERSION)
	shareWalletHandler := controllers.NewShareWalletHttpHandler(
		userRepository,
		walletRepository,
		VERSION,
	)
	createWalletHandler := controllers.NewCreateWalletHttpHandler(
		userRepository,
		walletRepository,
		VERSION,
	)
	registerTransactionHandler := controllers.NewRegisterTransactionHttpHandler(
		userRepository,
		walletRepository,
		VERSION,
	)
	listUserWallets := controllers.NewListUserWalletsHttpHandler(
		userRepository,
		walletRepository,
		VERSION,
	)
	createCategoryHandler := controllers.NewCreateCategoryHttpHandler(
		walletRepository,
		VERSION,
	)
	listCategoryHandler := controllers.NewListCategoryHttpHandler(
		walletRepository,
		VERSION,
	)

	router.Handle("/health", healthHandler).Methods("GET")
	router.Handle("/sign-up", signUpHandler).Methods("POST")
	router.Handle("/sign-in", signInHander).Methods("POST")

	// Authenticated routes
	router.Handle("/wallets", controllers.JWTAuthMiddleware(createWalletHandler)).Methods("POST")
	router.Handle("/wallets", controllers.JWTAuthMiddleware(listUserWallets)).Methods("GET")
	router.Handle("/wallets/{wallet_id}/share", controllers.JWTAuthMiddleware(shareWalletHandler)).Methods("POST")

	// Transactions
	router.Handle("/wallets/{wallet_id}/transactions", controllers.JWTAuthMiddleware(registerTransactionHandler)).Methods("POST")

	// Categories
	router.Handle("/wallets/{wallet_id}/categories", controllers.JWTAuthMiddleware(createCategoryHandler)).Methods("POST")
	router.Handle("/wallets/{wallet_id}/categories", controllers.JWTAuthMiddleware(listCategoryHandler)).Methods("GET")

	return router
}
