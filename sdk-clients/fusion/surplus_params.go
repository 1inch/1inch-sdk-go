package fusion

import (
	"errors"
	"math/big"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/constants"
)

type SurplusParams struct {
	EstimatedTakerAmount *big.Int
	ProtocolFee          *Bps
}

// SurplusParamsNoFee is equivalent to SurplusParams.NO_FEE in TS
var SurplusParamsNoFee, _ = NewSurplusParams(constants.Uint256Max, fusionorder.BpsZero)

// NewSurplusParams validates that the protocolFee is in whole percent increments
func NewSurplusParams(estimatedTakerAmount *big.Int, protocolFee *Bps) (*SurplusParams, error) {
	if new(big.Int).Rem(protocolFee.Value(), big.NewInt(100)).Sign() != 0 {
		return nil, errors.New("protocol fee must be an integer percent")
	}
	return &SurplusParams{
		EstimatedTakerAmount: new(big.Int).Set(estimatedTakerAmount),
		ProtocolFee:          protocolFee,
	}, nil
}
