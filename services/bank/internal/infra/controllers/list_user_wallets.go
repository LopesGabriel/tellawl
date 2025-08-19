package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/application/use-cases"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers/presenter"
)

type listUserWalletsHttpHandler struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
	version          string
}

func NewListUserWalletsHttpHandler(userRepository repository.UserRepository, walletRepository repository.WalletRepository, version string) *listUserWalletsHttpHandler {
	return &listUserWalletsHttpHandler{
		userRepository:   userRepository,
		walletRepository: walletRepository,
		version:          version,
	}
}

func (c *listUserWalletsHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.Header().Add("Content-Type", "application/json")

	useCase := usecases.NewListUserWalletsUseCase(c.userRepository, c.walletRepository)

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	userId, err := claims.GetSubject()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get token subject",
			"error":   err.Error(),
		})
		return
	}

	wallets, err := useCase.Execute(usecases.ListUserWalletsUseCaseInput{
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
		httpWallets[i] = presenter.NewHTTPWallet(*wallet)
	}

	jsonData, _ := json.Marshal(httpWallets)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
