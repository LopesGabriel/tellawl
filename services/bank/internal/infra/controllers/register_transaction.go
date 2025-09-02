package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

type RegisterTransactionRequest struct {
	Amount          int    `json:"amount"`
	CategoryId      string `json:"category_id"`
	Offset          int    `json:"offset"`
	TransactionType string `json:"transaction_type"`
	Description     string `json:"description"`
}

type registerTransactionHttpHandler struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
	version          string
}

func NewRegisterTransactionHttpHandler(userRepository repository.UserRepository, walletRepository repository.WalletRepository, version string) *registerTransactionHttpHandler {
	return &registerTransactionHttpHandler{
		userRepository:   userRepository,
		walletRepository: walletRepository,
		version:          version,
	}
}

func (c *registerTransactionHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.Header().Add("Content-Type", "application/json")

	useCase := usecases.NewRegisterTransactionUseCase(c.userRepository, c.walletRepository)

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	creatorId, err := claims.GetSubject()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get token subject",
			"error":   err.Error(),
		})
		return
	}

	var data RegisterTransactionRequest
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

	transaction, err := useCase.Execute(usecases.RegisterTransactionUseCaseInput{
		TransactionRegisteredByUserId: creatorId,
		WalletId:                      walletId,
		Amount:                        data.Amount,
		CategoryId:                    data.CategoryId,
		Offset:                        data.Offset,
		TransactionType:               data.TransactionType,
		Description:                   data.Description,
	})

	if err != nil {
		WriteError(w, http.StatusBadRequest, map[string]any{
			"message": "Could not register the transaction",
			"error":   err.Error(),
		})
		return
	}

	httpTransaction := presenter.NewHTTPTransaction(*transaction)

	w.Header().Add("Location", "/wallets/"+walletId+"/transactions/"+transaction.Id)
	w.WriteHeader(http.StatusCreated)
	w.Write(httpTransaction.ToJSON())
}
