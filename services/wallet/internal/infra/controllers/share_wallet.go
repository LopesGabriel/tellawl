package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/packages/logger"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/errx"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
)

type shareWalletRequest struct {
	UserEmail string `json:"user_email"`
}

func (handler *APIHandler) HandleShareWallet(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tracer.Start(r.Context(), "HandleShareWallet")
	defer span.End()

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	creatorId, err := claims.GetSubject()
	if err != nil {
		logger.Error(ctx, "Could not get token subject", slog.String("error", err.Error()))
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
		logger.Error(ctx, "Could not decode the request body", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not parse the request body, are you sending a JSON?",
			"error":   err.Error(),
		})
		return
	}

	vars := mux.Vars(r)
	walletId := vars["wallet_id"]
	if walletId == "" {
		logger.Error(ctx, "Could not get wallet id from path")
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get wallet id from path",
		})
		return
	}

	wallet, err := handler.usecases.ShareWallet(ctx, usecases.ShareWalletUseCaseInput{
		WalletCreatorId: creatorId,
		WalletId:        walletId,
		SharedUserEmail: data.UserEmail,
	})
	if err != nil {
		if errors.Is(err, errx.ErrInsufficientPermissions) {
			logger.Error(ctx, "Insufficient permissions", slog.String("error", err.Error()))
			WriteError(w, http.StatusForbidden, map[string]any{
				"message": "You do not have permission to share this wallet",
			})
			return
		}

		logger.Error(ctx, "Could not share the wallet", slog.String("error", err.Error()))
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
