package deploykit

import (
	"errors"
)

var ErrValidation = errors.New("validation error")

type ValidationError string
type ValidationErrors map[string][]ValidationError

func (e *ValidationErrors) Error() string {
	return ErrValidation.Error()
}
