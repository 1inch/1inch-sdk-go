package fusionplus

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type MerkleTree struct {
	tree   []string
	leaves []string
}

func NewMerkleTree(leaves []string) *MerkleTree {
	// Keep the caller's order untouched. sort.Strings reorders its argument in
	// place, so sort a private copy, never the caller's slice. The stored
	// leavesUnsorted preserves index N -> leaf N for GetProof.
	leavesUnsorted := make([]string, len(leaves))
	copy(leavesUnsorted, leaves)

	sortedLeaves := make([]string, len(leaves))
	copy(sortedLeaves, leaves)
	sort.Strings(sortedLeaves)

	tree := make([][]byte, len(sortedLeaves)*2-1)
	for i, leaf := range sortedLeaves {
		tree[len(tree)-1-i] = hexutil.MustDecode(leaf)
	}

	for i := len(tree) - len(sortedLeaves) - 1; i >= 0; i-- {
		left := tree[leftChildIndex(i)]
		rightIndex := rightChildIndex(i)
		var right []byte

		// Check if the right child is out of bounds and skip if necessary
		if rightIndex >= len(tree) {
			right = []byte{}
		} else {
			right = tree[rightIndex]
		}

		tree[i] = Keccak256SortedHash(left, right)
	}

	finalTree := make([]string, len(tree))
	for i, node := range tree {
		nodeAsHex := fmt.Sprintf("0x%x", node)
		finalTree[i] = nodeAsHex
	}

	return &MerkleTree{
		tree:   finalTree,
		leaves: leavesUnsorted,
	}
}

func GetProof(leaves []string, index int) ([]string, error) {
	if index < 0 || index >= len(leaves) {
		return nil, errors.New("index out of bounds")
	}

	tree := NewMerkleTree(leaves)

	leafToProve := tree.leaves[index]

	// Leaves occupy the last len(tree.leaves) slots of tree.tree. Search only
	// that leaf region. A full scan also matches internal nodes, which is a
	// second route to the wrong node when a leaf value repeats an inner hash.
	leafRegionStart := len(tree.tree) - len(tree.leaves)
	var leafIndexInTree int
	var foundLeaf bool
	for i := leafRegionStart; i < len(tree.tree); i++ {
		if tree.tree[i] == leafToProve {
			foundLeaf = true
			leafIndexInTree = i
			break
		}
	}
	if !foundLeaf {
		return nil, errors.New("leaf not found in tree")
	}

	currentIndex := leafIndexInTree
	var proof []string

	// Traverse up the tree to build the proof.
	for currentIndex > 0 {
		siblingIndex, err := getSiblingIndex(currentIndex)
		if err != nil {
			return nil, err
		}

		// Add the sibling hash to the proof.
		if siblingIndex < len(tree.tree) {
			siblingHash := tree.tree[siblingIndex]
			proof = append(proof, siblingHash)
		}

		// Move to the parent index.
		currentIndex, err = parentIndex(currentIndex)
		if err != nil {
			return nil, err
		}
	}

	return proof, nil
}

func getSiblingIndex(i int) (int, error) {
	if i <= 0 {
		return 0, errors.New("root has no siblings")
	}
	return i - int(math.Pow(-1, float64(i%2))), nil
}

// parentIndex returns the parent index of a given index in the tree.
func parentIndex(i int) (int, error) {
	if i > 0 {
		return (i - 1) / 2, nil
	}
	return 0, errors.New("root has no parent")
}

// Deprecated: Use MerkleTree instead.
type MyMerkleTree = MerkleTree

// Deprecated: Use NewMerkleTree instead.
func MakeTree(leaves []string) *MerkleTree { return NewMerkleTree(leaves) }
