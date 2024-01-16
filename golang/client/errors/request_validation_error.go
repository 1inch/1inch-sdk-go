package clienterrors

import (
	"errors"
	"fmt"
)

type requestValidationError struct {
	WrappedError error
}

func NewRequestValidationError(message string) *requestValidationError {
	return &requestValidationError{WrappedError: errors.New(message)}
}

func (e *requestValidationError) Error() string {
	return fmt.Errorf("request validation error: %v", e.WrappedError).Error()
}
