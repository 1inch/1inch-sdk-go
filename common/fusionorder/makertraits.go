package fusionorder

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

// MakerTraitsParams contains the parameters needed to create MakerTraits
// This is shared between fusion and fusionplus packages
type MakerTraitsParams struct {
	AuctionStartTime     uint32
	AuctionDuration      uint32
	OrderExpirationDelay uint32
	Nonce                *big.Int
	AllowPartialFills    bool
	AllowMultipleFills   bool
	UnwrapWeth           bool
	EnablePermit2        bool
}

// CreateMakerTraits creates MakerTraits from the provided parameters
// This is shared between fusion and fusionplus packages
func CreateMakerTraits(params MakerTraitsParams) (*orderbook.MakerTraits, error) {
	deadline := params.AuctionStartTime + params.AuctionDuration + params.OrderExpirationDelay
	var nonce int64
	if params.Nonce != nil {
		nonce = params.Nonce.Int64()
	}
	makerTraitParams := orderbook.MakerTraitsParams{
		Expiry:             int64(deadline),
		AllowPartialFills:  params.AllowPartialFills,
		AllowMultipleFills: params.AllowMultipleFills,
		HasPostInteraction: true,
		UnwrapWeth:         params.UnwrapWeth,
		UsePermit2:         params.EnablePermit2,
		HasExtension:       true,
		Nonce:              nonce,
	}
	makerTraits, err := orderbook.NewMakerTraits(makerTraitParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create maker traits: %w", err)
	}
	if makerTraits.IsBitInvalidatorMode() {
		if params.Nonce == nil || params.Nonce.Cmp(big.NewInt(0)) == 0 {
			return nil, errors.New("nonce required when partial fill or multiple fill disallowed")
		}
	}
	return makerTraits, nil
}
