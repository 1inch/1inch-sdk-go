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
	a := new(big.Int)
	a, _ = a.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

	testcases := []struct {
		description       string
		fromToken         string
		name              string
		version           string
		publicAddress     string
		chainId           int
		spender           string
		amount            *big.Int
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
			amount:            a,
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
		chainId:    big.NewInt(int64(1)),
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
		amount            *big.Int
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
			amount:            big.NewInt(1000000000),
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
		chainId:    big.NewInt(int64(1)),
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
		amount        *big.Int
		//expectedSignatureString        string
		nonce                  int64
		deadline               int64
		IsDomainWithoutVersion bool
	}{
		{
			description:          "Create Permit 1inch BSC",
			fromToken:            "0x111111111117dc0aa78b770fa6a738034120c302",
			publicAddress:        "0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
			chainId:              56,
			name:                 "1INCH Token",
			version:              "1",
			nonce:                0,
			spender:              "0x11111112542d85b3ef69ae05771c2dccff4faa26",
			amount:               big.NewInt(1000000000),
			deadline:             192689033,
			expectedPermitString: "0x0000000000000000000000002c9b2DBdbA8A9c969Ac24153f5C1c23CB0e6391400000000000000000000000011111112542d85b3ef69ae05771c2dccff4faa26000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000b7c3389000000000000000000000000000000000000000000000000000000000000001c3b448216a78f91e84db06cf54eb1e3758425bd97ffb9d6941ce437ec7a9c2c174c94f1fa492007dea3a3c305353bf3430b1ca506dd630ce1fd3da09bd387b2f3",
		},
		{
			description:          "Create Permit 1inch ETH",
			fromToken:            "0x111111111117dC0aa78b770fA6A738034120C302",
			publicAddress:        "0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
			chainId:              1,
			name:                 "1INCH Token",
			version:              "1",
			nonce:                1,
			spender:              "0x1111111254EEB25477B68fb85Ed929f73A960582",
			amount:               big.NewInt(100000),
			deadline:             1713453855,
			expectedPermitString: "0x0000000000000000000000002c9b2DBdbA8A9c969Ac24153f5C1c23CB0e639140000000000000000000000001111111254EEB25477B68fb85Ed929f73A96058200000000000000000000000000000000000000000000000000000000000186a00000000000000000000000000000000000000000000000000000000066213b1f000000000000000000000000000000000000000000000000000000000000001b5c8a5ea8fba76eb6f7ad00260b345420d2340ef5226a66d0d0124fc715c0b95538cf9fc2730782c297378d5f8d17694ad389a8270d0513742b1d6a796189c358",
		},
		{
			description:            "Create Permit Gitcoin ETH without version",
			fromToken:              "0xDe30da39c46104798bB5aA3fe8B9e0e1F348163F",
			publicAddress:          "0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
			chainId:                1,
			name:                   "Gitcoin",
			version:                "1",
			nonce:                  0,
			spender:                "0x1111111254EEB25477B68fb85Ed929f73A960582",
			amount:                 big.NewInt(100000),
			deadline:               1713454178,
			IsDomainWithoutVersion: true,
			expectedPermitString:   "0x0000000000000000000000002c9b2DBdbA8A9c969Ac24153f5C1c23CB0e639140000000000000000000000001111111254EEB25477B68fb85Ed929f73A96058200000000000000000000000000000000000000000000000000000000000186a00000000000000000000000000000000000000000000000000000000066213c62000000000000000000000000000000000000000000000000000000000000001b156cb83f6df524a321d7288c57411815bd15852f622583a585ad6679b9c162d263cffe30a293174924b8665fcc123298e0019e2c0a2846d048c7d03004a67e22",
		},
		{
			description:          "Create Permit LQTY ETH",
			fromToken:            "0x6DEA81C8171D0bA574754EF6F8b412F2Ed88c54D",
			publicAddress:        "0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
			chainId:              1,
			name:                 "LQTY",
			version:              "1",
			nonce:                0,
			spender:              "0x1111111254EEB25477B68fb85Ed929f73A960582",
			amount:               big.NewInt(100000),
			deadline:             1713456277,
			expectedPermitString: "0x0000000000000000000000002c9b2DBdbA8A9c969Ac24153f5C1c23CB0e639140000000000000000000000001111111254EEB25477B68fb85Ed929f73A96058200000000000000000000000000000000000000000000000000000000000186a00000000000000000000000000000000000000000000000000000000066214495000000000000000000000000000000000000000000000000000000000000001bcfe6e19c5525b4a49ea2b963e36effa06151db0e40c731ffb183c2f5019ce2b763be0425573471fa7b63cbedee5a1dc1f833353ed52021e02c23d8f53b96aeb8",
		},
		{
			description:          "Create Permit USDC ETH",
			fromToken:            "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			publicAddress:        "0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
			chainId:              1,
			name:                 "USD Coin",
			version:              "2",
			nonce:                0,
			spender:              "0x1111111254EEB25477B68fb85Ed929f73A960582",
			amount:               big.NewInt(100000),
			deadline:             1713457338,
			expectedPermitString: "0x0000000000000000000000002c9b2DBdbA8A9c969Ac24153f5C1c23CB0e639140000000000000000000000001111111254EEB25477B68fb85Ed929f73A96058200000000000000000000000000000000000000000000000000000000000186a000000000000000000000000000000000000000000000000000000000662148ba000000000000000000000000000000000000000000000000000000000000001be1e089fc6e42874b2d369a080af21ddc227181943d680f401da42d2c50ca8d646785db86ee82618b6c5b93eaa0954f12b30b086063001e7f66161157aaae652f",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			d := common.ContractPermitData{
				FromToken:              tc.fromToken,
				Spender:                tc.spender,
				Name:                   tc.name,
				Version:                tc.version,
				PublicAddress:          tc.publicAddress,
				ChainId:                tc.chainId,
				Nonce:                  tc.nonce,
				Deadline:               tc.deadline,
				Amount:                 tc.amount,
				IsDomainWithoutVersion: tc.IsDomainWithoutVersion,
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
