package aggregation

import (
	"errors"

	"github.com/1inch/1inch-sdk-go/internal/onchain"
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

type SwapTokensParams struct {
	ApprovalType onchain.ApprovalType
	ChainId      int
	Address      string
	WalletKey    string
	AggregationControllerGetSwapParams
}

// TODO Add validation to all optional parameters here

func (params *SwapTokensParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(int(params.ApprovalType), "approvalType", validate.CheckApprovalType, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.Address, "publicAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.WalletKey, "walletKey", validate.CheckPrivateKeyRequired, validationErrors)
	validationErrors = validate.Parameter(params.Src, "src", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Dst, "dst", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(params.From, "from", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Slippage, "slippage", validate.CheckSlippageRequired, validationErrors)
	validationErrors = validate.Parameter(params.Protocols, "protocols", validate.CheckProtocols, validationErrors)
	validationErrors = validate.Parameter(params.Fee, "fee", validate.CheckFee, validationErrors)
	validationErrors = validate.Parameter(params.GasPrice, "gasPrice", validate.CheckBigInt, validationErrors)
	validationErrors = validate.Parameter(params.ComplexityLevel, "complexityLevel", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.Parts, "parts", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.MainRouteParts, "mainRouteParts", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.GasLimit, "gasLimit", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.ConnectorTokens, "connectorTokens", validate.CheckConnectorTokens, validationErrors)
	validationErrors = validate.Parameter(params.Permit, "permit", validate.CheckPermitHash, validationErrors)
	validationErrors = validate.Parameter(params.Receiver, "receiver", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Referrer, "referrer", validate.CheckEthereumAddress, validationErrors)
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
	validationErrors = validate.Parameter(params.Protocols, "protocols", validate.CheckProtocols, validationErrors)
	validationErrors = validate.Parameter(params.Fee, "fee", validate.CheckFee, validationErrors)
	validationErrors = validate.Parameter(params.GasPrice, "gasPrice", validate.CheckBigInt, validationErrors)
	validationErrors = validate.Parameter(params.ComplexityLevel, "complexityLevel", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.Parts, "parts", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.MainRouteParts, "mainRouteParts", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.GasLimit, "gasLimit", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.ConnectorTokens, "connectorTokens", validate.CheckConnectorTokens, validationErrors)
	if len(params.Src) > 0 && len(params.Dst) > 0 && params.Src == params.Dst {
		validationErrors = append(validationErrors, errors.New("src and dst tokens must be different"))
	}
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetSwapParams struct {
	ChainId      int
	SkipWarnings bool
	AggregationControllerGetSwapParams
}

func (params *GetSwapParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.Src, "src", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Dst, "dst", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "amount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(params.From, "from", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Slippage, "slippage", validate.CheckSlippageRequired, validationErrors)
	validationErrors = validate.Parameter(params.Protocols, "protocols", validate.CheckProtocols, validationErrors)
	validationErrors = validate.Parameter(params.Fee, "fee", validate.CheckFee, validationErrors)
	validationErrors = validate.Parameter(params.GasPrice, "gasPrice", validate.CheckBigInt, validationErrors)
	validationErrors = validate.Parameter(params.ComplexityLevel, "complexityLevel", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.Parts, "parts", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.MainRouteParts, "mainRouteParts", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.GasLimit, "gasLimit", validate.CheckFloat32NonNegativeWhole, validationErrors)
	validationErrors = validate.Parameter(params.ConnectorTokens, "connectorTokens", validate.CheckConnectorTokens, validationErrors)
	validationErrors = validate.Parameter(params.Permit, "permit", validate.CheckPermitHash, validationErrors)
	validationErrors = validate.Parameter(params.Receiver, "receiver", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Referrer, "referrer", validate.CheckEthereumAddress, validationErrors)
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
