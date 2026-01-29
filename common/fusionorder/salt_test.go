package fusionorder

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeccak256Hash(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "String data",
			data: "hello world",
		},
		{
			name: "Struct data",
			data: struct {
				Field1 string
				Field2 int
			}{"test", 123},
		},
		{
			name: "Map data",
			data: map[string]int{"a": 1, "b": 2},
		},
		{
			name: "Array data",
			data: []int{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Keccak256Hash(tc.data)
			require.NoError(t, err)

			// Result should be a 256-bit (32 byte) hash
			require.NotNil(t, result)
			assert.Equal(t, 32, len(result.Bytes()), "Hash should be 32 bytes")

			// Same input should produce same output
			result2, err := Keccak256Hash(tc.data)
			require.NoError(t, err)
			assert.Equal(t, result, result2)
		})
	}
}

func TestKeccak256HashDifferentInputs(t *testing.T) {
	hash1, err := Keccak256Hash("hello")
	require.NoError(t, err)
	hash2, err := Keccak256Hash("world")
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "Different inputs should produce different hashes")
}

func TestGenerateSaltWithExtensionEmpty(t *testing.T) {
	salt, err := GenerateSaltWithExtension(nil, true)
	require.NoError(t, err)
	require.NotNil(t, salt)

	// When isEmpty is true, salt should be less than 2^96
	maxValue := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 96), big.NewInt(1))
	assert.True(t, salt.Cmp(maxValue) <= 0, "Empty extension salt should be <= 2^96 - 1")
	assert.True(t, salt.Sign() >= 0, "Salt should be non-negative")
}

func TestGenerateSaltWithExtensionNonEmpty(t *testing.T) {
	extensionHash := big.NewInt(0x123456789ABCDEF0)

	salt, err := GenerateSaltWithExtension(extensionHash, false)
	require.NoError(t, err)
	require.NotNil(t, salt)

	// When isEmpty is false, salt should incorporate the extension hash
	// The lowest 160 bits should match the lowest 160 bits of extensionHash
	uint160Max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
	saltLow160 := new(big.Int).And(salt, uint160Max)
	hashLow160 := new(big.Int).And(extensionHash, uint160Max)

	assert.Equal(t, hashLow160, saltLow160, "Low 160 bits should match extension hash")
}

func TestGenerateSaltWithExtensionRandomness(t *testing.T) {
	extensionHash := big.NewInt(12345)

	// Generate multiple salts and ensure they're different (random component)
	salts := make(map[string]bool)
	for i := 0; i < 10; i++ {
		salt, err := GenerateSaltWithExtension(extensionHash, false)
		require.NoError(t, err)
		salts[salt.String()] = true
	}

	// With high probability, we should get different salts
	// (allowing for some small chance of collision)
	assert.Greater(t, len(salts), 1, "Multiple salts should be different due to random component")
}

func TestGenerateSaltWithExtensionEmptyRandomness(t *testing.T) {
	// Generate multiple empty salts and ensure they're different (random)
	salts := make(map[string]bool)
	for i := 0; i < 10; i++ {
		salt, err := GenerateSaltWithExtension(nil, true)
		require.NoError(t, err)
		salts[salt.String()] = true
	}

	// With high probability, we should get different salts
	assert.Greater(t, len(salts), 1, "Multiple empty salts should be different")
}
