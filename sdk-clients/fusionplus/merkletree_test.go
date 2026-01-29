package fusionplus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeTree_SingleLeaf(t *testing.T) {
	leaves := []string{"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}

	tree := MakeTree(leaves)
	require.NotNil(t, tree)
	assert.Len(t, tree.leaves, 1)
	assert.Len(t, tree.tree, 1)
}

func TestMakeTree_TwoLeaves(t *testing.T) {
	leaves := []string{
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
	}

	tree := MakeTree(leaves)
	require.NotNil(t, tree)
	assert.Len(t, tree.leaves, 2)
	assert.Len(t, tree.tree, 3) // 2 leaves + 1 root
}

func TestMakeTree_FourLeaves(t *testing.T) {
	leaves := []string{
		"0x1111111111111111111111111111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333333333333333333333333333",
		"0x4444444444444444444444444444444444444444444444444444444444444444",
	}

	tree := MakeTree(leaves)
	require.NotNil(t, tree)
	assert.Len(t, tree.leaves, 4)
	assert.Len(t, tree.tree, 7) // 4 leaves + 2 intermediate + 1 root
}

func TestMakeTree_PreservesOriginalOrderInLeaves(t *testing.T) {
	leaves := []string{
		"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	tree := MakeTree(leaves)
	require.NotNil(t, tree)

	// The tree.leaves field stores the original unsorted order
	// This is used for GetProof to find the leaf by index in the original array
	assert.Len(t, tree.leaves, 2)
	// The tree itself is built with sorted leaves for consistent hashing
}

func TestGetProof_SingleLeaf(t *testing.T) {
	leaves := []string{"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}

	proof, err := GetProof(leaves, 0)
	require.NoError(t, err)
	assert.Empty(t, proof) // Single leaf has no proof (it's the root)
}

func TestGetProof_TwoLeaves(t *testing.T) {
	leaves := []string{
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
	}

	// Get proof for first leaf
	proof0, err := GetProof(leaves, 0)
	require.NoError(t, err)
	assert.Len(t, proof0, 1) // Should have one sibling in proof

	// Get proof for second leaf
	proof1, err := GetProof(leaves, 1)
	require.NoError(t, err)
	assert.Len(t, proof1, 1) // Should have one sibling in proof
}

func TestGetProof_FourLeaves(t *testing.T) {
	leaves := []string{
		"0x1111111111111111111111111111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333333333333333333333333333",
		"0x4444444444444444444444444444444444444444444444444444444444444444",
	}

	for i := 0; i < len(leaves); i++ {
		proof, err := GetProof(leaves, i)
		require.NoError(t, err)
		assert.Len(t, proof, 2) // Depth 2 tree should have 2 elements in proof
	}
}

func TestGetProof_IndexOutOfBounds(t *testing.T) {
	leaves := []string{
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	}

	// Test negative index
	_, err := GetProof(leaves, -1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "index out of bounds")

	// Test index >= len(leaves)
	_, err = GetProof(leaves, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "index out of bounds")
}

func TestGetSiblingIndex(t *testing.T) {
	// Index 1 (left child of root) -> sibling is 2
	idx, err := getSiblingIndex(1)
	require.NoError(t, err)
	assert.Equal(t, 2, idx)

	// Index 2 (right child of root) -> sibling is 1
	idx, err = getSiblingIndex(2)
	require.NoError(t, err)
	assert.Equal(t, 1, idx)

	// Index 3 -> sibling is 4
	idx, err = getSiblingIndex(3)
	require.NoError(t, err)
	assert.Equal(t, 4, idx)

	// Index 4 -> sibling is 3
	idx, err = getSiblingIndex(4)
	require.NoError(t, err)
	assert.Equal(t, 3, idx)
}

func TestGetSiblingIndex_Root_ReturnsError(t *testing.T) {
	_, err := getSiblingIndex(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "root has no siblings")
}

func TestParentIndex(t *testing.T) {
	tests := []struct {
		name        string
		index       int
		expected    int
		expectError bool
	}{
		{
			name:        "Index 1 -> parent 0",
			index:       1,
			expected:    0,
			expectError: false,
		},
		{
			name:        "Index 2 -> parent 0",
			index:       2,
			expected:    0,
			expectError: false,
		},
		{
			name:        "Index 3 -> parent 1",
			index:       3,
			expected:    1,
			expectError: false,
		},
		{
			name:        "Index 4 -> parent 1",
			index:       4,
			expected:    1,
			expectError: false,
		},
		{
			name:        "Index 5 -> parent 2",
			index:       5,
			expected:    2,
			expectError: false,
		},
		{
			name:        "Index 6 -> parent 2",
			index:       6,
			expected:    2,
			expectError: false,
		},
		{
			name:        "Root (index 0) has no parent",
			index:       0,
			expected:    0,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParentIndex(tc.index)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestLeftChildIndex(t *testing.T) {
	// leftChildIndex(0) = 2*0 + 1 = 1
	assert.Equal(t, 1, leftChildIndex(0))

	// leftChildIndex(1) = 2*1 + 1 = 3
	assert.Equal(t, 3, leftChildIndex(1))

	// leftChildIndex(2) = 2*2 + 1 = 5
	assert.Equal(t, 5, leftChildIndex(2))
}

func TestRightChildIndex(t *testing.T) {
	// rightChildIndex(0) = 2*0 + 2 = 2
	assert.Equal(t, 2, rightChildIndex(0))

	// rightChildIndex(1) = 2*1 + 2 = 4
	assert.Equal(t, 4, rightChildIndex(1))

	// rightChildIndex(2) = 2*2 + 2 = 6
	assert.Equal(t, 6, rightChildIndex(2))
}
