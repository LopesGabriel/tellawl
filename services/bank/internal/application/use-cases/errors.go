package usecases

import (
	"errors"
	"strings"
)

var ErrInvalidCreatorID = errors.New("invalid creator id")

func MissingRequiredFieldsError(fields ...string) error {
	return errors.New("missing required field: " + strings.Join(fields, ","))
}
