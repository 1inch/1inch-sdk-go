package fusionorder

import (
	"fmt"
	"strings"

	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Interaction represents an interaction with a target contract
type Interaction struct {
	Target common.Address
	Data   string
}

// NewInteraction creates a new Interaction with validated hex data
func NewInteraction(target common.Address, data string) (*Interaction, error) {
	if _, err := hexutil.Decode(data); err != nil {
		return nil, fmt.Errorf("failed to decode interaction data: %w", err)
	}
	return &Interaction{
		Target: target,
		Data:   data,
	}, nil
}

// Encode encodes the interaction as a hex string (lowercase target + data)
func (i *Interaction) Encode() string {
	return strings.ToLower(i.Target.String()) + hexadecimal.Trim0x(i.Data)
}

// DecodeInteraction decodes a hex string into an Interaction
func DecodeInteraction(bytes string) (*Interaction, error) {
	if !hexadecimal.IsHexBytes(bytes) {
		return nil, fmt.Errorf("invalid interaction hex: %s", bytes)
	}

	if len(bytes) < 42 {
		return nil, fmt.Errorf("interaction data too short: requires at least 20 bytes")
	}

	return &Interaction{
		Target: common.HexToAddress(bytes[:42]),
		Data:   fmt.Sprintf("0x%s", bytes[42:]),
	}, nil
}
