package swap

import (
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"1inch-sdk-golang/helpers/consts/amounts"
)

func TestCreatePermitSignature(t *testing.T) {
	testcases := []struct {
		description       string
		fromToken         string
		name              string
		publicAddress     string
		chainId           int
		key               string
		nonce             int64
		deadline          int64
		expectedSignature string
	}{
		{
			description:       "Create Signature",
			fromToken:         "0x45c32fA6DF82ead1e2EF74d17b76547EDdFaFF89",
			publicAddress:     "0x2a250893f86Dc8497E131508f680338ac647B498",
			chainId:           137,
			key:               "ad21c0552a3b52e94520da713455cc347e4e89628a334be24d85b8083848434f",
			name:              "Frax",
			nonce:             0,
			deadline:          1704250835,
			expectedSignature: "0x0d95c0246c1356df4653606e586e97447a516c937b5dd758fa0e56f2f8dd1f952b222c24a337e89dfbe20a8e112a7c6d004a3170598b9d4941aa38126920c9ed1b",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {

			config := &PermitSignatureConfig{
				FromToken:     tc.fromToken,
				Name:          tc.name,
				PublicAddress: tc.publicAddress,
				ChainId:       tc.chainId,
				Key:           tc.key,
				Nonce:         tc.nonce,
				Deadline:      tc.deadline,
			}

			result, err := CreatePermitSignature(config)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSignature, result)
		})
	}
}

func TestCreatePermitParams(t *testing.T) {
	testcases := []struct {
		description          string
		Owner                string
		Spender              string
		Value                *big.Int
		Deadline             int64
		Signature            string
		expectedPermitString string
	}{
		{
			description:          "Create Permit parameter",
			Owner:                "0x50c5df26654b5efbdd0c54a062dfa6012933defe",
			Spender:              "0x1111111254eeb25477b68fb85ed929f73a960582",
			Value:                amounts.BigMaxUint256,
			Deadline:             1704250835,
			Signature:            "c8dcab9ab2ce2055e61c0718117f8d77a56cd0a8b8370d8f5e16932a60d21a3e0eb0214dcbe4e7c5131cc45fd552e12f5bcef3b9c7fcb47ace9d4f694a496d471b",
			expectedPermitString: "0x00000000000000000000000050c5df26654b5efbdd0c54a062dfa6012933defe0000000000000000000000001111111254eeb25477b68fb85ed929f73a960582ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000000000006594cdd3000000000000000000000000000000000000000000000000000000000000001bc8dcab9ab2ce2055e61c0718117f8d77a56cd0a8b8370d8f5e16932a60d21a3e0eb0214dcbe4e7c5131cc45fd552e12f5bcef3b9c7fcb47ace9d4f694a496d47",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {

			config := &PermitParamsConfig{
				Owner:     tc.Owner,
				Spender:   tc.Spender,
				Value:     tc.Value,
				Deadline:  tc.Deadline,
				Signature: tc.Signature,
			}

			result := CreatePermitParams(config)
			require.Equal(t, tc.expectedPermitString, result)
		})
	}
}

func TestConvertSignatureToVRSString(t *testing.T) {
	testcases := []struct {
		description             string
		signature               string
		expectedSignatureString string
	}{
		{
			description:             "Create Signature",
			signature:               "c8dcab9ab2ce2055e61c0718117f8d77a56cd0a8b8370d8f5e16932a60d21a3e0eb0214dcbe4e7c5131cc45fd552e12f5bcef3b9c7fcb47ace9d4f694a496d471b",
			expectedSignatureString: "000000000000000000000000000000000000000000000000000000000000001bc8dcab9ab2ce2055e61c0718117f8d77a56cd0a8b8370d8f5e16932a60d21a3e0eb0214dcbe4e7c5131cc45fd552e12f5bcef3b9c7fcb47ace9d4f694a496d47",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			result := ConvertSignatureToVRSString(tc.signature)
			require.Equal(t, tc.expectedSignatureString, result)
		})
	}
}

func TestPadStringWithZeroes(t *testing.T) {
	testcases := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "String shorter than 64 characters",
			input:       "abc",
			expected:    "0000000000000000000000000000000000000000000000000000000000000abc",
		},
		{
			description: "String exactly 64 characters",
			input:       strings.Repeat("a", 64),
			expected:    strings.Repeat("a", 64),
		},
		{
			description: "String longer than 64 characters",
			input:       strings.Repeat("b", 65),
			expected:    strings.Repeat("b", 65),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			result := padStringWithZeroes(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestRemove0xPrefix(t *testing.T) {
	testcases := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "String with 0x prefix",
			input:       "0x12345",
			expected:    "12345",
		},
		{
			description: "String without 0x prefix",
			input:       "12345",
			expected:    "12345",
		},
		{
			description: "Empty string",
			input:       "",
			expected:    "",
		},
		{
			description: "String with only 0x",
			input:       "0x",
			expected:    "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			result := Remove0xPrefix(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
