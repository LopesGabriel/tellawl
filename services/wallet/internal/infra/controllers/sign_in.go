package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler *APIHandler) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tracer.Start(r.Context(), "HandleSignIn")
	defer span.End()

	var data signInRequest
	// Read the requst body
	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		logger.Error(ctx, "Could not decode the request body", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not parse request body, did you provided a valid JSON?",
			"error":   err.Error(),
		})
		return
	}

	token, err := handler.usecases.AuthenticateUser(ctx, usecases.AuthenticateUserUseCaseInput{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		if errors.Is(err, errx.ErrInvalidCredentials) {
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Invalid Credentials",
			})
			return
		}

		if errors.Is(err, errx.ErrInvalidInput) {
			WriteError(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid Input: email and password are required",
				"error":   err.Error(),
			})
			return
		}

		if errors.Is(err, errx.ErrNotFound) {
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Unauthorized",
			})
			return
		}

		logger.Error(ctx, "Could not authenticate the user", slog.String("error", err.Error()))
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Something unexpected happened",
			"error":   err.Error(),
		})
		return
	}

	result, _ := json.Marshal(map[string]string{
		"token": token,
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
