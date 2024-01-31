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

func AggregateValidationErorrs(validationErrors []error) error {
	builder := strings.Builder{}
	builder.WriteString("request config errors: \n")
	for _, err := range validationErrors {
		builder.WriteString(err.Error())
		builder.WriteString("\n")
	}
	return errors.New(builder.String())
}
