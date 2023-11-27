package fusion

import (
	"errors"

	clienterrors "1inch-sdk-golang/client/errors"
	"1inch-sdk-golang/helpers"
)

func (params *OrderApiControllerGetActiveOrdersParams) Validate() error {
	if params.Page != nil {
		if *params.Page < 1 {
			return clienterrors.NewRequestValidationError("page must be greater than 0")
		}
	}
	if params.Limit != nil {
		if *params.Limit < 1 {
			return clienterrors.NewRequestValidationError("limit must be greater than 0")
		}
	}
	return nil
}

func (params *QuoterControllerGetQuoteParams) Validate() error {
	// Validate Ethereum addresses
	if !helpers.IsEthereumAddress(params.FromTokenAddress) {
		return errors.New("fromTokenAddress must be a valid Ethereum address")
	}
	if !helpers.IsEthereumAddress(params.ToTokenAddress) {
		return errors.New("toTokenAddress must be a valid Ethereum address")
	}
	if !helpers.IsEthereumAddress(params.WalletAddress) {
		return errors.New("walletAddress must be a valid Ethereum address")
	}

	// Validate amount
	if params.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if params.Fee != nil {
		if *params.Fee < 0 {
			return errors.New("fee must not be negative")
		}
	}

	return nil
}
