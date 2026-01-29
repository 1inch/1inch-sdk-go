package fusionorder

import (
	"math/big"
)

// ExtraParams contains additional parameters for order creation
type ExtraParams struct {
	Nonce                *big.Int
	Permit               string
	AllowPartialFills    bool
	AllowMultipleFills   bool
	OrderExpirationDelay uint32
	EnablePermit2        bool
	Source               string
	UnwrapWeth           bool
}

// DetailsBase contains the common fields for order details
type DetailsBase struct {
	Auction            *AuctionDetails
	Whitelist          []AuctionWhitelistItem
	ResolvingStartTime *big.Int
}

// IsNonceRequired returns true if a nonce is required based on fill settings
func IsNonceRequired(allowPartialFills, allowMultipleFills bool) bool {
	return !allowPartialFills || !allowMultipleFills
}
