package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/application/use-cases"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers/presenter"
)

type CreateWalletRequest struct {
	WalletName string `json:"wallet_name"`
}

type createWalletHttpHandler struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
	version          string
}

func NewCreateWalletHttpHandler(userRepository repository.UserRepository, walletRepository repository.WalletRepository) *createWalletHttpHandler {
	return &createWalletHttpHandler{
		userRepository:   userRepository,
		walletRepository: walletRepository,
		version:          "1.0.0",
	}
}

func (c *createWalletHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	useCase := usecases.NewCreateWalletUseCase(c.userRepository, c.walletRepository)

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	creatorId, err := claims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data CreateWalletRequest
	// Read the requst body
	err = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wallet, err := useCase.Execute(usecases.CreateWalletUseCaseInput{
		CreatorID: creatorId,
		Name:      data.WalletName,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	httpWallet := presenter.NewHTTPWallet(*wallet)

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Location", "/wallets/"+wallet.Id)
	w.WriteHeader(http.StatusCreated)
	w.Write(httpWallet.ToJSON())
}
