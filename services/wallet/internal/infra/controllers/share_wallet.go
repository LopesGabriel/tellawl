package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

type shareWalletRequest struct {
	UserEmail string `json:"user_email"`
}

func (handler *APIHandler) HandleShareWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", handler.version)
	w.Header().Add("Content-Type", "application/json")

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	creatorId, err := claims.GetSubject()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get token subject",
			"error":   err.Error(),
		})
		return
	}

	var data shareWalletRequest
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

	wallet, err := handler.usecases.ShareWallet(r.Context(), usecases.ShareWalletUseCaseInput{
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
