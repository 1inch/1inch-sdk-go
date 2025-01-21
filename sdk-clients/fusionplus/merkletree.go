package fusionplus

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type MyMerkleTree struct {
	tree   []string
	leaves []string
}

func MakeTree(leaves []string) *MyMerkleTree {
	leavesUnsorted := make([]string, len(leaves))
	copy(leavesUnsorted, leaves)
	sort.Strings(leaves)

	tree := make([][]byte, len(leaves)*2-1)
	for i, leaf := range leaves {
		tree[len(tree)-1-i] = hexutil.MustDecode(leaf)
	}

	for i := len(tree) - len(leaves) - 1; i >= 0; i-- {
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

	return &MyMerkleTree{
		tree:   finalTree,
		leaves: leavesUnsorted,
	}
}

func GetProof(leaves []string, index int) ([]string, error) {
	if index < 0 || index >= len(leaves) {
		return nil, errors.New("index out of bounds")
	}

	tree := MakeTree(leaves)

	leafToProve := tree.leaves[index]

	var leafIndexInTree int
	var foundLeaf bool
	for i, leaf := range tree.tree {
		if leaf == leafToProve {
			foundLeaf = true
			leafIndexInTree = i
			break
		}
	}
	if !foundLeaf {
		panic("Leaf not found in tree")
	}

	currentIndex := leafIndexInTree
	var proof []string

	// Traverse up the tree to build the proof.
	for currentIndex > 0 {
		siblingIndex := getSiblingIndex(currentIndex)

		// Add the sibling hash to the proof.
		if siblingIndex < len(tree.tree) {
			siblingHash := tree.tree[siblingIndex]
			proof = append(proof, siblingHash)
		}

		// Move to the parent index.
		var err error
		currentIndex, err = ParentIndex(currentIndex)
		if err != nil {
			return nil, err
		}
	}

	return proof, nil
}

func getSiblingIndex(i int) int {
	if i <= 0 {
		panic("Root has no siblings")
	}
	return i - int(math.Pow(-1, float64(i%2)))
}

// ParentIndex returns the parent index of a given index in the tree.
func ParentIndex(i int) (int, error) {
	if i > 0 {
		return (i - 1) / 2, nil
	}
	return 0, errors.New("root has no parent")
}
