package orderbook

import "github.com/1inch/1inch-sdk/golang/client/validate"

type CreateOrderParams struct {
	ChainId      int
	WalletKey    string
	SourceWallet string
	FromToken    string
	ToToken      string
	TakingAmount string
	MakingAmount string
	Receiver     string
	SkipWarnings bool
}

func (params *CreateOrderParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.WalletKey, "walletKey", validate.CheckPrivateKeyRequired, validationErrors)
	validationErrors = validate.Parameter(params.SourceWallet, "sourceWallet", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.FromToken, "fromToken", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.ToToken, "toToken", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.TakingAmount, "takingAmount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(params.MakingAmount, "makingAmount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(params.Receiver, "receiver", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetOrdersByCreatorAddressParams struct {
	ChainId        int
	CreatorAddress string
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

func (params *GetOrdersByCreatorAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.CreatorAddress, "creatorAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Page, "page", validate.CheckPage, validationErrors)
	validationErrors = validate.Parameter(params.Limit, "limit", validate.CheckLimit, validationErrors)
	validationErrors = validate.Parameter(params.Statuses, "statuses", validate.CheckStatusesInts, validationErrors)
	validationErrors = validate.Parameter((string)(params.SortBy), "sortBy", validate.CheckSortBy, validationErrors)
	validationErrors = validate.Parameter(params.TakerAsset, "takerAsset", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.MakerAsset, "makerAsset", validate.CheckEthereumAddress, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetAllOrdersParams struct {
	ChainId int
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

func (params *GetAllOrdersParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.Page, "page", validate.CheckPage, validationErrors)
	validationErrors = validate.Parameter(params.Limit, "limit", validate.CheckLimit, validationErrors)
	validationErrors = validate.Parameter(params.Statuses, "statuses", validate.CheckStatusesInts, validationErrors)
	validationErrors = validate.Parameter((string)(params.SortBy), "sortBy", validate.CheckSortBy, validationErrors)
	validationErrors = validate.Parameter(params.TakerAsset, "takerAsset", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.MakerAsset, "makerAsset", validate.CheckEthereumAddress, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetCountParams struct {
	ChainId int
	LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams
}

func (params *GetCountParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.Statuses, "statuses", validate.CheckStatusesStrings, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetEventParams struct {
	ChainId   int
	OrderHash string
}

func (params *GetEventParams) Validate() error { // TODO Find validation criteria for OrderHash
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.OrderHash, "orderHash", validate.CheckOrderHashRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetEventsParams struct {
	ChainId int
	LimitOrderV3SubscribedApiControllerGetEventsParams
}

func (params *GetEventsParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.Limit, "limit", validate.CheckLimit, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetActiveOrdersWithPermitParams struct {
	ChainId int
	Wallet  string
	Token   string
}

func (params *GetActiveOrdersWithPermitParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainIdRequired, validationErrors)
	validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckPrivateKeyRequired, validationErrors)
	validationErrors = validate.Parameter(params.Token, "token", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
