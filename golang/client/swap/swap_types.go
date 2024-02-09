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
	validationErrors = validate.Parameter(int(params.ApprovalType), "approvalType", validate.CheckApprovalType, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.PublicAddress, "publicAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.WalletKey, "walletKey", validate.CheckPrivateKeyRequired, validationErrors)
	validationErrors = validate.Parameter(params.Src, "src", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Dst, "dst", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(params.From, "from", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Slippage, "slippage", validate.CheckSlippageRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type ApproveAllowanceParams struct {
	ChainId int
	ApproveControllerGetAllowanceParams
}

func (params *ApproveAllowanceParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.TokenAddress, "tokenAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.WalletAddress, "walletAddress", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type ApproveSpenderParams struct {
	ChainId int
}

func (params *ApproveSpenderParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type ApproveTransactionParams struct {
	ChainId int
	ApproveControllerGetCallDataParams
}

func (params *ApproveTransactionParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.TokenAddress, "tokenAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.CheckBigInt, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetLiquiditySourcesParams struct {
	ChainId int
}

func (params *GetLiquiditySourcesParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetQuoteParams struct {
	ChainId int
	AggregationControllerGetQuoteParams
}

func (params *GetQuoteParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.Src, "src", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Dst, "dst", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.CheckBigIntRequired, validationErrors)
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
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.Src, "src", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Dst, "dst", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(params.From, "from", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Slippage, "slippage", validate.CheckSlippageRequired, validationErrors)
	if params.Src == params.Dst && len(params.Src) > 0 && len(params.Dst) > 0 {
		validationErrors = append(validationErrors, errors.New("src and dst tokens must be different"))
	}
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetTokensParams struct {
	ChainId int
}

func (params *GetTokensParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
