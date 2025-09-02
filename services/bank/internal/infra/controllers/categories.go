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
	"github.com/lopesgabriel/tellawl/services/bank/internal/use-cases/category"
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

	useCase := category.NewCreateCategoryUseCase(c.walletRepository)

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

	category, err := useCase.Execute(category.CreateCategoryUseCaseInput{
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

	vars := mux.Vars(r)
	walletId := vars["wallet_id"]
	if walletId == "" {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get wallet id from path",
		})
		return
	}

	listCategoriesUseCase := category.NewListCategoriesUseCase(c.walletRepository)
	categories, err := listCategoriesUseCase.Execute(category.ListCategoriesUseCaseInput{
		WalletId: walletId,
	})

	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			WriteError(w, http.StatusNotFound, map[string]any{
				"message": "Wallet not found",
				"error":   err.Error(),
			})
			return
		}

		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not list categories",
			"error":   err.Error(),
		})
		return
	}

	var httpCategories []presenter.HTTPCategory
	for _, category := range categories {
		httpCategories = append(httpCategories, presenter.NewHTTPCategory(category))
	}

	data, _ := json.Marshal(httpCategories)

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type UpdateCategoryRequest struct {
	CategoryName string `json:"name"`
}

type getCategoryHttpHandler struct {
	walletRepository repository.WalletRepository
	version          string
}

func NewGetCategoryHttpHandler(walletRepository repository.WalletRepository, version string) *getCategoryHttpHandler {
	return &getCategoryHttpHandler{
		walletRepository: walletRepository,
		version:          version,
	}
}

func (c *getCategoryHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server-Version", c.version)
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	walletId := vars["wallet_id"]
	categoryId := vars["category_id"]

	if walletId == "" || categoryId == "" {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get wallet id or category id from path",
		})
		return
	}

	getCategoryUseCase := category.NewGetCategoryUseCase(c.walletRepository)
	categoryResult, err := getCategoryUseCase.Execute(category.GetCategoryUseCaseInput{
		WalletId:   walletId,
		CategoryId: categoryId,
	})

	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) || strings.Contains(err.Error(), "not found") {
			WriteError(w, http.StatusNotFound, map[string]any{
				"message": "Category not found",
				"error":   err.Error(),
			})
			return
		}

		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get category",
			"error":   err.Error(),
		})
		return
	}

	httpCategory := presenter.NewHTTPCategory(*categoryResult)
	w.WriteHeader(http.StatusOK)
	w.Write(httpCategory.ToJSON())
}

type updateCategoryHttpHandler struct {
	walletRepository repository.WalletRepository
	version          string
}

func NewUpdateCategoryHttpHandler(walletRepository repository.WalletRepository, version string) *updateCategoryHttpHandler {
	return &updateCategoryHttpHandler{
		walletRepository: walletRepository,
		version:          version,
	}
}

func (c *updateCategoryHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var data UpdateCategoryRequest
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
	categoryId := vars["category_id"]

	if walletId == "" || categoryId == "" {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get wallet id or category id from path",
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
			"message": "You are not allowed to update categories in this wallet",
		})
		return
	}

	updateCategoryUseCase := category.NewUpdateCategoryUseCase(c.walletRepository)
	categoryResult, err := updateCategoryUseCase.Execute(category.UpdateCategoryUseCaseInput{
		WalletId:   walletId,
		CategoryId: categoryId,
		Name:       data.CategoryName,
	})

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteError(w, http.StatusNotFound, map[string]any{
				"message": "Category not found",
				"error":   err.Error(),
			})
			return
		}

		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not update category",
			"error":   err.Error(),
		})
		return
	}

	httpCategory := presenter.NewHTTPCategory(*categoryResult)
	w.WriteHeader(http.StatusOK)
	w.Write(httpCategory.ToJSON())
}

type deleteCategoryHttpHandler struct {
	walletRepository repository.WalletRepository
	version          string
}

func NewDeleteCategoryHttpHandler(walletRepository repository.WalletRepository, version string) *deleteCategoryHttpHandler {
	return &deleteCategoryHttpHandler{
		walletRepository: walletRepository,
		version:          version,
	}
}

func (c *deleteCategoryHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	categoryId := vars["category_id"]

	if walletId == "" || categoryId == "" {
		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not get wallet id or category id from path",
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
			"message": "You are not allowed to delete categories in this wallet",
		})
		return
	}

	deleteCategoryUseCase := category.NewDeleteCategoryUseCase(c.walletRepository)
	err = deleteCategoryUseCase.Execute(category.DeleteCategoryUseCaseInput{
		WalletId:   walletId,
		CategoryId: categoryId,
	})

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			WriteError(w, http.StatusNotFound, map[string]any{
				"message": "Category not found",
				"error":   err.Error(),
			})
			return
		}

		WriteError(w, http.StatusInternalServerError, map[string]any{
			"message": "Could not delete category",
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
