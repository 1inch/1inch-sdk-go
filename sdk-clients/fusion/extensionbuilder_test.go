package fusion

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
)

func TestGenerateSalt(t *testing.T) {
	// Save the original function
	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc

	// Monkey patch the function
	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
		return big.NewInt(123456), nil
	}

	// Restore the original function after the test
	defer func() {
		random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc
	}()

	tests := []struct {
		name      string
		extension *Extension
		expected  string
		expectErr bool
	}{
		{
			name: "Generate salt when extension is not empty",
			extension: &Extension{
				MakerAssetSuffix: "suffix1",
				TakerAssetSuffix: "suffix2",
				MakingAmountData: "data1",
				TakingAmountData: "data2",
				Predicate:        "predicate",
				MakerPermit:      "permit",
				PreInteraction:   "pre",
				PostInteraction:  "post",
				CustomData:       "custom",
			},
			expected:  "180431658011416401710119735245975317914670388782711199",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expected, err := BigIntFromString(tc.expected)
			require.NoError(t, err)

			result, err := tc.extension.GenerateSalt()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expected, result)
			}
		})
	}
}
