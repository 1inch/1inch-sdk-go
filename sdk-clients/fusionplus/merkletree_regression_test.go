package fusionplus

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// genSecrets returns n distinct 32-byte hex secrets. The Merkle leaves derived
// from these are Keccak digests, so they land in an order that is not sorted,
// which is the condition that triggers the in-place-sort regression.
func genSecrets(n int) []string {
	secrets := make([]string, n)
	for i := 0; i < n; i++ {
		secrets[i] = fmt.Sprintf("0x%064x", (i+1)*0x1111)
	}
	return secrets
}

// merkleRoot builds the tree over a private copy and returns the root as the
// SDK stores it (a "0x" hex string).
func merkleRoot(leaves []string) string {
	cp := make([]string, len(leaves))
	copy(cp, leaves)
	return NewMerkleTree(cp).tree[0]
}

// proofAuthenticates walks the proof from the given leaf using the same
// sibling-sorted hashing the tree uses, and reports whether it reaches root.
func proofAuthenticates(leaf string, proof []string, root string) bool {
	cur := hexutil.MustDecode(leaf)
	for _, sibling := range proof {
		cur = Keccak256SortedHash(cur, hexutil.MustDecode(sibling))
	}
	return fmt.Sprintf("0x%x", cur) == root
}

// TestNewMerkleTree_DoesNotMutateArgument locks in the fix: building a tree
// must not reorder the caller's slice.
func TestNewMerkleTree_DoesNotMutateArgument(t *testing.T) {
	tests := []struct {
		name   string
		leaves []string
	}{
		{
			name: "two leaves in descending order",
			leaves: []string{
				"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			},
		},
		{
			name:   "four secret-derived leaves",
			leaves: mustLeaves(t, genSecrets(4)),
		},
		{
			name:   "five secret-derived leaves",
			leaves: mustLeaves(t, genSecrets(5)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			before := make([]string, len(tc.leaves))
			copy(before, tc.leaves)

			_ = NewMerkleTree(tc.leaves)

			assert.True(t, reflect.DeepEqual(before, tc.leaves),
				"NewMerkleTree must not reorder the caller's slice")
		})
	}
}

// TestGetProof_AuthenticatesRequestedLeaf is the core regression. For every
// index, and across repeated calls over the same slice, the returned proof must
// authenticate the leaf for THAT index and no other. Before the fix, the
// in-place sort permuted the slice, so later indices resolved to the wrong leaf.
func TestGetProof_AuthenticatesRequestedLeaf(t *testing.T) {
	tests := []struct {
		name    string
		secrets []string
	}{
		{name: "three secrets", secrets: genSecrets(3)},
		{name: "four secrets", secrets: genSecrets(4)},
		{name: "five secrets", secrets: genSecrets(5)},
		{name: "eight secrets", secrets: genSecrets(8)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			leaves := mustLeaves(t, tc.secrets)
			root := merkleRoot(leaves)

			for i := range leaves {
				proof, err := GetProof(leaves, i)
				require.NoError(t, err)

				assert.True(t, proofAuthenticates(leaves[i], proof, root),
					"proof for index %d must authenticate leaf %d", i, i)

				for j := range leaves {
					if j == i {
						continue
					}
					assert.False(t, proofAuthenticates(leaves[j], proof, root),
						"proof for index %d must not authenticate leaf %d", i, j)
				}
			}
		})
	}
}

// TestForMultipleFills_DoesNotMutateArgument covers the reachable maker path:
// building the hashlock must not reorder the caller's slice, so a later GetProof
// on the same slice still maps index N to secret N.
func TestForMultipleFills_DoesNotMutateArgument(t *testing.T) {
	tests := []struct {
		name    string
		secrets []string
	}{
		{name: "three secrets", secrets: genSecrets(3)},
		{name: "four secrets", secrets: genSecrets(4)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			leaves := mustLeaves(t, tc.secrets)
			before := make([]string, len(leaves))
			copy(before, leaves)
			root := merkleRoot(leaves)

			_, err := ForMultipleFills(leaves)
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(before, leaves),
				"ForMultipleFills must not reorder the caller's slice")

			// The first proof after building the hashlock must still be correct.
			proof, err := GetProof(leaves, 0)
			require.NoError(t, err)
			assert.True(t, proofAuthenticates(leaves[0], proof, root),
				"GetProof(0) after ForMultipleFills must authenticate leaf 0")
		})
	}
}

func mustLeaves(t *testing.T, secrets []string) []string {
	t.Helper()
	leaves, err := GetMerkleLeaves(secrets)
	require.NoError(t, err)
	return leaves
}
