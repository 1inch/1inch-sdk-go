package web3_provider

import (
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"
)

func Test_createPermitSignature(t *testing.T) {
	erc20ABI, err := abi.JSON(strings.NewReader(constants.Erc20ABI)) // Make a generic version of this ABI
	require.NoError(t, err)

	privateKey, err := crypto.HexToECDSA("ad21c0552a3b52e94520da713455cc347e4e89628a334be24d85b8083848434f")
	require.NoError(t, err)

	publicKey := privateKey.Public()
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	w := &Wallet{
		ethClient:  &ethclient.Client{},
		address:    &address,
		privateKey: privateKey,
		chainId:    big.NewInt(int64(137)),
		erc20ABI:   &erc20ABI,
	}

	testcases := []struct {
		description       string
		fromToken         string
		name              string
		version           string
		publicAddress     string
		chainId           int
		spender           string
		amount            string
		nonce             int64
		deadline          int64
		expectedSignature string
	}{
		{
			description:       "Create Signature",
			fromToken:         "0x45c32fA6DF82ead1e2EF74d17b76547EDdFaFF89",
			publicAddress:     "0x2a250893f86Dc8497E131508f680338ac647B498",
			chainId:           137,
			name:              "Frax",
			version:           "1",
			nonce:             0,
			spender:           "0x1111111254eeb25477b68fb85ed929f73a960582",
			amount:            "115792089237316195423570985008687907853269984665640564039457584007913129639935",
			deadline:          1704250835,
			expectedSignature: "0d95c0246c1356df4653606e586e97447a516c937b5dd758fa0e56f2f8dd1f952b222c24a337e89dfbe20a8e112a7c6d004a3170598b9d4941aa38126920c9ed1b",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			d := common.ContractPermitData{
				FromToken:     tc.fromToken,
				Spender:       tc.spender,
				Name:          tc.name,
				Version:       tc.version,
				PublicAddress: tc.publicAddress,
				ChainId:       tc.chainId,
				Nonce:         tc.nonce,
				Deadline:      tc.deadline,
				Amount:        tc.amount,
			}

			result, err := w.createPermitSignature(&d)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSignature, result)
		})
	}
}

func Test_createPermitSignature2(t *testing.T) {
	erc20ABI, err := abi.JSON(strings.NewReader(constants.Erc20ABI)) // Make a generic version of this ABI
	require.NoError(t, err)

	privateKey, err := crypto.HexToECDSA("965e092fdfc08940d2bd05c7b5c7e1c51e283e92c7f52bbf1408973ae9a9acb7")
	require.NoError(t, err)

	publicKey := privateKey.Public()
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	w := &Wallet{
		ethClient:  &ethclient.Client{},
		address:    &address,
		privateKey: privateKey,
		chainId:    big.NewInt(int64(137)),
		erc20ABI:   &erc20ABI,
	}

	testcases := []struct {
		description       string
		fromToken         string
		name              string
		version           string
		publicAddress     string
		chainId           int
		spender           string
		amount            string
		nonce             int64
		deadline          int64
		expectedSignature string
	}{
		{
			description:       "Create Signature 2",
			fromToken:         "0x111111111117dc0aa78b770fa6a738034120c302",
			publicAddress:     "0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
			chainId:           56,
			name:              "1INCH Token",
			version:           "1",
			nonce:             0,
			spender:           "0x11111112542d85b3ef69ae05771c2dccff4faa26",
			amount:            "1000000000",
			deadline:          192689033,
			expectedSignature: "3b448216a78f91e84db06cf54eb1e3758425bd97ffb9d6941ce437ec7a9c2c174c94f1fa492007dea3a3c305353bf3430b1ca506dd630ce1fd3da09bd387b2f31c",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			d := common.ContractPermitData{
				FromToken:     tc.fromToken,
				Spender:       tc.spender,
				Name:          tc.name,
				Version:       tc.version,
				PublicAddress: tc.publicAddress,
				ChainId:       tc.chainId,
				Nonce:         tc.nonce,
				Deadline:      tc.deadline,
				Amount:        tc.amount,
			}

			result, err := w.createPermitSignature(&d)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSignature, result)
		})
	}
}

func TestTokenPermit(t *testing.T) {
	erc20ABI, err := abi.JSON(strings.NewReader(constants.Erc20ABI)) // Make a generic version of this ABI
	require.NoError(t, err)

	privateKey, err := crypto.HexToECDSA("965e092fdfc08940d2bd05c7b5c7e1c51e283e92c7f52bbf1408973ae9a9acb7")
	require.NoError(t, err)

	publicKey := privateKey.Public()
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	w := &Wallet{
		ethClient:  &ethclient.Client{},
		address:    &address,
		privateKey: privateKey,
		chainId:    big.NewInt(int64(137)),
		erc20ABI:   &erc20ABI,
	}

	testcases := []struct {
		description          string
		expectedPermitString string

		fromToken     string
		name          string
		version       string
		publicAddress string
		chainId       int
		spender       string
		amount        string
		//expectedSignatureString        string
		nonce    int64
		deadline int64
	}{
		{
			description:   "Create Permit parameter",
			fromToken:     "0x111111111117dc0aa78b770fa6a738034120c302",
			publicAddress: "0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
			chainId:       56,
			name:          "1INCH Token",
			version:       "1",
			nonce:         0,
			spender:       "0x11111112542d85b3ef69ae05771c2dccff4faa26",
			amount:        "1000000000",
			deadline:      192689033,
			//expectedSignatureString: "0x3b448216a78f91e84db06cf54eb1e3758425bd97ffb9d6941ce437ec7a9c2c174c94f1fa492007dea3a3c305353bf3430b1ca506dd630ce1fd3da09bd387b2f31c",
			expectedPermitString: "0x0000000000000000000000002c9b2DBdbA8A9c969Ac24153f5C1c23CB0e6391400000000000000000000000011111112542d85b3ef69ae05771c2dccff4faa26000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000b7c3389000000000000000000000000000000000000000000000000000000000000001c3b448216a78f91e84db06cf54eb1e3758425bd97ffb9d6941ce437ec7a9c2c174c94f1fa492007dea3a3c305353bf3430b1ca506dd630ce1fd3da09bd387b2f3",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			d := common.ContractPermitData{
				FromToken:     tc.fromToken,
				Spender:       tc.spender,
				Name:          tc.name,
				Version:       tc.version,
				PublicAddress: tc.publicAddress,
				ChainId:       tc.chainId,
				Nonce:         tc.nonce,
				Deadline:      tc.deadline,
				Amount:        tc.amount,
			}

			permit, err := w.TokenPermit(d)
			require.NoError(t, err)
			// temp, need to work for good example
			// finish it soon
			require.Equal(t, tc.expectedPermitString, permit)
		})
	}
}

func Test_convertSignatureToVRSString(t *testing.T) {
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
			result := convertSignatureToVRSString(tc.signature)
			if tc.expectedSignature != "" {
				require.Equal(t, tc.expectedSignature, result)
			}
		})
	}
}

func Test_padStringWithZeroes(t *testing.T) {
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

func Test_remove0xPrefix(t *testing.T) {
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
			result := remove0xPrefix(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
