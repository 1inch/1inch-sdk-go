package swap

import (
	clienterrors "github.com/1inch/1inch-sdk/golang/client/errors"
)

func (params *AggregationControllerGetQuoteParams) Validate() error {
	if params.Src == "" {
		return clienterrors.NewRequestValidationError("src is required")
	}
	if params.Dst == "" {
		return clienterrors.NewRequestValidationError("dst is required")
	}
	if params.Amount == "" {
		return clienterrors.NewRequestValidationError("amount is required")
	}
	if params.Src == params.Dst {
		return clienterrors.NewRequestValidationError("src and dst tokens must be different")
	}
	return nil
}

func (params *AggregationControllerGetSwapParams) Validate() error {
	if params.Src == "" {
		return clienterrors.NewRequestValidationError("src is required")
	}
	if params.Dst == "" {
		return clienterrors.NewRequestValidationError("dst is required")
	}
	if params.Amount == "" {
		return clienterrors.NewRequestValidationError("amount is required")
	}
	if params.From == "" {
		return clienterrors.NewRequestValidationError("from is required")
	}
	if params.Src == params.Dst {
		return clienterrors.NewRequestValidationError("src and dst tokens must be different")
	}
	return nil
}

func (params *ApproveControllerGetCallDataParams) Validate() error {
	if params.TokenAddress == "" {
		return clienterrors.NewRequestValidationError("tokenAddress is required")
	}
	return nil
}

func (params *ApproveControllerGetAllowanceParams) Validate() error {
	if params.TokenAddress == "" {
		return clienterrors.NewRequestValidationError("tokenAddress is required")
	}
	if params.WalletAddress == "" {
		return clienterrors.NewRequestValidationError("walletAddress is required")
	}
	return nil
}
