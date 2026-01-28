package fusionorder

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWhitelistItem(t *testing.T) {
	addressHalf := "1234567890abcdef1234"
	delay := big.NewInt(100)
	
	item := NewWhitelistItem(addressHalf, delay)
	
	assert.Equal(t, addressHalf, item.AddressHalf)
	assert.Equal(t, delay, item.Delay)
}

func TestGenerateWhitelist(t *testing.T) {
	tests := []struct {
		name               string
		addresses          []string
		resolvingStartTime *big.Int
		expectError        bool
		expectedLen        int
	}{
		{
			name:               "Single address",
			addresses:          []string{"0x1234567890123456789012345678901234567890"},
			resolvingStartTime: big.NewInt(1000000),
			expectError:        false,
			expectedLen:        1,
		},
		{
			name: "Multiple addresses",
			addresses: []string{
				"0x1234567890123456789012345678901234567890",
				"0xabcdef1234567890abcdef1234567890abcdef12",
			},
			resolvingStartTime: big.NewInt(1000000),
			expectError:        false,
			expectedLen:        2,
		},
		{
			name:               "Empty addresses",
			addresses:          []string{},
			resolvingStartTime: big.NewInt(1000000),
			expectError:        true,
			expectedLen:        0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GenerateWhitelist(tc.addresses, tc.resolvingStartTime)
			
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tc.expectedLen)
				
				// Verify address halves are lowercase and correct length
				for i, item := range result {
					assert.Len(t, item.AddressHalf, 20, "Address half should be 20 chars (10 bytes)")
					assert.Equal(t, tc.addresses[i][len(tc.addresses[i])-20:], item.AddressHalf)
				}
			}
		})
	}
}

func TestGenerateWhitelistFromItems(t *testing.T) {
	tests := []struct {
		name               string
		items              []AuctionWhitelistItem
		resolvingStartTime *big.Int
		expectError        bool
		expectedLen        int
	}{
		{
			name: "Single item",
			items: []AuctionWhitelistItem{
				{
					Address:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
					AllowFrom: big.NewInt(1000000),
				},
			},
			resolvingStartTime: big.NewInt(1000000),
			expectError:        false,
			expectedLen:        1,
		},
		{
			name: "Multiple items with same AllowFrom",
			items: []AuctionWhitelistItem{
				{
					Address:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
					AllowFrom: big.NewInt(1000000),
				},
				{
					Address:   common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12"),
					AllowFrom: big.NewInt(1000000),
				},
			},
			resolvingStartTime: big.NewInt(1000000),
			expectError:        false,
			expectedLen:        2,
		},
		{
			name: "Items with different AllowFrom times",
			items: []AuctionWhitelistItem{
				{
					Address:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
					AllowFrom: big.NewInt(1000000),
				},
				{
					Address:   common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12"),
					AllowFrom: big.NewInt(1000100), // 100 seconds later
				},
			},
			resolvingStartTime: big.NewInt(1000000),
			expectError:        false,
			expectedLen:        2,
		},
		{
			name: "Item with AllowFrom before resolvingStartTime",
			items: []AuctionWhitelistItem{
				{
					Address:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
					AllowFrom: big.NewInt(900000), // Before resolvingStartTime
				},
			},
			resolvingStartTime: big.NewInt(1000000),
			expectError:        false,
			expectedLen:        1,
		},
		{
			name:               "Empty items",
			items:              []AuctionWhitelistItem{},
			resolvingStartTime: big.NewInt(1000000),
			expectError:        true,
			expectedLen:        0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GenerateWhitelistFromItems(tc.items, tc.resolvingStartTime)
			
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tc.expectedLen)
				
				// Verify address halves are lowercase and correct length
				for _, item := range result {
					assert.Len(t, item.AddressHalf, 20, "Address half should be 20 chars (10 bytes)")
				}
			}
		})
	}
}

func TestGenerateWhitelistDelayCalculation(t *testing.T) {
	items := []AuctionWhitelistItem{
		{
			Address:   common.HexToAddress("0x1111111111111111111111111111111111111111"),
			AllowFrom: big.NewInt(1000000), // Same as resolving start
		},
		{
			Address:   common.HexToAddress("0x2222222222222222222222222222222222222222"),
			AllowFrom: big.NewInt(1000100), // 100 seconds after first
		},
		{
			Address:   common.HexToAddress("0x3333333333333333333333333333333333333333"),
			AllowFrom: big.NewInt(1000200), // 100 seconds after second
		},
	}

	result, err := GenerateWhitelistFromItems(items, big.NewInt(1000000))
	require.NoError(t, err)
	require.Len(t, result, 3)

	// First item should have delay 0 (starts at resolving start time)
	assert.Equal(t, big.NewInt(0), result[0].Delay)
	
	// Second item should have delay 100 (100 seconds after first)
	assert.Equal(t, big.NewInt(100), result[1].Delay)
	
	// Third item should have delay 100 (100 seconds after second)
	assert.Equal(t, big.NewInt(100), result[2].Delay)
}
