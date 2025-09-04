package presenter

import (
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
)

type HTTPMonetary struct {
	Value  int `json:"value"`
	Offset int `json:"offset"`
}

func NewHTTPMonetary(monetary models.Monetary) HTTPMonetary {
	return HTTPMonetary{
		Value:  monetary.Value,
		Offset: monetary.Offset,
	}
}
