package swap

import (
	"errors"

	"github.com/1inch/1inch-sdk/golang/client/validate"
)

type SwapTokensParams struct {
	ApprovalType  ApprovalType
	ChainId       int
	SkipWarnings  bool
	PublicAddress string
	WalletKey     string
	AggregationControllerGetSwapParams
}

func (params *SwapTokensParams) Validate() error {
	var validationErrors []error
	if err := validate.ApprovalType(int(params.ApprovalType), "approvalType"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.PublicAddress, "publicAddress"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.PrivateKey(params.WalletKey, "walletKey"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.Src, "src"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.Dst, "dst"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.BigInt(params.Amount, "amount"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.From, "from"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.Slippage(params.Slippage, "slippage"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
}

type ApproveAllowanceParams struct {
	ChainId int
	ApproveControllerGetAllowanceParams
}

func (params *ApproveAllowanceParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.TokenAddress, "tokenAddress"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.WalletAddress, "walletAddress"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
}

type ApproveSpenderParams struct {
	ChainId int
}

func (params *ApproveSpenderParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
}

type ApproveTransactionParams struct {
	ChainId int
	ApproveControllerGetCallDataParams
}

func (params *ApproveTransactionParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.TokenAddress, "tokenAddress"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.BigIntPointer(params.Amount, "amount"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
}

type GetLiquiditySourcesParams struct {
	ChainId int
}

func (params *GetLiquiditySourcesParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
}

type GetQuoteParams struct {
	ChainId int
	AggregationControllerGetQuoteParams
}

func (params *GetQuoteParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.Src, "src"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.Dst, "dst"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.BigInt(params.Amount, "amount"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if params.Src == params.Dst {
		validationErrors = append(validationErrors, errors.New("src and dst tokens must be different"))
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
}

type GetSwapDataParams struct {
	ChainId      int
	SkipWarnings bool
	AggregationControllerGetSwapParams
}

func (params *GetSwapDataParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.Src, "src"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.Dst, "dst"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.BigInt(params.Amount, "amount"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.From, "from"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.Slippage(params.Slippage, "slippage"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if params.Src == params.Dst {
		validationErrors = append(validationErrors, errors.New("src and dst tokens must be different"))
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
}

type GetTokensParams struct {
	ChainId int
}

func (params *GetTokensParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
}
