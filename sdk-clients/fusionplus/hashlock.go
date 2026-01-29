package fusionplus

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
	"github.com/1inch/1inch-sdk-go/internal/keccak"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

type HashLock struct {
	Value string
}

func ForSingleFill(secret string) (*HashLock, error) {
	hashlock, err := HashSecret(secret)
	if err != nil {
		return nil, err
	}
	return &HashLock{
		hashlock,
	}, nil
}

func ForMultipleFills(leaves []string) (*HashLock, error) {
	// Assertion to check the number of leaves
	if len(leaves) <= 2 {
		return nil, errors.New("leaves array requires more than 2 elements")
	}

	tree := MakeTree(leaves)

	root := tree.tree[0]

	// Convert root from []byte to big.Int
	rootAsBytes := hexutil.MustDecode(root)
	rootBig := new(big.Int).SetBytes(rootAsBytes)

	// Specify the mask starting bit and length (for example, bits 240 to 256 for a 256-bit value)
	maskStart := 240
	maskLength := 16

	countValue := big.NewInt(int64(len(leaves) - 1))

	// Apply the mask using the SetMask function
	rootWithCount := SetMask(rootBig, uint(maskStart), uint(maskLength), countValue)

	// Convert the modified root back to []byte
	rootWithCountBytes := rootWithCount.Bytes()

	// Ensure the byte slice has the correct length (32 bytes)
	if len(rootWithCountBytes) < 32 {
		// Pad the byte slice to 32 bytes if needed
		padding := make([]byte, 32-len(rootWithCountBytes))
		rootWithCountBytes = append(padding, rootWithCountBytes...)
	}

	// Create and return the HashLock
	return &HashLock{fmt.Sprintf("0x%x", rootWithCountBytes)}, nil
}

func GetMerkleLeaves(secrets []string) ([]string, error) {
	secretHashes := make([]string, len(secrets))
	for i, secret := range secrets {
		hash, err := HashSecret(secret)
		if err != nil {
			return nil, err
		}
		secretHashes[i] = hash
	}
	return GetMerkleLeavesFromSecretHashes(secretHashes)
}

func GetMerkleLeavesFromSecretHashes(secretHashes []string) ([]string, error) {
	var leaves []string
	for idx, s := range secretHashes {
		hash, err := solidityPackedKeccak256([]string{"uint64", "bytes32"}, []interface{}{idx, s})
		if err != nil {
			return nil, err
		}
		leaves = append(leaves, hash)
	}
	return leaves, nil
}

func GetRandomBytes32() (string, error) {
	// Create a byte slice with a length of 32
	bytes := make([]byte, 32)

	// Read random bytes into the slice
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Convert bytes to a hexadecimal string and prepend "0x"
	return "0x" + hex.EncodeToString(bytes), nil
}

func SetMask(original *big.Int, maskStart, maskLength uint, value *big.Int) *big.Int {
	// Create a mask that only affects the bits we want to modify
	mask := new(big.Int).Lsh(big.NewInt(1), maskLength)
	mask.Sub(mask, big.NewInt(1))
	mask.Lsh(mask, maskStart)

	// Clear the bits of the original value that are set in the mask
	originalCleared := new(big.Int).AndNot(original, mask)

	// Shift the value into the correct position
	maskedValue := new(big.Int).Lsh(value, maskStart)

	// Combine cleared original with the masked value
	result := new(big.Int).Or(originalCleared, maskedValue)

	return result
}

// Keccak256SortedHash takes two byte slices, sorts them lexicographically, concatenates them, and hashes the result.
func Keccak256SortedHash(a, b []byte) []byte {
	// Sort the inputs lexicographically
	if bytes.Compare(a, b) > 0 {
		a, b = b, a
	}

	// Concatenate the sorted byte slices
	concatenated := append(a, b...)

	// Apply the Keccak-256 hash on the concatenated result
	return crypto.Keccak256(concatenated)
}

func leftChildIndex(i int) int {
	return 2*i + 1
}
func rightChildIndex(i int) int {
	return 2*i + 2
}

func getBytesCount(hex string) int {
	return len(hexadecimal.Trim0x(hex)) / 2
}

func HashSecret(secret string) (string, error) {
	if !hexadecimal.IsHexBytes(secret) || getBytesCount(secret) != 32 {
		return "", fmt.Errorf("invalid secret length: expected 32 bytes hex encoded, got %s", secret)
	}

	hexBytes, err := hexutil.Decode(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %w", err)
	}

	return keccak.Keccak256Legacy(hexBytes), nil
}

func hexlify(data []byte) string {
	return "0x" + hex.EncodeToString(data)
}

func concat(datas [][]byte) []byte {
	result := []byte{}
	for _, d := range datas {
		result = append(result, d...)
	}
	return result
}

func solidityPacked(types []string, values []interface{}) ([]byte, error) {
	if len(types) != len(values) {
		return nil, fmt.Errorf("value count mismatch: expected %d", len(types))
	}

	var tight [][]byte
	for i, t := range types {
		packed, err := pack(t, values[i])
		if err != nil {
			return nil, err
		}
		tight = append(tight, packed)
	}

	return concat(tight), nil
}

func solidityPackedKeccak256(types []string, values []interface{}) (string, error) {
	packed, err := solidityPacked(types, values)
	if err != nil {
		return "", err
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write(packed)
	hashed := hash.Sum(nil)

	return hexlify(hashed), nil
}

func pack(typ string, value interface{}) ([]byte, error) {
	switch typ {
	case "uint64":
		// Pack uint64 as big-endian 8-byte array
		v, ok := value.(int)
		if !ok {
			return nil, fmt.Errorf("invalid uint64 type: expected int, got %T", value)
		}
		bigInt := big.NewInt(int64(v))
		packed := bigInt.FillBytes(make([]byte, 8))
		return packed, nil

	case "bytes32":
		// Pack bytes32 as exactly 32 bytes
		s, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid bytes32 type: expected string, got %T", value)
		}
		bytes, err := hexutil.Decode(s)
		if err != nil {
			return nil, err
		}
		if len(bytes) != 32 {
			return nil, fmt.Errorf("invalid bytes32 length: expected 32 bytes, got %d", len(bytes))
		}
		return bytes, nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", typ)
	}
}
