package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

func (handler *APIHandler) HandleListUserWallets(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", handler.version)
	w.Header().Add("Content-Type", "application/json")

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	userId, err := claims.GetSubject()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get token subject",
			"error":   err.Error(),
		})
		return
	}

	wallets, err := handler.usecases.ListUserWallets(r.Context(), usecases.ListUserWalletsUseCaseInput{
		UserId: userId,
	})

	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not list user wallets",
			"error":   err.Error(),
		})
		return
	}

	httpWallets := make([]presenter.HTTPWallet, len(wallets))
	for i, wallet := range wallets {
		httpWallets[i] = presenter.NewHTTPWallet(wallet)
	}

	jsonData, _ := json.Marshal(httpWallets)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
