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
	validationErrors = validate.Parameter(int(params.ApprovalType), "approvalType", validate.ApprovalType, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.ChainID, validationErrors)
	validationErrors = validate.Parameter(params.PublicAddress, "publicAddress", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.WalletKey, "walletKey", validate.PrivateKey, validationErrors)
	validationErrors = validate.Parameter(params.Src, "src", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Dst, "dst", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.BigInt, validationErrors)
	validationErrors = validate.Parameter(params.From, "from", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Slippage, "slippage", validate.Slippage, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type ApproveAllowanceParams struct {
	ChainId int
	ApproveControllerGetAllowanceParams
}

func (params *ApproveAllowanceParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.ChainID, validationErrors)
	validationErrors = validate.Parameter(params.TokenAddress, "tokenAddress", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.WalletAddress, "walletAddress", validate.EthereumAddress, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type ApproveSpenderParams struct {
	ChainId int
}

func (params *ApproveSpenderParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.ChainID, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type ApproveTransactionParams struct {
	ChainId int
	ApproveControllerGetCallDataParams
}

func (params *ApproveTransactionParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.ChainID, validationErrors)
	validationErrors = validate.Parameter(params.TokenAddress, "tokenAddress", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.BigIntPointer, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetLiquiditySourcesParams struct {
	ChainId int
}

func (params *GetLiquiditySourcesParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.ChainID, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetQuoteParams struct {
	ChainId int
	AggregationControllerGetQuoteParams
}

func (params *GetQuoteParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.ChainID, validationErrors)
	validationErrors = validate.Parameter(params.Src, "src", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Dst, "dst", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.BigInt, validationErrors)
	if len(params.Src) > 0 && len(params.Dst) > 0 && params.Src == params.Dst {
		validationErrors = append(validationErrors, errors.New("src and dst tokens must be different"))
	}
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetSwapDataParams struct {
	ChainId      int
	SkipWarnings bool
	AggregationControllerGetSwapParams
}

func (params *GetSwapDataParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.ChainID, validationErrors)
	validationErrors = validate.Parameter(params.Src, "src", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Dst, "dst", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.BigInt, validationErrors)
	validationErrors = validate.Parameter(params.From, "from", validate.EthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Slippage, "slippage", validate.Slippage, validationErrors)
	if len(params.Src) > 0 && len(params.Dst) > 0 && params.Src == params.Dst {
		validationErrors = append(validationErrors, errors.New("src and dst tokens must be different"))
	}
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetTokensParams struct {
	ChainId int
}

func (params *GetTokensParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.ChainID, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
