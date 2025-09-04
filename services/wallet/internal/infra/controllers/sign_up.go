package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

type signUpRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (handler *APIHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", handler.version)
	w.Header().Add("Content-Type", "application/json")

	var data signUpRequest
	// Read the requst body
	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not parse the request body, are you sending a JSON?",
			"error":   err.Error(),
		})
		return
	}

	user, err := handler.usecases.CreateUser(r.Context(), usecases.CreateUserUseCaseInput{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	})

	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not sign up",
			"error":   err.Error(),
		})
		return
	}

	httpUser := presenter.NewHTTPUser(*user)

	token, err := handler.usecases.AuthenticateUser(r.Context(), usecases.AuthenticateUserUseCaseInput{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not generate the token",
			"error":   err.Error(),
		})
		return
	}

	w.Header().Add("Location", "/users/"+httpUser.Id)
	w.Header().Add("token", token)
	w.WriteHeader(http.StatusCreated)
	w.Write(httpUser.ToJSON())
}
