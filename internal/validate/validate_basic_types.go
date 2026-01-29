package validate

import "fmt"

func CheckStringRequired(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("'%s' must be a string", variableName)
	}

	if value == "" {
		return NewParameterMissingError(variableName)
	}

	return CheckString(value, variableName)
}

func CheckString(parameter interface{}, variableName string) error {
	_, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("'%s' must be a string", variableName)
	}

	return nil
}

func CheckBoolean(parameter interface{}, variableName string) error {
	_, ok := parameter.(bool)
	if !ok {
		return fmt.Errorf("'%s' must be a boolean", variableName)
	}

	return nil
}
