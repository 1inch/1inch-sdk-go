package fusionorder

import "math/big"

// ExtraData contains additional parameters for order creation
type ExtraData struct {
	UnwrapWETH           bool
	Nonce                *big.Int
	Permit               string
	AllowPartialFills    bool
	AllowMultipleFills   bool
	OrderExpirationDelay uint32
	EnablePermit2        bool
	Source               string
}
