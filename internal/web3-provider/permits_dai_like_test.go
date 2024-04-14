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

func Test_createPermitSignatureDaiLike(t *testing.T) {
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
		allowed           bool
		nonce             int64
		deadline          int64
		expectedSignature string
	}{
		{
			description:       "Create Signature",
			fromToken:         "0x111111111117dc0aa78b770fa6a738034120c302",
			publicAddress:     "0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
			chainId:           56,
			name:              "1INCH Token",
			version:           "1",
			allowed:           true,
			nonce:             0,
			spender:           "0x11111112542d85b3ef69ae05771c2dccff4faa26",
			deadline:          192689033,
			expectedSignature: "cdcf508eed2f330082c6a19ba5931ebbab16efd470dee2072440aee35c064b736b31b4eed202958a43e250f0a5321db09185f1525776015ecaa8975ca7cf157d1b",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			d := common.ContractPermitDataDaiLike{
				FromToken: tc.fromToken,
				Spender:   tc.spender,
				Name:      tc.name,
				Version:   tc.version,
				Holder:    tc.publicAddress,
				ChainId:   tc.chainId,
				Nonce:     tc.nonce,
				Expiry:    tc.deadline,
				Allowed:   tc.allowed,
			}

			result, err := w.createPermitSignatureDaiLike(&d)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSignature, result)
		})
	}
}

func TestTokenPermitDaiLike(t *testing.T) {
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
			deadline:      192689033,
			//expectedSignatureString: "0x3b448216a78f91e84db06cf54eb1e3758425bd97ffb9d6941ce437ec7a9c2c174c94f1fa492007dea3a3c305353bf3430b1ca506dd630ce1fd3da09bd387b2f31c",
			expectedPermitString: "0x0000000000000000000000002c9b2DBdbA8A9c969Ac24153f5C1c23CB0e6391400000000000000000000000011111112542d85b3ef69ae05771c2dccff4faa26000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000b7c3389000000000000000000000000000000000000000000000000000000000000001c3b448216a78f91e84db06cf54eb1e3758425bd97ffb9d6941ce437ec7a9c2c174c94f1fa492007dea3a3c305353bf3430b1ca506dd630ce1fd3da09bd387b2f3",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			d := common.ContractPermitDataDaiLike{
				FromToken: tc.fromToken,
				Spender:   tc.spender,
				Name:      tc.name,
				Version:   tc.version,
				Holder:    tc.publicAddress,
				ChainId:   tc.chainId,
				Nonce:     tc.nonce,
				Expiry:    tc.deadline,
			}

			permit, err := w.TokenPermitDaiLike(d)
			require.NoError(t, err)
			// temp, need to work for good example
			// finish it soon
			require.Equal(t, tc.expectedPermitString, permit)
		})
	}
}
