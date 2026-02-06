package bytesiterator

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bigIntFromHex(hexStr string) (*big.Int, error) {
	b, err := hexutil.Decode(hexStr)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(b), nil
}

func uint32FromHex(hexStr string) (uint32, error) {
	b, err := hexutil.Decode(hexStr)
	if err != nil {
		return 0, err
	}
	if len(b) > 4 {
		return 0, fmt.Errorf("more than 4 bytes: %v", hexStr)
	}
	var val uint32
	for _, bb := range b {
		val = (val << 8) | uint32(bb)
	}
	return val, nil
}

func TestUint32FromHex(t *testing.T) {
	tests := []struct {
		name        string
		hexStr      string
		expected    uint32
		expectError bool
	}{
		{
			name:        "Empty string",
			hexStr:      "0x",
			expected:    0,
			expectError: false,
		},
		{
			name:        "Single byte zero",
			hexStr:      "0x00",
			expected:    0,
			expectError: false,
		},
		{
			name:        "Two bytes 0x0001",
			hexStr:      "0x0001",
			expected:    1,
			expectError: false,
		},
		{
			name:        "Four bytes 0x01020304",
			hexStr:      "0x01020304",
			expected:    16909060,
			expectError: false,
		},
		{
			name:        "More than four bytes 0x01020304AA",
			hexStr:      "0x01020304AA",
			expectError: true,
		},
		{
			name:        "Invalid hex characters",
			hexStr:      "0xGG",
			expected:    0,
			expectError: true,
		},
		{
			name:        "Odd length hex string",
			hexStr:      "0x123",
			expected:    0,
			expectError: true,
		},
		{
			name:        "Four bytes all zeros",
			hexStr:      "0x00000000",
			expected:    0,
			expectError: false,
		},
		{
			name:        "Four bytes all ones",
			hexStr:      "0xFFFFFFFF",
			expected:    4294967295,
			expectError: false,
		},
		{
			name:        "Three bytes 0x0000FF",
			hexStr:      "0x0000FF",
			expected:    255,
			expectError: false,
		},
		{
			name:        "Four bytes with leading zero 0x000000FF",
			hexStr:      "0x000000FF",
			expected:    255,
			expectError: false,
		},
		{
			name:        "Four bytes 0x01000000",
			hexStr:      "0x01000000",
			expected:    16777216,
			expectError: false,
		},
		{
			name:        "Four bytes 0x0A0B0C0D",
			hexStr:      "0x0A0B0C0D",
			expected:    168496141,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := uint32FromHex(tt.hexStr)
			if tt.expectError {
				assert.Error(t, err, "Expected an error for input: %s", tt.hexStr)
			} else {
				require.NoError(t, err, "Did not expect an error for input: %s", tt.hexStr)
				assert.Equal(t, tt.expected, val, "Mismatch for input: %s", tt.hexStr)
			}
		})
	}
}

func TestBytesIter_NextByte(t *testing.T) {
	tests := []struct {
		name        string
		hexData     string
		readCount   int
		expectedHex string
		expectError bool
	}{
		{
			name:        "Exact Bytes",
			hexData:     "0x010203",
			readCount:   3,
			expectedHex: "0x010203",
			expectError: false,
		},
		{
			name:        "Read Beyond",
			hexData:     "0x0102",
			readCount:   3,
			expectError: true,
		},
		{
			name:        "Single Byte Exact",
			hexData:     "0xFF",
			readCount:   1,
			expectedHex: "0xFF",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			var result []byte
			var readErr error
			for i := 0; i < tt.readCount; i++ {
				b, err := iter.NextByte()
				if err != nil {
					readErr = err
					break
				}
				result = append(result, b)
			}

			if tt.expectError {
				assert.Error(t, readErr)
			} else {
				expected, err := hexutil.Decode(tt.expectedHex)
				require.NoError(t, err)
				assert.NoError(t, readErr)
				assert.Equal(t, expected, result)
			}
		})
	}
}

func TestBytesIter_NextUint16(t *testing.T) {
	tests := []struct {
		name        string
		hexData     string
		expectedHex string
		expectError bool
	}{
		{
			name:        "Valid",
			hexData:     "0x00FFAB",
			expectedHex: "0x00FF",
			expectError: false,
		},
		{
			name:        "Insufficient Bytes",
			hexData:     "0x00",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			val, err := iter.NextUint16()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				expectedVal, err := bigIntFromHex(tt.expectedHex)
				require.NoError(t, err)
				assert.Equal(t, expectedVal, val)
			}
		})
	}
}

func TestBytesIter_NextUint24(t *testing.T) {
	tests := []struct {
		name        string
		hexData     string
		expectedHex string
		expectError bool
	}{
		{
			name:        "Valid",
			hexData:     "0x01020304",
			expectedHex: "0x010203",
			expectError: false,
		},
		{
			name:        "Insufficient Bytes",
			hexData:     "0x01",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			val, err := iter.NextUint24()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				expectedVal, err := uint32FromHex(tt.expectedHex)
				require.NoError(t, err)
				assert.Equal(t, expectedVal, val)
			}
		})
	}
}

func TestBytesIter_NextUint32(t *testing.T) {
	tests := []struct {
		name        string
		hexData     string
		expectedHex string
		expectError bool
	}{
		{
			name:        "Valid",
			hexData:     "0xFFEEDDCCAB",
			expectedHex: "0xFFEEDDCC",
			expectError: false,
		},
		{
			name:        "Insufficient Bytes",
			hexData:     "0xFFEE",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			val, err := iter.NextUint32()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				expectedVal, err := bigIntFromHex(tt.expectedHex)
				require.NoError(t, err)
				assert.Equal(t, expectedVal, val)
			}
		})
	}
}

func TestBytesIter_NextUint160(t *testing.T) {
	twentyBytes := strings.Repeat("FF", 20) // 20 bytes of 0xFF
	tests := []struct {
		name        string
		hexData     string
		expectError bool
	}{
		{
			name:    "Valid",
			hexData: "0x" + twentyBytes,
		},
		{
			name:        "Insufficient Bytes",
			hexData:     "0xFFEE",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			val, err := iter.NextUint160()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				expectedVal, err := bigIntFromHex(tt.hexData)
				require.NoError(t, err)
				assert.Equal(t, expectedVal, val)
			}
		})
	}
}

func TestBytesIter_NextUint256(t *testing.T) {
	thirtyTwoBytes := strings.Repeat("AA", 32) // 32 bytes of 0xAA
	tests := []struct {
		name        string
		hexData     string
		expectError bool
	}{
		{
			name:    "Valid",
			hexData: "0x" + thirtyTwoBytes,
		},
		{
			name:        "Insufficient Bytes",
			hexData:     "0xAA",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			val, err := iter.NextUint256()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				expectedVal, err := bigIntFromHex(tt.hexData)
				require.NoError(t, err)
				assert.Equal(t, expectedVal, val)
			}
		})
	}
}

func TestBytesIter_NextBytes(t *testing.T) {
	tests := []struct {
		name        string
		hexData     string
		readSize    int
		expectedHex string
		expectError bool
	}{
		{
			name:        "Valid Read",
			hexData:     "0x10203040",
			readSize:    2,
			expectedHex: "0x1020",
			expectError: false,
		},
		{
			name:        "Insufficient Data",
			hexData:     "0x10",
			readSize:    2,
			expectError: true,
		},
		{
			name:        "Negative Length",
			hexData:     "0x1020",
			readSize:    -1,
			expectError: true,
		},
		{
			name:        "Exact Length Bytes",
			hexData:     "0xABCDEE",
			readSize:    3,
			expectedHex: "0xABCDEE",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			val, err := iter.NextBytes(tt.readSize)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				expected, err := hexutil.Decode(tt.expectedHex)
				require.NoError(t, err)
				assert.Equal(t, expected, val)
			}
		})
	}
}

func TestBytesIter_NextString(t *testing.T) {
	tests := []struct {
		name        string
		hexData     string
		readSize    int
		expectedHex string
		expectError bool
	}{
		{
			name:        "Valid String",
			hexData:     "0x48656C6C6F",
			readSize:    5,
			expectedHex: "0x48656C6C6F",
			expectError: false,
		},
		{
			name:        "Insufficient Data",
			hexData:     "0x48",
			readSize:    5,
			expectError: true,
		},
		{
			name:        "Exact Length String",
			hexData:     "0x416263",
			readSize:    3,
			expectedHex: "0x416263",
			expectError: false,
		},
		{
			name:        "One Byte String",
			hexData:     "0x41",
			readSize:    1,
			expectedHex: "0x41",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			strVal, err := iter.NextString(tt.readSize)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				expectedBytes, err := hexutil.Decode(tt.expectedHex)
				require.NoError(t, err)
				expectedStr := string(expectedBytes)
				assert.Equal(t, expectedStr, strVal)
			}
		})
	}
}

func TestBytesIter_Rest(t *testing.T) {
	tests := []struct {
		name        string
		hexData     string
		readSize    int
		expectedHex string
	}{
		{
			name:        "Read Part, Then Rest",
			hexData:     "0x010203",
			readSize:    1,
			expectedHex: "0x0203",
		},
		{
			name:        "Read All, Then Rest",
			hexData:     "0x0102",
			readSize:    2,
			expectedHex: "0x", // no rest
		},
		{
			name:        "No Read, Return All",
			hexData:     "0x050607",
			readSize:    0,
			expectedHex: "0x050607",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			if tt.readSize > 0 {
				_, _ = iter.NextBytes(tt.readSize)
			}

			rest, err := iter.Rest()
			require.NoError(t, err)

			if tt.expectedHex == "0x" {
				// Expect no remaining bytes; rest can be nil or empty
				assert.Empty(t, rest, "Expected no remaining bytes, but some were found")
			} else {
				expected, err := hexutil.Decode(tt.expectedHex)
				require.NoError(t, err)
				assert.Equal(t, expected, rest)
			}

			// Subsequent calls should yield empty
			rest, err = iter.Rest()
			require.NoError(t, err)
			assert.Empty(t, rest, "Expected no remaining bytes on subsequent Rest call")
		})
	}
}

func TestBytesIter_IsEmpty(t *testing.T) {
	tests := []struct {
		name          string
		hexData       string
		reads         int
		expectedEmpty bool
	}{
		{
			name:          "Initially Not Empty",
			hexData:       "0x0102",
			reads:         0,
			expectedEmpty: false,
		},
		{
			name:          "Partially Read",
			hexData:       "0x0102",
			reads:         1,
			expectedEmpty: false,
		},
		{
			name:          "Fully Read",
			hexData:       "0x0102",
			reads:         2,
			expectedEmpty: true,
		},
		{
			name:          "One Byte Only",
			hexData:       "0xFF",
			reads:         1,
			expectedEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hexutil.Decode(tt.hexData)
			require.NoError(t, err)
			iter := New(data)

			for i := 0; i < tt.reads; i++ {
				_, _ = iter.NextByte()
			}

			assert.Equal(t, tt.expectedEmpty, iter.IsEmpty())
		})
	}
}
