package controllers

import (
	"encoding/json"
	"net/http"

	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/application/use-cases"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers/presenter"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type createUserHttpHandler struct {
	version        string
	userRepository repository.UserRepository
}

func NewCreateUserHttpHandler(userRepository repository.UserRepository) *createUserHttpHandler {
	return &createUserHttpHandler{
		version:        "1.0.0",
		userRepository: userRepository,
	}
}

func (c *createUserHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.Header().Add("Content-Type", "application/json")

	useCase := usecases.NewCreateUserUseCase(c.userRepository)
	signInUseCase := usecases.NewAuthenticateUserUseCase(c.userRepository, jwtSecret)

	var data CreateUserRequest
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

	user, err := useCase.Execute(usecases.CreateUserUseCaseInput{
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

	token, err := signInUseCase.Execute(usecases.AuthenticateUserUseCaseInput{
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
