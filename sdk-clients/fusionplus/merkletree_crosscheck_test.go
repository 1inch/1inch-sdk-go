package fusionplus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The multiple-fill hashlock, its leaves, and its proofs must match the
// canonical 1inch cross-chain reference, the JS SDK @1inch/cross-chain-sdk
// (HashLock in src/domains/hash-lock/hash-lock.ts). testdata holds golden
// vectors generated from that SDK at version 2.2.2 for the deterministic
// secrets fmt.Sprintf("0x%064x", (i+1)*0x1111). Regenerate with:
//
//	const {HashLock} = require('@1inch/cross-chain-sdk')
//	const leaves = HashLock.getMerkleLeaves(secrets)
//	HashLock.forMultipleFills(leaves).toString(); HashLock.getProof(leaves, i)
type jsHashlockVector struct {
	Secrets  []string   `json:"secrets"`
	Leaves   []string   `json:"leaves"`
	Hashlock string     `json:"hashlock"`
	Proofs   [][]string `json:"proofs"`
}

func loadJSVectors(t *testing.T) map[string]jsHashlockVector {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "jsref_hashlock_vectors.json"))
	require.NoError(t, err)
	var vectors map[string]jsHashlockVector
	require.NoError(t, json.Unmarshal(raw, &vectors))
	require.NotEmpty(t, vectors)
	return vectors
}

// TestHashLock_MatchesJSReferenceSDK is the canonical cross-SDK check. For each
// vector, the Go leaves, hashlock, proofs, and parts count must match the JS
// SDK byte-for-byte, and every proof must authenticate the leaf for its index.
func TestHashLock_MatchesJSReferenceSDK(t *testing.T) {
	for name, vec := range loadJSVectors(t) {
		t.Run("n="+name, func(t *testing.T) {
			leaves, err := GetMerkleLeaves(vec.Secrets)
			require.NoError(t, err)
			assert.Equal(t, vec.Leaves, leaves, "leaves must match the JS SDK")

			hashLock, err := ForMultipleFills(leaves)
			require.NoError(t, err)
			assert.Equal(t, vec.Hashlock, hashLock.Value, "hashlock must match the JS SDK")

			assert.Equal(t, uint64(len(leaves)-1), hashLock.GetPartsCount(),
				"parts count is len(leaves)-1")

			root := merkleRoot(leaves)
			for i := range leaves {
				proof, err := GetProof(leaves, i)
				require.NoError(t, err)
				assert.Equal(t, vec.Proofs[i], proof,
					"proof for index %d must match the JS SDK", i)
				assert.True(t, proofAuthenticates(leaves[i], proof, root),
					"proof for index %d must authenticate leaf %d", i, i)
			}
		})
	}
}

// TestForMultipleFills_RequiresLeavesNotSecrets pins the example-and-canary
// defect and its fix. Passing raw secrets (the old example path) builds the
// tree over the wrong values and yields a hashlock that does NOT match the
// reference. The correct path, over GetMerkleLeaves, matches the JS SDK.
func TestForMultipleFills_RequiresLeavesNotSecrets(t *testing.T) {
	for name, vec := range loadJSVectors(t) {
		t.Run("n="+name, func(t *testing.T) {
			// Pre-fix behavior: raw secrets in, wrong hashlock out.
			wrong, err := ForMultipleFills(vec.Secrets)
			require.NoError(t, err)
			assert.NotEqual(t, vec.Hashlock, wrong.Value,
				"raw secrets must not produce the reference hashlock")

			// Post-fix behavior: leaves in, reference hashlock out.
			leaves, err := GetMerkleLeaves(vec.Secrets)
			require.NoError(t, err)
			correct, err := ForMultipleFills(leaves)
			require.NoError(t, err)
			assert.Equal(t, vec.Hashlock, correct.Value,
				"leaves must produce the reference hashlock")
		})
	}
}

// TestGetPartsCount_MatchesReference checks the parts count that a maker reads
// back from a multiple-fill hashlock equals the count the JS SDK packs in.
func TestGetPartsCount_MatchesReference(t *testing.T) {
	tests := []struct {
		name    string
		secrets []string
	}{
		{name: "three", secrets: genSecrets(3)},
		{name: "four", secrets: genSecrets(4)},
		{name: "eight", secrets: genSecrets(8)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			leaves, err := GetMerkleLeaves(tc.secrets)
			require.NoError(t, err)
			hashLock, err := ForMultipleFills(leaves)
			require.NoError(t, err)
			assert.Equal(t, uint64(len(tc.secrets)-1), hashLock.GetPartsCount(),
				fmt.Sprintf("parts count for %d secrets", len(tc.secrets)))
		})
	}
}
