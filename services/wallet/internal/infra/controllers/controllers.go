package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/packages/tracing"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	"go.opentelemetry.io/otel/trace"
)

type APIHandler struct {
	usecases *usecases.UseCase
	version  string
	tracer   trace.Tracer
}

func NewAPIHandler(usecases *usecases.UseCase, version string) *APIHandler {
	return &APIHandler{
		usecases: usecases,
		version:  version,
		tracer:   tracing.GetTracer("wallet"),
	}
}

func (handler *APIHandler) registerEndpoints() *mux.Router {
	router := mux.NewRouter()

	router.Use(handler.requestInterceptorMiddleware)

	router.HandleFunc("/health", handler.HandleHealthCheck).Methods("GET")

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
		logger.Debug(r.Context(), "processing new request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}
