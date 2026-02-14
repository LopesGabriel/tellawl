package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	"go.opentelemetry.io/otel/attribute"
)

func (h *APIHandler) HandleListUserWallets(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HandleListUserWallets")
	defer span.End()

	member := r.Context().Value(memberContextKey).(*models.Member)

	span.SetAttributes(attribute.String("member.id", member.Id))
	wallets, err := h.usecases.ListUserWallets(ctx, usecases.ListUserWalletsUseCaseInput{
		UserId: member.Id,
	})
	if err != nil {
		h.logger.Error(ctx, "Could not list user wallets", slog.String("error", err.Error()))
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
