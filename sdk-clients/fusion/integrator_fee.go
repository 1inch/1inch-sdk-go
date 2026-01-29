package fusion

import (
	"errors"
	"fmt"

	"github.com/1inch/1inch-sdk-go/internal/addresses"
	"github.com/1inch/1inch-sdk-go/sdk-clients/fusionorder"
)

type IntegratorFee struct {
	Integrator string
	Protocol   string
	Fee        *Bps
	Share      *Bps
}

// IntegratorFeeZero is a safe default with all zero values
var IntegratorFeeZero = &IntegratorFee{
	Integrator: addresses.ZeroAddress,
	Protocol:   addresses.ZeroAddress,
	Fee:        fusionorder.BpsZero,
	Share:      fusionorder.BpsZero,
}

// NewIntegratorFee constructs a validated IntegratorFee or returns an error
func NewIntegratorFee(integrator, protocol string, fee, share *Bps) (*IntegratorFee, error) {
	if fee.IsZero() {
		if !share.IsZero() {
			return nil, errors.New("zero fee requires zero integrator share")
		}
		if integrator != addresses.ZeroAddress {
			return nil, errors.New("zero fee requires zero integrator address")
		}
		if protocol != addresses.ZeroAddress {
			return nil, errors.New("zero fee requires zero protocol address")
		}
	}

	if (integrator == addresses.ZeroAddress || protocol == addresses.ZeroAddress) && !fee.IsZero() {
		return nil, errors.New("non-zero fee requires non-zero integrator and protocol addresses")
	}

	return &IntegratorFee{
		Integrator: integrator,
		Protocol:   protocol,
		Fee:        fee,
		Share:      share,
	}, nil
}

func (f *IntegratorFee) String() string {
	return fmt.Sprintf("IntegratorFee{Integrator: %s, Protocol: %s, Fee: %s, Share: %s}",
		f.Integrator, f.Protocol, f.Fee.String(), f.Share.String())
}
