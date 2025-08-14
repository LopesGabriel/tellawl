package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/application/use-cases"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers/presenter"
)

type ShareWalletRequest struct {
	UserEmail string `json:"user_email"`
}

type shareWalletHttpHandler struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
	version          string
}

func NewShareWalletHttpHandler(userRepository repository.UserRepository, walletRepository repository.WalletRepository, version string) *shareWalletHttpHandler {
	return &shareWalletHttpHandler{
		userRepository:   userRepository,
		walletRepository: walletRepository,
		version:          version,
	}
}

func (c *shareWalletHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.Header().Add("Content-Type", "application/json")

	useCase := usecases.NewShareWalletUseCase(c.userRepository, c.walletRepository)

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	creatorId, err := claims.GetSubject()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get token subject",
			"error":   err.Error(),
		})
		return
	}

	var data ShareWalletRequest
	// Read the requst body
	err = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not parse the request body, are you sending a JSON?",
			"error":   err.Error(),
		})
		return
	}

	vars := mux.Vars(r)
	walletId := vars["wallet_id"]
	if walletId == "" {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get wallet id from path",
		})
		return
	}

	wallet, err := useCase.Execute(usecases.ShareWalletUseCaseInput{
		WalletCreatorId: creatorId,
		WalletId:        walletId,
		SharedUserEmail: data.UserEmail,
	})

	if err != nil {
		if errors.Is(err, usecases.ErrInsufficientPermissions) {
			WriteError(w, http.StatusForbidden, map[string]any{
				"message": "You do not have permission to share this wallet",
			})
			return
		}

		if errors.Is(err, repository.ErrWalletNotFound) {
			WriteError(w, http.StatusNotFound, map[string]any{
				"message": "Wallet not found",
			})
			return
		}

		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not share the wallet",
			"error":   err.Error(),
		})
		return
	}

	httpWallet := presenter.NewHTTPWallet(*wallet)

	w.WriteHeader(http.StatusOK)
	w.Write(httpWallet.ToJSON())
}
