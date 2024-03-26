package onchain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/helpers/consts/chains"
)

func TestCreatePermitSignature(t *testing.T) {
	testcases := []struct {
		description       string
		fromToken         string
		name              string
		version           string
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
			chainId:           chains.Polygon,
			key:               "ad21c0552a3b52e94520da713455cc347e4e89628a334be24d85b8083848434f",
			name:              "Frax",
			version:           "1",
			nonce:             0,
			deadline:          1704250835,
			expectedSignature: "0x55dcf81b7366e4f6e5e5f5b335340164b4f4e6d2b4a63c7dbaf937a2f1a4ec380263e1676c814d72ef969a8ce4da2eb78d8b4d60b8b70a1b820a214cdf21acb21b",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {

			config := &PermitSignatureConfig{
				FromToken:     tc.fromToken,
				Name:          tc.name,
				Version:       tc.version,
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

func TestConvertSignatureToVRSString(t *testing.T) {
	testcases := []struct {
		description       string
		signature         string
		expectedSignature string
	}{
		{
			description:       "Create Signature",
			signature:         "c8dcab9ab2ce2055e61c0718117f8d77a56cd0a8b8370d8f5e16932a60d21a3e0eb0214dcbe4e7c5131cc45fd552e12f5bcef3b9c7fcb47ace9d4f694a496d471b",
			expectedSignature: "000000000000000000000000000000000000000000000000000000000000001bc8dcab9ab2ce2055e61c0718117f8d77a56cd0a8b8370d8f5e16932a60d21a3e0eb0214dcbe4e7c5131cc45fd552e12f5bcef3b9c7fcb47ace9d4f694a496d47",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			result := ConvertSignatureToVRSString(tc.signature)
			if tc.expectedSignature != "" {
				require.Equal(t, tc.expectedSignature, result)
			}
		})
	}
}
