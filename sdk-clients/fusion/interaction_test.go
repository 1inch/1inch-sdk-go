package fusion

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInteraction(t *testing.T) {
	tests := []struct {
		name   string
		target common.Address
		data   string
	}{
		{
			name:   "Encode/Decode Interaction",
			target: common.BigToAddress(big.NewInt(1337)),
			data:   "0xdeadbeef",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			interaction := NewInteraction(tc.target, tc.data)
			fmt.Printf("interaction: %v\n", interaction)
			encoded := interaction.Encode()
			fmt.Printf("encoded: %v\n", encoded)
			decoded, err := DecodeInteraction(encoded)
			require.NoError(t, err)
			fmt.Printf("decoded: %v\n", decoded)

			assert.Equal(t, interaction.Target, decoded.Target)
			assert.Equal(t, interaction.Data, decoded.Data)
		})
	}
}
