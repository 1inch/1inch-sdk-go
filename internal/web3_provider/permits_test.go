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

	"github.com/1inch/1inch-sdk-go/internal/common"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/abis"
)

func TestCreatePermitSignature(t *testing.T) {
	erc20ABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
	require.NoError(t, err)

	privateKey, err := crypto.HexToECDSA("ad21c0552a3b52e94520da713455cc347e4e89628a334be24d85b8083848434f")
	require.NoError(t, err)

	publicKey := privateKey.Public()
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	w := &Wallet{
		ethClient:  &ethclient.Client{},
		address:    &address,
		privateKey: privateKey,
		ChainId:    big.NewInt(int64(137)),
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
			expectedSignature: "0x0d95c0246c1356df4653606e586e97447a516c937b5dd758fa0e56f2f8dd1f952b222c24a337e89dfbe20a8e112a7c6d004a3170598b9d4941aa38126920c9ed1b",
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

			result, err := w.TokenPermit(d)
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
