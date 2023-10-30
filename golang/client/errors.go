package client

import "fmt"

type requestValidationError struct {
	WrappedError error
}

func NewRequestValidationError(e error) *requestValidationError {
	return &requestValidationError{WrappedError: e}
}

func (e *requestValidationError) Error() string {
	return fmt.Errorf("request validation error: %v", e.WrappedError).Error()
}
