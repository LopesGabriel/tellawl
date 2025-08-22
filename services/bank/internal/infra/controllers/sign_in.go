package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

type signInHttpHandler struct {
	userRepository repository.UserRepository
	version        string
}

func NewSignInHttpHandler(userRepository repository.UserRepository, version string) *signInHttpHandler {
	return &signInHttpHandler{
		userRepository: userRepository,
		version:        version,
	}
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *signInHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.Header().Add("Content-Type", "application/json")

	useCase := usecases.NewAuthenticateUserUseCase(c.userRepository, jwtSecret)

	var data createUserRequest
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

	token, err := useCase.Execute(usecases.AuthenticateUserUseCaseInput{
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

		if errors.Is(err, repository.ErrUserNotFound) {
			WriteError(w, http.StatusUnauthorized, map[string]any{
				"message": "Invalid credentials",
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
