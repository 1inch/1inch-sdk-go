package helpers

func BoolPtr(b bool) *bool {
	return &b
}

func Float32Ptr(v int) *float32 {
	f := float32(v)
	return &f
}

// GetPtr returns a pointer to the provided value.
func GetPtr[T any](value T) *T {
	return &value
}
