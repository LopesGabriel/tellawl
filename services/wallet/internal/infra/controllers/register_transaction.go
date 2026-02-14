package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/wallet/internal/use-cases"
	"go.opentelemetry.io/otel/attribute"
)

type registerTransactionRequest struct {
	Amount          int    `json:"amount"`
	Offset          int    `json:"offset"`
	TransactionType string `json:"transaction_type"`
	Description     string `json:"description"`
}

func (h *APIHandler) HandleRegisterTransaction(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HandleRegisterTransaction")
	defer span.End()

	member := r.Context().Value(memberContextKey).(*models.Member)
	creatorId := member.Id

	var data registerTransactionRequest
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

	span.SetAttributes(
		attribute.String("wallet_id", walletId),
		attribute.String("user_id", creatorId),
	)
	transaction, err := h.usecases.RegisterTransaction(ctx, usecases.RegisterTransactionUseCaseInput{
		TransactionRegisteredByUserId: creatorId,
		WalletId:                      walletId,
		Amount:                        data.Amount,
		Offset:                        data.Offset,
		TransactionType:               data.TransactionType,
		Description:                   data.Description,
	})
	if err != nil {
		h.logger.Error(ctx, "Could not register transaction", slog.String("error", err.Error()))
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not register the transaction",
			"error":   err.Error(),
		})
		return
	}

	httpTransaction := presenter.NewHTTPTransaction(*transaction)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Location", "/wallets/"+walletId+"/transactions/"+transaction.Id)
	w.WriteHeader(http.StatusCreated)
	w.Write(httpTransaction.ToJSON())
}
