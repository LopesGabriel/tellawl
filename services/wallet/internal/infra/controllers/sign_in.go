package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler *APIHandler) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", handler.version)
	w.Header().Add("Content-Type", "application/json")

	var data signInRequest
	// Read the requst body
	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not parse request body, did you provided a valid JSON?",
			"error":   err.Error(),
		})
		return
	}

	token, err := handler.usecases.AuthenticateUser(usecases.AuthenticateUserUseCaseInput{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Invalid Credentials",
			})
			return
		}

		if errors.Is(err, usecases.ErrInvalidInput) {
			WriteError(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid Input: email and password are required",
				"error":   err.Error(),
			})
			return
		}

		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Something unexpected happened",
			"error":   err.Error(),
		})
		return
	}

	result, _ := json.Marshal(map[string]string{
		"token": token,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
