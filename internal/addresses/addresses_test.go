package addresses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroAddress(t *testing.T) {
	// Verify the zero address constant is correctly defined
	assert.Equal(t, "0x0000000000000000000000000000000000000000", ZeroAddress)

	// Verify it has the correct length (42 characters: 0x + 40 hex chars)
	assert.Len(t, ZeroAddress, 42)

	// Verify it starts with 0x
	assert.Equal(t, "0x", ZeroAddress[:2])

	// Verify all characters after 0x are zeros
	for _, c := range ZeroAddress[2:] {
		assert.Equal(t, '0', c)
	}
}
