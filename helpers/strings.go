package helpers

import (
	"math/big"
	"regexp"
)

// IsScientificNotation checks if the string is in scientific notation (like 1e+18).
func IsScientificNotation(s string) bool {
	// This regular expression matches strings in the format of "1e+18", "2.3e-4", etc.
	re := regexp.MustCompile(`^[+-]?\d+(\.\d+)?[eE][+-]?\d+$`)
	return re.MatchString(s)
}

func ExpandScientificNotation(s string) (string, error) {
	f, _, err := big.ParseFloat(s, 10, 0, big.ToNearestEven)
	if err != nil {
		return "", err
	}

	// Use a precision that is sufficient to handle small numbers.
	// The precision here is set to a large number to ensure accuracy for small decimal values.
	f.SetPrec(64)

	return f.Text('f', -1), nil // -1 ensures that insignificant zeroes are not omitted
}
