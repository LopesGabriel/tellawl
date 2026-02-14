package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

type shareWalletRequest struct {
	UserEmail string `json:"user_email"`
}

func (h *APIHandler) HandleShareWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HandleShareWallet")
	defer span.End()

	member := r.Context().Value(memberContextKey).(*models.Member)
	creatorId := member.Id

	var data shareWalletRequest
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

	vars := mux.Vars(r)
	walletId := vars["wallet_id"]
	if walletId == "" {
		h.logger.Error(ctx, "Could not get wallet id from path")
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get wallet id from path",
		})
		return
	}

	wallet, err := h.usecases.ShareWallet(ctx, usecases.ShareWalletUseCaseInput{
		WalletCreatorId: creatorId,
		WalletId:        walletId,
		SharedUserEmail: data.UserEmail,
	})
	if err != nil {
		if errors.Is(err, errx.ErrInsufficientPermissions) {
			h.logger.Error(ctx, "Insufficient permissions", slog.String("error", err.Error()))
			WriteError(w, http.StatusForbidden, map[string]any{
				"message": "You do not have permission to share this wallet",
			})
			return
		}

		h.logger.Error(ctx, "Could not share the wallet", slog.String("error", err.Error()))
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not share the wallet",
			"error":   err.Error(),
		})
		return
	}

	httpWallet := presenter.NewHTTPWallet(*wallet)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(httpWallet.ToJSON())
}
