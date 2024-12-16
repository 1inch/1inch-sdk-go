package fusionplus

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Interaction struct {
	Target common.Address
	Data   string
}

func NewInteraction(target common.Address, data string) *Interaction {
	if !isHexBytes(data) {
		panic("Interaction data must be valid hex bytes")
	}
	return &Interaction{
		Target: target,
		Data:   data,
	}
}

func (i *Interaction) Encode() string {
	return i.Target.String() + trim0x(i.Data)
}

func DecodeInteraction(bytes string) (*Interaction, error) {
	if !isHexBytes(bytes) {
		return nil, fmt.Errorf("invalid hex bytes: %s", bytes)
	}

	return &Interaction{
		Target: common.HexToAddress(bytes[:42]),
		Data:   fmt.Sprintf("0x%s", bytes[42:]),
	}, nil
}

func isHexBytes(s string) bool {
	_, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	return err == nil
}
