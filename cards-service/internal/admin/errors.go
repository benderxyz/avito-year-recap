package admin

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
)

type ValidationError struct {
	Message string
	Fields  map[string]string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type ReferenceError struct {
	Message    string
	References []string
}

func (e *ReferenceError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.References)
}

type fieldErrors map[string]string

func (f fieldErrors) add(field, message string) {
	if _, exists := f[field]; !exists {
		f[field] = message
	}
}

func (f fieldErrors) err(message string) error {
	if len(f) == 0 {
		return nil
	}
	return &ValidationError{Message: message, Fields: f}
}

func invalidBody(message string) error {
	return &ValidationError{Message: message}
}
