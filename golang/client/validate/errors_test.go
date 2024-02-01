package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetValidatorErrorsCount(t *testing.T) {
	testCases := []struct {
		description   string
		validationErr error
		expectedCount int
	}{
		{
			description:   "No validation errors",
			validationErr: nil,
			expectedCount: 0,
		},
		{
			description: "Single validation error",
			validationErr: AggregateValidationErorrs([]error{
				NewParameterValidationError("chainId", "is required"),
			}),
			expectedCount: 1,
		},
		{
			description: "Multiple validation errors",
			validationErr: AggregateValidationErorrs([]error{
				NewParameterValidationError("chainId", "is required"),
				NewParameterValidationError("walletKey", "not a valid private key"),
				NewParameterValidationError("sourceWallet", "not a valid Ethereum address"),
			}),
			expectedCount: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			count := GetValidatorErrorsCount(tc.validationErr)
			require.Equal(t, tc.expectedCount, count, "The count of validator errors should match the expected value")
		})
	}
}
