package bytesbuilder

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	bb := New()
	require.NotNil(t, bb)
	assert.Equal(t, "", bb.AsHex())
}

func TestBytesBuilder_AddUint256(t *testing.T) {
	tests := []struct {
		name     string
		value    *big.Int
		expected string
	}{
		{
			name:     "Zero value",
			value:    big.NewInt(0),
			expected: "0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:     "Small value - 1",
			value:    big.NewInt(1),
			expected: "0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			name:     "Small value - 255",
			value:    big.NewInt(255),
			expected: "00000000000000000000000000000000000000000000000000000000000000ff",
		},
		{
			name:     "Large value - 1 ETH in wei",
			value:    new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
			expected: "0000000000000000000000000000000000000000000000000de0b6b3a7640000",
		},
		{
			name:     "Max uint64",
			value:    new(big.Int).SetUint64(^uint64(0)),
			expected: "000000000000000000000000000000000000000000000000ffffffffffffffff",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			bb.AddUint256(tc.value)
			assert.Equal(t, tc.expected, bb.AsHex())
		})
	}
}

func TestBytesBuilder_AddUint24(t *testing.T) {
	tests := []struct {
		name     string
		value    *big.Int
		expected string
	}{
		{
			name:     "Zero value",
			value:    big.NewInt(0),
			expected: "000000",
		},
		{
			name:     "Small value - 1",
			value:    big.NewInt(1),
			expected: "000001",
		},
		{
			name:     "Max uint24 - 16777215",
			value:    big.NewInt(16777215),
			expected: "ffffff",
		},
		{
			name:     "Middle value - 256",
			value:    big.NewInt(256),
			expected: "000100",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			bb.AddUint24(tc.value)
			assert.Equal(t, tc.expected, bb.AsHex())
		})
	}
}

func TestBytesBuilder_AddUint32(t *testing.T) {
	tests := []struct {
		name     string
		value    *big.Int
		expected string
	}{
		{
			name:     "Zero value",
			value:    big.NewInt(0),
			expected: "00000000",
		},
		{
			name:     "Small value - 1",
			value:    big.NewInt(1),
			expected: "00000001",
		},
		{
			name:     "Max uint32 - 4294967295",
			value:    big.NewInt(4294967295),
			expected: "ffffffff",
		},
		{
			name:     "Middle value - 65536",
			value:    big.NewInt(65536),
			expected: "00010000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			bb.AddUint32(tc.value)
			assert.Equal(t, tc.expected, bb.AsHex())
		})
	}
}

func TestBytesBuilder_AddUint16(t *testing.T) {
	tests := []struct {
		name     string
		value    *big.Int
		expected string
	}{
		{
			name:     "Zero value",
			value:    big.NewInt(0),
			expected: "0000",
		},
		{
			name:     "Small value - 1",
			value:    big.NewInt(1),
			expected: "0001",
		},
		{
			name:     "Max uint16 - 65535",
			value:    big.NewInt(65535),
			expected: "ffff",
		},
		{
			name:     "Middle value - 256",
			value:    big.NewInt(256),
			expected: "0100",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			bb.AddUint16(tc.value)
			assert.Equal(t, tc.expected, bb.AsHex())
		})
	}
}

func TestBytesBuilder_AddUint8(t *testing.T) {
	tests := []struct {
		name     string
		value    uint8
		expected string
	}{
		{
			name:     "Zero value",
			value:    0,
			expected: "00",
		},
		{
			name:     "Small value - 1",
			value:    1,
			expected: "01",
		},
		{
			name:     "Max uint8 - 255",
			value:    255,
			expected: "ff",
		},
		{
			name:     "Middle value - 128",
			value:    128,
			expected: "80",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			bb.AddUint8(tc.value)
			assert.Equal(t, tc.expected, bb.AsHex())
		})
	}
}

func TestBytesBuilder_AddAddress(t *testing.T) {
	tests := []struct {
		name     string
		address  common.Address
		expected string
	}{
		{
			name:     "Zero address",
			address:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
			expected: "0000000000000000000000000000000000000000",
		},
		{
			name:     "DAI address",
			address:  common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
			expected: "6b175474e89094c44da98b954eedeac495271d0f",
		},
		{
			name:     "WETH address",
			address:  common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			expected: "c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			bb.AddAddress(tc.address)
			assert.Equal(t, tc.expected, bb.AsHex())
		})
	}
}

func TestBytesBuilder_AddBytes(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		expected    string
		expectError bool
	}{
		{
			name:        "Valid hex with 0x prefix",
			data:        "0x1234567890abcdef",
			expected:    "1234567890abcdef",
			expectError: false,
		},
		{
			name:        "Valid hex without 0x prefix",
			data:        "1234567890abcdef",
			expected:    "1234567890abcdef",
			expectError: false,
		},
		{
			name:        "Empty hex",
			data:        "0x",
			expected:    "",
			expectError: false,
		},
		{
			name:        "Empty string",
			data:        "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "Invalid hex - odd length",
			data:        "0x123",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid hex - non-hex characters",
			data:        "0xGHIJ",
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			err := bb.AddBytes(tc.data)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, bb.AsHex())
			}
		})
	}
}

func TestBytesBuilder_AddNativeUint16(t *testing.T) {
	tests := []struct {
		name     string
		value    uint16
		expected string
	}{
		{
			name:     "Zero value",
			value:    0,
			expected: "0000",
		},
		{
			name:     "Small value - 1",
			value:    1,
			expected: "0001",
		},
		{
			name:     "Max uint16 - 65535",
			value:    65535,
			expected: "ffff",
		},
		{
			name:     "Middle value - 256",
			value:    256,
			expected: "0100",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			bb.AddNativeUint16(tc.value)
			assert.Equal(t, tc.expected, bb.AsHex())
		})
	}
}

func TestBytesBuilder_AddNativeUint24(t *testing.T) {
	tests := []struct {
		name     string
		value    uint32
		expected string
	}{
		{
			name:     "Zero value",
			value:    0,
			expected: "000000",
		},
		{
			name:     "Small value - 1",
			value:    1,
			expected: "000001",
		},
		{
			name:     "Max uint24 - 16777215",
			value:    16777215,
			expected: "ffffff",
		},
		{
			name:     "Middle value - 256",
			value:    256,
			expected: "000100",
		},
		{
			name:     "Value 10000",
			value:    10000,
			expected: "002710",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			bb.AddNativeUint24(tc.value)
			assert.Equal(t, tc.expected, bb.AsHex())
		})
	}
}

func TestBytesBuilder_AddNativeUint32(t *testing.T) {
	tests := []struct {
		name     string
		value    uint32
		expected string
	}{
		{
			name:     "Zero value",
			value:    0,
			expected: "00000000",
		},
		{
			name:     "Small value - 1",
			value:    1,
			expected: "00000001",
		},
		{
			name:     "Max uint32 - 4294967295",
			value:    4294967295,
			expected: "ffffffff",
		},
		{
			name:     "Middle value - 65536",
			value:    65536,
			expected: "00010000",
		},
		{
			name:     "Timestamp - 1673548149",
			value:    1673548149,
			expected: "63c05175",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bb := New()
			bb.AddNativeUint32(tc.value)
			assert.Equal(t, tc.expected, bb.AsHex())
		})
	}
}

func TestBytesBuilder_AsBytes(t *testing.T) {
	bb := New()
	bb.AddUint8(0x01)
	bb.AddUint8(0x02)
	bb.AddUint8(0x03)

	result := bb.AsBytes()
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, result)
}

func TestBytesBuilder_AsBytesEmpty(t *testing.T) {
	bb := New()
	result := bb.AsBytes()
	assert.Equal(t, []byte{}, result)
}

func TestBytesBuilder_NativeMatchesBigInt(t *testing.T) {
	// Verify that native methods produce identical output to big.Int methods
	tests := []struct {
		name  string
		value uint32
	}{
		{"Zero", 0},
		{"One", 1},
		{"Thousand", 1000},
		{"Large", 1673548149},
		{"Max24", 16777215},
		{"Max32", 4294967295},
	}

	for _, tc := range tests {
		t.Run(tc.name+"_uint32", func(t *testing.T) {
			native := New()
			native.AddNativeUint32(tc.value)

			bigint := New()
			bigint.AddUint32(big.NewInt(int64(tc.value)))

			assert.Equal(t, bigint.AsHex(), native.AsHex())
		})

		if tc.value <= 16777215 {
			t.Run(tc.name+"_uint24", func(t *testing.T) {
				native := New()
				native.AddNativeUint24(tc.value)

				bigint := New()
				bigint.AddUint24(big.NewInt(int64(tc.value)))

				assert.Equal(t, bigint.AsHex(), native.AsHex())
			})
		}

		if tc.value <= 65535 {
			t.Run(tc.name+"_uint16", func(t *testing.T) {
				native := New()
				native.AddNativeUint16(uint16(tc.value))

				bigint := New()
				bigint.AddUint16(big.NewInt(int64(tc.value)))

				assert.Equal(t, bigint.AsHex(), native.AsHex())
			})
		}
	}
}

func TestBytesBuilder_ChainedOperations(t *testing.T) {
	bb := New()

	// Add a uint8 (1 byte)
	bb.AddUint8(0x01)

	// Add a uint16 (2 bytes)
	bb.AddUint16(big.NewInt(0x0203))

	// Add a uint24 (3 bytes)
	bb.AddUint24(big.NewInt(0x040506))

	// Add a uint32 (4 bytes)
	bb.AddUint32(big.NewInt(0x0708090a))

	expected := "01" + "0203" + "040506" + "0708090a"
	assert.Equal(t, expected, bb.AsHex())
}

func TestBytesBuilder_ComplexBuild(t *testing.T) {
	bb := New()

	// Simulate building transaction data
	address := common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")
	amount := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 1 ETH

	bb.AddAddress(address)
	bb.AddUint256(amount)

	result := bb.AsHex()

	// Address should be 20 bytes (40 hex chars)
	// Uint256 should be 32 bytes (64 hex chars)
	assert.Len(t, result, 104)

	// Verify address portion
	assert.Equal(t, "6b175474e89094c44da98b954eedeac495271d0f", result[:40])

	// Verify amount portion
	assert.Equal(t, "0000000000000000000000000000000000000000000000000de0b6b3a7640000", result[40:])
}
