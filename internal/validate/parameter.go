package validate

// Parameter validates a parameter and appends any error to the validationErrors slice.
func Parameter[T any](parameter T, variableName string, validationFunc func(T, string) error, validationErrors []error) []error {
	err := validationFunc(parameter, variableName)
	if err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validationErrors
}
