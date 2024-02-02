package validate

type ValidationFunc func(parameter interface{}, variableName string) error

func Parameter(parameter interface{}, variableName string, validationFunc ValidationFunc, validationErrors []error) []error {
	var err error
	err = validationFunc(parameter, variableName)
	if err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validationErrors
}
