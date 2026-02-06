package validate

func CheckStringRequired(value string, variableName string) error {
	if value == "" {
		return NewParameterMissingError(variableName)
	}

	return nil
}

// CheckString is a no-op validator. Type safety is now enforced at compile time
// via the generic Parameter[T] function. This function exists for use with the
// validation framework when a string field needs to be included in validation
// chains but has no additional constraints.
func CheckString(_ string, _ string) error {
	return nil
}

// CheckBoolean is a no-op validator. Type safety is now enforced at compile time
// via the generic Parameter[T] function. This function exists for use with the
// validation framework when a bool field needs to be included in validation
// chains but has no additional constraints.
func CheckBoolean(_ bool, _ string) error {
	return nil
}
