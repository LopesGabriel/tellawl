package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
	"github.com/lopesgabriel/tellawl/services/bank/internal/infra/controllers/presenter"
	usecases "github.com/lopesgabriel/tellawl/services/bank/internal/use-cases"
)

type CreateCategoryRequest struct {
	CategoryName string `json:"name"`
}

type createCategoryHttpHandler struct {
	walletRepository repository.WalletRepository
	version          string
}

func NewCreateCategoryHttpHandler(walletRepository repository.WalletRepository, version string) *createCategoryHttpHandler {
	return &createCategoryHttpHandler{
		walletRepository: walletRepository,
		version:          version,
	}
}

func (c *createCategoryHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.Header().Add("Content-Type", "application/json")

	useCase := usecases.NewCreateCategoryUseCase(c.walletRepository)

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	creatorId, err := claims.GetSubject()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get token subject",
			"error":   err.Error(),
		})
		return
	}

	var data CreateCategoryRequest
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

	wallet, err := c.walletRepository.FindById(walletId)
	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			WriteError(w, http.StatusNotFound, map[string]any{
				"message": "Wallet not found",
				"error":   err.Error(),
			})
			return
		}

		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not find wallet",
			"error":   err.Error(),
		})
		return
	}

	if !wallet.IsUserAllowedToRegisterTransactions(creatorId) {
		WriteError(w, http.StatusForbidden, map[string]any{
			"message": "You are not allowed to create categories in this wallet",
		})
		return
	}

	category, err := useCase.Execute(usecases.CreateCategoryUseCaseInput{
		Name:     data.CategoryName,
		WalletId: walletId,
	})
	if err != nil {
		if strings.Contains(err.Error(), "wallet not found") {
			WriteError(w, http.StatusNotFound, map[string]any{
				"message": "Wallet not found",
				"error":   err.Error(),
			})
			return
		}

		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not create the wallet",
			"error":   err.Error(),
		})
		return
	}

	httpCategory := presenter.NewHTTPCategory(*category)

	w.Header().Add("Location", "/wallets/"+walletId+"/category/"+category.Id)
	w.WriteHeader(http.StatusCreated)
	w.Write(httpCategory.ToJSON())
}

type listCategoriesHttpHandler struct {
	walletRepository repository.WalletRepository
	version          string
}

func NewListCategoryHttpHandler(walletRepository repository.WalletRepository, version string) *listCategoriesHttpHandler {
	return &listCategoriesHttpHandler{
		walletRepository: walletRepository,
		version:          version,
	}
}

func (c *listCategoriesHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.Header().Add("Content-Type", "application/json")

	claims := r.Context().Value(userContextKey).(jwt.MapClaims)
	creatorId, err := claims.GetSubject()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get token subject",
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

	wallet, err := c.walletRepository.FindById(walletId)
	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			WriteError(w, http.StatusNotFound, map[string]any{
				"message": "Wallet not found",
				"error":   err.Error(),
			})
			return
		}

		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not find wallet",
			"error":   err.Error(),
		})
		return
	}

	if !wallet.IsUserAllowedToRegisterTransactions(creatorId) {
		WriteError(w, http.StatusForbidden, map[string]any{
			"message": "You are not allowed to list categories in this wallet",
		})
		return
	}

	var httpCategories []presenter.HTTPCategory
	for _, category := range wallet.Categories {
		httpCategories = append(httpCategories, presenter.NewHTTPCategory(category))
	}

	data, _ := json.Marshal(httpCategories)

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
