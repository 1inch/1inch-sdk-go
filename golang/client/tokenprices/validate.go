package tokenprices

import (
	"fmt"

	clienterrors "github.com/1inch/1inch-sdk/golang/client/errors"
)

func (params *ChainControllerByAddressesParams) Validate() error {
	if params.Currency != nil && *params.Currency != "" {
		if !contains(currencies, *params.Currency) {
			return clienterrors.NewRequestValidationError(fmt.Sprintf("currency value %s is not valid", string(*params.Currency)))
		}
	}
	return nil
}

func contains(slice []ChainControllerByAddressesParamsCurrency, item ChainControllerByAddressesParamsCurrency) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

var currencies = []ChainControllerByAddressesParamsCurrency{
	AED, ARS, AUD, BDT, BHD, BMD, BRL, CAD, CHF, CLP,
	CNY, CZK, DKK, EUR, GBP, HKD, HUF, IDR, ILS, INR,
	JPY, KRW, KWD, LKR, MMK, MXN, MYR, NGN, NOK, NZD,
	PHP, PKR, PLN, RUB, SAR, SEK, SGD, THB, TRY, TWD,
	UAH, USD, VEF, VND, ZAR,
}
