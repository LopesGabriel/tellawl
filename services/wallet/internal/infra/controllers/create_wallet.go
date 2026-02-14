package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

type createWalletRequest struct {
	WalletName string `json:"name"`
}

func (h *APIHandler) HandleCreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HandleCreateWallet")
	defer span.End()

	member := r.Context().Value(memberContextKey).(*models.Member)
	creatorId := member.Id

	var data createWalletRequest
	// Read the requst body
	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		h.logger.Error(ctx, "Could not decode the request body", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not parse the request body, are you sending a JSON?",
			"error":   err.Error(),
		})
		return
	}

	wallet, err := h.usecases.CreateWallet(ctx, usecases.CreateWalletUseCaseInput{
		CreatorID: creatorId,
		Creator:   member,
		Name:      data.WalletName,
	})
	if err != nil {
		h.logger.Error(ctx, "Could not create the wallet", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not create the wallet",
			"error":   err.Error(),
		})
		return
	}

	httpWallet := presenter.NewHTTPWallet(*wallet)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Location", "/wallets/"+wallet.Id)
	w.WriteHeader(http.StatusCreated)
	w.Write(httpWallet.ToJSON())
}
