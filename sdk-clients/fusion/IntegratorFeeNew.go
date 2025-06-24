package fusion

import (
	"errors"
	"fmt"
)

const ZeroAddress string = "0x0000000000000000000000000000000000000000"

type IntegratorFeeNew struct {
	Integrator string
	Protocol   string
	Fee        *Bps
	Share      *Bps
}

// IntegratorFeeZero is a safe default with all zero values
var IntegratorFeeZero = &IntegratorFeeNew{
	Integrator: ZeroAddress,
	Protocol:   ZeroAddress,
	Fee:        BpsZero,
	Share:      BpsZero,
}

// NewIntegratorFee constructs a validated IntegratorFeeNew or returns an error
func NewIntegratorFee(integrator, protocol string, fee, share *Bps) (*IntegratorFeeNew, error) {
	if fee.IsZero() {
		if !share.IsZero() {
			return nil, errors.New("integrator share must be zero if fee is zero")
		}
		if integrator != ZeroAddress {
			return nil, errors.New("integrator address must be zero if fee is zero")
		}
		if protocol != ZeroAddress {
			return nil, errors.New("protocol address must be zero if fee is zero")
		}
	}

	if (integrator == ZeroAddress || protocol == ZeroAddress) && !fee.IsZero() {
		return nil, errors.New("fee must be zero if integrator or protocol is zero address")
	}

	return &IntegratorFeeNew{
		Integrator: integrator,
		Protocol:   protocol,
		Fee:        fee,
		Share:      share,
	}, nil
}

func (f *IntegratorFeeNew) String() string {
	return fmt.Sprintf("IntegratorFeeNew{Integrator: %s, Protocol: %s, Fee: %s, Share: %s}",
		f.Integrator, f.Protocol, f.Fee.String(), f.Share.String())
}
