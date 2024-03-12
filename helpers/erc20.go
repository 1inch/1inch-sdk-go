package helpers

import "strings"

func SimplifyValue(input string, decimalPlaces int) string {
	// Padding with zeros if necessary
	for len(input) < decimalPlaces {
		input = "0" + input
	}

	// Inserting decimal point
	pointIndex := len(input) - decimalPlaces
	if pointIndex == 0 {
		return "0." + strings.TrimRight(input, "0")
	}

	// Form the result and trim trailing zeros after the decimal point
	result := input[:pointIndex] + "." + input[pointIndex:]
	return strings.TrimRight(strings.TrimRight(result, "0"), ".")
}
