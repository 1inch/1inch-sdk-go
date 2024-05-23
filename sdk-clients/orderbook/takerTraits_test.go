package orderbook

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func TestTakerTraitsEncode(t *testing.T) {

	tests := []struct {
		name                string
		takerTraitParams    TakerTraitsParams
		expectedTakerTraits string
		expectedTakerArgs   string
	}{
		{
			name: "Extension",
			takerTraitParams: TakerTraitsParams{
				Extension: "0x000000f4000000f4000000f4000000000000000000000000000000000000000045c32fa6df82ead1e2ef74d17b76547eddfaff8900000000000000000000000050c5df26654b5efbdd0c54a062dfa6012933defe000000000000000000000000111111125421ca6dc452d289314280a0f8842a65000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000663a478b000000000000000000000000000000000000000000000000000000000000001bdf138a0d223e2ef8635075f5fe68efa8a2da1d890fdc3825b754c7ba2083ca0464494f534829f576cd67b966059657c51aaf53edbd6498d51cbd07da8bdb256b",
			},
			expectedTakerTraits: "7440945280133576583328096164017418065923851860621198004784596428783616",
			expectedTakerArgs:   "0x000000f4000000f4000000f4000000000000000000000000000000000000000045c32fa6df82ead1e2ef74d17b76547eddfaff8900000000000000000000000050c5df26654b5efbdd0c54a062dfa6012933defe000000000000000000000000111111125421ca6dc452d289314280a0f8842a65000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000663a478b000000000000000000000000000000000000000000000000000000000000001bdf138a0d223e2ef8635075f5fe68efa8a2da1d890fdc3825b754c7ba2083ca0464494f534829f576cd67b966059657c51aaf53edbd6498d51cbd07da8bdb256b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			expectedTakerTraitsBig, err := validate.BigIntFromString(tc.expectedTakerTraits)
			require.NoError(t, err)

			takerTraits := NewTakerTraits(tc.takerTraitParams)
			takerTraitsEnocded := takerTraits.Encode()
			assert.True(t, expectedTakerTraitsBig.Cmp(takerTraitsEnocded.TraitFlags) == 0, fmt.Sprintf("Expected %x, got %x", expectedTakerTraitsBig, takerTraitsEnocded.TraitFlags))

			expectedTakerArgs := common.FromHex(tc.expectedTakerArgs)
			require.NoError(t, err)
			assert.Equal(t, expectedTakerArgs, takerTraitsEnocded.Args)
		})
	}
}
