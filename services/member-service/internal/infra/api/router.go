package api

import (
	"github.com/gorilla/mux"
)

func (h *apiHandler) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.Use(h.requestInterceptorMiddleware)

	router.HandleFunc("/health", h.HandleHealthCheck).Methods("GET")

	// Use error middleware for all routes
	router.HandleFunc("/signup", h.ErrorMiddleware(h.HandleSignUp)).Methods("POST")

	return router
}
