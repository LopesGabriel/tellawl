package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

type APIHandler struct {
	usecases *usecases.UseCase
	version  string
}

func NewAPIHandler(usecases *usecases.UseCase, version string) *APIHandler {
	return &APIHandler{
		usecases: usecases,
		version:  version,
	}
}

func (handler *APIHandler) registerEndpoints() *mux.Router {
	router := mux.NewRouter()

	router.Use(handler.requestInterceptorMiddleware)

	router.HandleFunc("/health", handler.HandleHealthCheck).Methods("GET")
	router.HandleFunc("/sign-up", handler.HandleSignUp).Methods("POST")
	router.HandleFunc("/sign-in", handler.HandleSignIn).Methods("POST")

	// Authenticated routes
	router.Handle("/wallets", jwtAuthMiddleware(
		http.HandlerFunc(handler.HandleCreateWallet))).Methods("POST")
	router.Handle("/wallets", jwtAuthMiddleware(
		http.HandlerFunc(handler.HandleListUserWallets))).Methods("GET")
	router.Handle("/wallets/{wallet_id}/share", jwtAuthMiddleware(
		http.HandlerFunc(handler.HandleShareWallet))).Methods("POST")

	// Transactions
	router.Handle("/wallets/{wallet_id}/transactions", jwtAuthMiddleware(
		http.HandlerFunc(handler.HandleRegisterTransaction))).Methods("POST")

	return router
}

func (handler *APIHandler) Listen(port int) error {
	router := handler.registerEndpoints()
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func WriteError(w http.ResponseWriter, statusCode int, payload map[string]any) {
	result, _ := json.Marshal(payload)

	w.Header().Del("Content-Type")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(result)
}

func (handler *APIHandler) requestInterceptorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Server-Version", handler.version)
		slog.Debug(
			"Received new request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}
