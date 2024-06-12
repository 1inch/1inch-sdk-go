package fusion

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestInteraction(t *testing.T) {
	tests := []struct {
		name   string
		target common.Address
		data   string
	}{
		{
			name:   "Encode/Decode Interaction",
			target: common.HexToAddress("0x0000000000000000000000000000000000000539"), // 1337 in hexadecimal
			data:   "0xdeadbeef",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			interaction := NewInteraction(tc.target, tc.data)
			encoded := interaction.Encode()
			decoded := DecodeInteraction(encoded)

			assert.Equal(t, interaction.Target, decoded.Target)
			assert.Equal(t, interaction.Data, decoded.Data)
		})
	}
}
