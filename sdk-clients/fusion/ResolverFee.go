package fusion

import (
	"errors"
	"fmt"
	"math/big"
)

// ResolverFee represents a fee paid to a receiver, with optional whitelist discount
type ResolverFee struct {
	Receiver          string
	Fee               *Bps
	WhitelistDiscount *Bps
}

// ResolverFeeZero is the default zero instance
var ResolverFeeZero = &ResolverFee{
	Receiver:          ZeroAddress,
	Fee:               BpsZero,
	WhitelistDiscount: BpsZero,
}

// NewResolverFee constructs a validated ResolverFee instance
func NewResolverFee(receiver string, fee *Bps, whitelistDiscount *Bps) (*ResolverFee, error) {
	if (receiver == "" || receiver == ZERO_ADDRESS) && !fee.IsZero() {
		return nil, errors.New("fee must be zero if receiver is zero address")
	}
	if !(receiver == "" || receiver == ZERO_ADDRESS) && fee.IsZero() {
		return nil, errors.New("receiver must be zero address if fee is zero")
	}
	if fee.IsZero() && !whitelistDiscount.IsZero() {
		return nil, errors.New("whitelist discount must be zero if fee is zero")
	}

	// Check percent precision: must be divisible by 100 (i.e. 1%)
	if new(big.Int).Rem(whitelistDiscount.value, big.NewInt(100)).Sign() != 0 {
		return nil, errors.New("whitelist discount must have percent precision: 1%, 2% and so on")
	}

	return &ResolverFee{
		Receiver:          receiver,
		Fee:               fee,
		WhitelistDiscount: whitelistDiscount,
	}, nil
}

func (r *ResolverFee) String() string {
	return fmt.Sprintf("ResolverFee{Receiver: %s, Fee: %s, WhitelistDiscount: %s}",
		r.Receiver, r.Fee.String(), r.WhitelistDiscount.String())
}
