package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	"go.opentelemetry.io/otel/attribute"
)

func (handler *APIHandler) HandleListUserWallets(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "HandleListUserWallets")
	defer span.End()

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	userId, err := claims.GetSubject()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get token subject",
			"error":   err.Error(),
		})
		return
	}

	span.SetAttributes(attribute.String("user_id", userId))
	wallets, err := handler.usecases.ListUserWallets(ctx, usecases.ListUserWalletsUseCaseInput{
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
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
