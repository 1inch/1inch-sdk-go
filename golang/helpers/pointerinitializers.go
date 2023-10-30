package helpers

// GetPtr returns a pointer to the provided value.
func GetPtr[T any](value T) *T {
	return &value
}
