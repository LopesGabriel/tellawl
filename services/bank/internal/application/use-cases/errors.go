package usecases

import (
	"errors"
	"strings"
)

var ErrInvalidCreatorID = errors.New("invalid creator id")
var ErrInvalidInput = errors.New("invalid input")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInsufficientPermissions = errors.New("insufficient permissions")

func MissingRequiredFieldsError(fields ...string) error {
	return errors.New("missing required field: " + strings.Join(fields, ","))
}
