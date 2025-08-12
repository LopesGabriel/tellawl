package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers"
)

func main() {
	router := setupHttpServer()
	http.ListenAndServe(":8080", router)
}

func setupHttpServer() *mux.Router {
	router := mux.NewRouter()
	controller := controllers.NewController()

	router.HandleFunc("/health", controller.HandleHealth).Methods("GET")

	return router
}
