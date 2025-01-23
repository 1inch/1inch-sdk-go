package fusion

import (
	"fmt"

	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Interaction struct {
	Target common.Address
	Data   string
}

func NewInteraction(target common.Address, data string) *Interaction {
	if _, err := hexutil.Decode(data); err != nil {
		panic("Interaction data must be valid hex bytes")
	}
	return &Interaction{
		Target: target,
		Data:   data,
	}
}

func (i *Interaction) Encode() string {
	return i.Target.String() + hexadecimal.Trim0x(i.Data)
}

func DecodeInteraction(bytes string) (*Interaction, error) {
	if !hexadecimal.IsHexBytes(bytes) {
		return nil, fmt.Errorf("invalid hex bytes: %s", bytes)
	}

	return &Interaction{
		Target: common.HexToAddress(bytes[:42]),
		Data:   fmt.Sprintf("0x%s", bytes[42:]),
	}, nil
}
