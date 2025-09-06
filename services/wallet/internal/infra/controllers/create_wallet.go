package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

type createWalletRequest struct {
	WalletName string `json:"name"`
}

func (handler *APIHandler) HandleCreateWallet(w http.ResponseWriter, r *http.Request) {
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

	var data createWalletRequest
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

	wallet, err := handler.usecases.CreateWallet(r.Context(), usecases.CreateWalletUseCaseInput{
		CreatorID: creatorId,
		Name:      data.WalletName,
	})

	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not create the wallet",
			"error":   err.Error(),
		})
		return
	}

	httpWallet := presenter.NewHTTPWallet(*wallet)

	w.Header().Add("Location", "/wallets/"+wallet.Id)
	w.WriteHeader(http.StatusCreated)
	w.Write(httpWallet.ToJSON())
}
