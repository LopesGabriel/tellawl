package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/application/use-cases"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type signInHttpHandler struct {
	userRepository repository.UserRepository
	version        string
}

func NewSignInHttpHandler(userRepository repository.UserRepository) *signInHttpHandler {
	return &signInHttpHandler{
		userRepository: userRepository,
		version:        "1.0.0",
	}
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *signInHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	useCase := usecases.NewAuthenticateUserUseCase(c.userRepository, jwtSecret)

	var data createUserRequest
	// Read the requst body
	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := useCase.Execute(usecases.AuthenticateUserUseCaseInput{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if errors.Is(err, usecases.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
