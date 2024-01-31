package swap

import "errors"

func (params *AggregationControllerGetSwapParams) Validate() error {
	if params.Src == "" {
		return errors.New("src is required")
	}
	if params.Dst == "" {
		return errors.New("dst is required")
	}
	if params.Amount == "" {
		return errors.New("amount is required")
	}
	if params.From == "" {
		return errors.New("from is required")
	}
	if params.Src == params.Dst {
		return errors.New("src and dst tokens must be different")
	}
	if params.Slippage == 0 {
		return errors.New("slippage is required")
	}
	return nil
}

func (params *ApproveControllerGetCallDataParams) Validate() error {
	if params.TokenAddress == "" {
		return errors.New("tokenAddress is required")
	}
	return nil
}

func (params *ApproveControllerGetAllowanceParams) Validate() error {
	if params.TokenAddress == "" {
		return errors.New("tokenAddress is required")
	}
	if params.WalletAddress == "" {
		return errors.New("walletAddress is required")
	}
	return nil
}
