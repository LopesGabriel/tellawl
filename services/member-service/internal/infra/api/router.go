package api

import (
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func (h *apiHandler) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.Use(h.requestInterceptorMiddleware)
	router.Use(otelmux.Middleware("member-service"))

	// Health check route
	router.HandleFunc("/internal/health", h.HandleHealthCheck).Methods("GET")

	// Authentication routes
	router.HandleFunc("/public/signup", h.ErrorMiddleware(h.HandleSignUp)).Methods("POST")
	router.HandleFunc("/public/signin", h.ErrorMiddleware(h.HandleSignIn)).Methods("POST")
	router.HandleFunc("/public/refresh-token", h.ErrorMiddleware(h.HandleRefreshToken)).Methods("POST")

	router.HandleFunc("/public/me", h.ErrorMiddleware(h.HandleMe)).Methods("GET")

	// Internal routes
	router.HandleFunc("/internal/members", h.ErrorMiddleware(h.HandleListMembers)).Methods("GET")
	router.HandleFunc("/internal/members/{id}", h.ErrorMiddleware(h.HandleGetMemberByID)).Methods("GET")

	return router
}
