package orderbook

import (
	"math/big"
	"strings"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	web3_provider "github.com/1inch/1inch-sdk-go/v4/internal/web3-provider"
)

func TestBuildPermit2Calldata(t *testing.T) {
	testKey := "d8d1f95deb28949ea0ecc4e9a0decf89e98422c2d76ab6e5f736792a388c56c7"

	tests := []struct {
		name        string
		params      Permit2PermitParams
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Known values",
			params: Permit2PermitParams{
				Token:       gethCommon.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
				Amount:      big.NewInt(100000000000000000),
				Expiration:  constants.Uint48Max,
				Nonce:       big.NewInt(0),
				Spender:     gethCommon.HexToAddress("0x111111125421cA6dc452d289314280a0f8842A65"),
				SigDeadline: constants.Uint48Max,
			},
			expected: "0x000000000000000000000000a07c1d51497fb6e66aa2329cecb86fca0a957fdb000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000000000000000000000000000000000ffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000111111125421ca6dc452d289314280a0f8842a650000000000000000000000000000000000000000000000000000ffffffffffff00000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000040977350b6dbbf5b1308a4b231bf73c8e779cc15d7be0c9c48bdfaed90528ae66da6fbe646b71e7f00265deac46e2c89bc2827f0d1029839cd18f82e24a8ee90ad",
		},
		{
			name: "Missing nonce",
			params: Permit2PermitParams{
				Token:       gethCommon.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
				Amount:      big.NewInt(1),
				Expiration:  constants.Uint48Max,
				Spender:     gethCommon.HexToAddress("0x111111125421cA6dc452d289314280a0f8842A65"),
				SigDeadline: constants.Uint48Max,
			},
			expectError: true,
			errorMsg:    "amount, expiration, nonce, and sig deadline are required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wallet, err := web3_provider.DefaultWalletOnlyProvider(testKey, 1)
			require.NoError(t, err)

			result, err := BuildPermit2Calldata(wallet, tc.params)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)

				// The Limit Order Protocol recognizes Permit2 permits by their exact
				// 352-byte length; the compact EIP-2098 signature keeps it there
				assert.Equal(t, 352, (len(result)-2)/2, "permit calldata must be 352 bytes")

				// The owner is the signing wallet
				ownerHex := strings.ToLower(wallet.Address().Hex()[2:])
				assert.Equal(t, ownerHex, result[26:66], "owner must be the wallet address")
			}
		})
	}
}

func TestBuildPermit2CalldataCompact(t *testing.T) {
	testKey := "d8d1f95deb28949ea0ecc4e9a0decf89e98422c2d76ab6e5f736792a388c56c7"

	baseParams := func() Permit2PermitParams {
		return Permit2PermitParams{
			Token:       gethCommon.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
			Amount:      big.NewInt(100000000000000000),
			Expiration:  constants.Uint48Max,
			Nonce:       big.NewInt(0),
			Spender:     gethCommon.HexToAddress("0x111111125421cA6dc452d289314280a0f8842A65"),
			SigDeadline: constants.Uint48Max,
		}
	}

	tests := []struct {
		name        string
		mutate      func(*Permit2PermitParams)
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name:     "Known values with unlimited timestamps",
			mutate:   func(p *Permit2PermitParams) {},
			expected: "0x000000000000000000000000016345785d8a0000000000000000000000000000977350b6dbbf5b1308a4b231bf73c8e779cc15d7be0c9c48bdfaed90528ae66da6fbe646b71e7f00265deac46e2c89bc2827f0d1029839cd18f82e24a8ee90ad",
		},
		{
			name: "Expiration too large for compact encoding",
			mutate: func(p *Permit2PermitParams) {
				p.Expiration = new(big.Int).Lsh(big.NewInt(1), 33)
			},
			expectError: true,
			errorMsg:    "expiration: value must be max uint48 or at most 2^32 - 2",
		},
		{
			name: "Nonce too large for compact encoding",
			mutate: func(p *Permit2PermitParams) {
				p.Nonce = new(big.Int).Lsh(big.NewInt(1), 33)
			},
			expectError: true,
			errorMsg:    "nonce must fit in uint32",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wallet, err := web3_provider.DefaultWalletOnlyProvider(testKey, 1)
			require.NoError(t, err)

			params := baseParams()
			tc.mutate(&params)

			result, err := BuildPermit2CalldataCompact(wallet, params)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
				assert.Equal(t, 96, (len(result)-2)/2, "compact permit calldata must be 96 bytes")
			}
		})
	}
}

func TestCompactPermit2Timestamp(t *testing.T) {
	tests := []struct {
		name        string
		value       *big.Int
		expected    uint32
		expectError bool
	}{
		{
			name:     "Max uint48 stores zero",
			value:    constants.Uint48Max,
			expected: 0,
		},
		{
			name:     "Regular timestamp stores value plus one",
			value:    big.NewInt(1715201499),
			expected: 1715201500,
		},
		{
			name:     "Largest encodable value",
			value:    big.NewInt(1<<32 - 2),
			expected: 1<<32 - 1,
		},
		{
			name:        "One above largest encodable value",
			value:       big.NewInt(1<<32 - 1),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := compactPermit2Timestamp(tc.value)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestBuildPermit2CalldataCompact_NegativeAmount(t *testing.T) {
	testKey := "d8d1f95deb28949ea0ecc4e9a0decf89e98422c2d76ab6e5f736792a388c56c7"
	wallet, err := web3_provider.DefaultWalletOnlyProvider(testKey, 1)
	require.NoError(t, err)

	_, err = BuildPermit2CalldataCompact(wallet, Permit2PermitParams{
		Token:       gethCommon.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		Amount:      big.NewInt(-1),
		Expiration:  constants.Uint48Max,
		Nonce:       big.NewInt(0),
		Spender:     gethCommon.HexToAddress("0x111111125421cA6dc452d289314280a0f8842A65"),
		SigDeadline: constants.Uint48Max,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "amount must fit in uint160")
}
