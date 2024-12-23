package fusionplus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashLockForSingleFill(t *testing.T) {
	result, err := ForSingleFill("0x531d1d2d7a594f1c7e413b074c7b693161486b5c495d457748144a01795c6a45")
	require.NoError(t, err)
	expected := "0x9f65fdcf781d4320c2dde70da02a1fe916d595dc1817149cc4758fd6a4bfd830"
	assert.Equal(t, expected, result.Value, "HashLock for single fill is incorrect")
}

func TestHashLockForMultipleFills(t *testing.T) {
	secrets := []string{
		"0x531d1d2d7a594f1c7e413b074c7b693161486b5c495d457748144a01795c6a45",
		"0x657812136b5000651d5e18516d764b5e661a681c760d3c3c4c15751020757823",
		"0x62071a322351281f04756576270c362a6e5b395e3b0f68027f231141555c3d43",
	}
	leaves, err := GetMerkleLeaves(secrets)
	require.NoError(t, err)
	result, err := ForMultipleFills(leaves)
	require.NoError(t, err)
	expected := "0x000292766d9172e4b4983ee4d4b6d511cdbcbef175c7e3e1b1554d513e1ab724"
	assert.Equal(t, expected, result.Value, "HashLock for multiple fills is incorrect")
}

func TestHashLockIsBytes32(t *testing.T) {
	secrets := []string{
		"0x6466643931343237333333313437633162386632316365646666323931643738",
		"0x3131353932633266343034343466363562333230313837353438356463616130",
		"0x6634376135663837653765303462346261616566383430303662303336386635",
	}
	leaves, err := GetMerkleLeaves(secrets)
	require.NoError(t, err)
	result, err := ForMultipleFills(leaves)
	require.NoError(t, err)
	bytes := getBytesCount(result.Value)
	assert.Equal(t, 32, bytes, "HashLock result length is not 32 bytes")
}

func TestHashLockGetProof(t *testing.T) {
	secrets := []string{
		"0x6466643931343237333333313437633162386632316365646666323931643738",
		"0x3131353932633266343034343466363562333230313837353438356463616130",
		"0x6634376135663837653765303462346261616566383430303662303336386635",
	}
	leaves, err := GetMerkleLeaves(secrets)
	require.NoError(t, err)
	result, err := GetProof(leaves, 0)
	require.NoError(t, err)
	expected := []string{
		"0x540daf363747246d40b31da95b3ef1c1497e22e9a56b70d117c835839822c95f",
	}
	assert.Equal(t, expected, result, "HashLock proof is incorrect")
}

func TestHashLockGetProof2(t *testing.T) {
	secrets := []string{
		"0x531d1d2d7a594f1c7e413b074c7b693161486b5c495d457748144a01795c6a45",
		"0x657812136b5000651d5e18516d764b5e661a681c760d3c3c4c15751020757823",
		"0x62071a322351281f04756576270c362a6e5b395e3b0f68027f231141555c3d43",
	}
	leaves, err := GetMerkleLeaves(secrets)
	require.NoError(t, err)
	result, err := GetProof(leaves, 0)
	require.NoError(t, err)
	expected := []string{
		"0xb19c79aa34d58e459ce8119c301a24f7a01b8080ced7f3d608093e9e67624729",
	}
	assert.Equal(t, expected, result, "HashLock proof is incorrect")
}
