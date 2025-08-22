package models

import (
	"time"

	"github.com/google/uuid"
)

type CategoryType string

const (
	CategoryTypeDefault CategoryType = "default"
	CategoryTypeCustom  CategoryType = "custom"
)

type Category struct {
	Id       string
	WalletId string
	Name     string
	Type     CategoryType

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func CreateDefaultCategory(walletId, name string) *Category {
	return &Category{
		Id:        uuid.NewString(),
		WalletId:  walletId,
		Name:      name,
		Type:      CategoryTypeDefault,
		CreatedAt: time.Now(),
	}
}

func CreateCustomCategory(walletId, name string) *Category {
	return &Category{
		Id:        uuid.NewString(),
		WalletId:  walletId,
		Name:      name,
		Type:      CategoryTypeCustom,
		CreatedAt: time.Now(),
	}
}

func GetDefaultCategories() []string {
	return []string{
		"Alimentação",
		"Transporte",
		"Moradia",
		"Saúde",
		"Educação",
		"Lazer",
		"Investimento",
		"Salário",
		"Imprevisto",
		"Outros",
	}
}
