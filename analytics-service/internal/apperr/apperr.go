package apperr

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func Validation(message string) error {
	return &ValidationError{Message: message}
}

func Validationf(format string, args ...any) error {
	return &ValidationError{Message: fmt.Sprintf(format, args...)}
}

func IsValidation(err error) bool {
	var target *ValidationError
	return errors.As(err, &target)
}
