package errx

import (
	"errors"
	"strings"
)

var (
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrInvalidCreatorID        = errors.New("invalid creator id")
	ErrInvalidInput            = errors.New("invalid input")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInvalidTransactionType  = errors.New("invalid transaction type")
	ErrNotFound                = errors.New("entry not found")
)

func MissingRequiredFieldsError(fields ...string) error {
	return errors.New("missing required field: " + strings.Join(fields, ","))
}
