package fusion

import (
	"math/big"
	"testing"

	"github.com/1inch/1inch-sdk-go/sdk-clients/fusionorder"
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
			interaction, err := fusionorder.NewInteraction(tc.target, tc.data)
			require.NoError(t, err)
			encoded := interaction.Encode()
			decoded, err := fusionorder.DecodeInteraction(encoded)
			require.NoError(t, err)
			assert.Equal(t, interaction.Target, decoded.Target)
			assert.Equal(t, interaction.Data, decoded.Data)
		})
	}
}
