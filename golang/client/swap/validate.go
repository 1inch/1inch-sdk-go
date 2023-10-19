package swap

import "fmt"

func (params *AggregationControllerGetQuoteParams) Validate() error {
	if params.Src == "" {
		return fmt.Errorf("src is required")
	}
	if params.Dst == "" {
		return fmt.Errorf("dst is required")
	}
	if params.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	if params.Src == params.Dst {
		return fmt.Errorf("src and dst tokens must be different")
	}
	return nil
}

func (params *AggregationControllerGetSwapParams) Validate() error {
	if params.Src == "" {
		return fmt.Errorf("src is required")
	}
	if params.Dst == "" {
		return fmt.Errorf("dst is required")
	}
	if params.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	if params.From == "" {
		return fmt.Errorf("from is required")
	}
	if params.Src == params.Dst {
		return fmt.Errorf("src and dst tokens must be different")
	}
	return nil
}

func (params *ApproveControllerGetCallDataParams) Validate() error {
	if params.TokenAddress == "" {
		return fmt.Errorf("tokenAddress is required")
	}
	return nil
}

func (params *ApproveControllerGetAllowanceParams) Validate() error {
	if params.TokenAddress == "" {
		return fmt.Errorf("tokenAddress is required")
	}
	if params.WalletAddress == "" {
		return fmt.Errorf("walletAddress is required")
	}
	return nil
}
