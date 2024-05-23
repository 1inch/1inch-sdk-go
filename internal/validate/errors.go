package validate

import (
	"errors"
	"fmt"
	"strings"
)

func NewParameterValidationError(variableName string, errMessage string) error {
	return fmt.Errorf("config validation error '%s': %s", variableName, errMessage)
}

func NewParameterMissingError(variableName string) error {
	return fmt.Errorf("config validation error '%s' is required in the request config", variableName)
}

func NewParameterCustomError(errorMessage string) error {
	return fmt.Errorf("config validation error: %s", errorMessage)
}

func ConsolidateValidationErorrs(validationErrors []error) error {
	if len(validationErrors) == 0 {
		return nil
	}
	builder := strings.Builder{}
	builder.WriteString("request config errors: \n")
	for _, err := range validationErrors {
		builder.WriteString(err.Error())
		builder.WriteString("\n")
	}
	return errors.New(builder.String())
}

// GetValidatorErrorsCount uses the number of newlines in the error message to know how many errors were returned
func GetValidatorErrorsCount(validationError error) int {
	if validationError == nil {
		return 0
	}
	return strings.Count(validationError.Error(), "\n") - 1
}
