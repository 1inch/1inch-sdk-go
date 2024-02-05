package orderbook

import "github.com/svanas/1inch-sdk/golang/client/validate"

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
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainId, validationErrors)
	validationErrors = validate.Parameter(params.WalletKey, "walletKey", validate.CheckPrivateKey, validationErrors)
	validationErrors = validate.Parameter(params.SourceWallet, "sourceWallet", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.FromToken, "fromToken", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.ToToken, "toToken", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.TakingAmount, "takingAmount", validate.CheckBigInt, validationErrors)
	validationErrors = validate.Parameter(params.MakingAmount, "makingAmount", validate.CheckBigInt, validationErrors)
	validationErrors = validate.Parameter(params.Receiver, "receiver", validate.CheckEthereumAddress, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetOrdersByCreatorAddressParams struct {
	ChainId        int
	CreatorAddress string
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

func (params *GetOrdersByCreatorAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainId, validationErrors)
	validationErrors = validate.Parameter(params.CreatorAddress, "creatorAddress", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Page, "page", validate.CheckPagePointer, validationErrors)
	validationErrors = validate.Parameter(params.Limit, "limit", validate.CheckLimitPointer, validationErrors)
	validationErrors = validate.Parameter(params.Statuses, "statuses", validate.CheckStatusesIntsPointer, validationErrors)
	validationErrors = validate.Parameter((*string)(params.SortBy), "sortBy", validate.CheckSortByPointer, validationErrors)
	validationErrors = validate.Parameter(params.TakerAsset, "takerAsset", validate.CheckEthereumAddressPointer, validationErrors)
	validationErrors = validate.Parameter(params.MakerAsset, "makerAsset", validate.CheckEthereumAddressPointer, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetAllOrdersParams struct {
	ChainId int
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

func (params *GetAllOrdersParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainId, validationErrors)
	validationErrors = validate.Parameter(params.Page, "page", validate.CheckPagePointer, validationErrors)
	validationErrors = validate.Parameter(params.Limit, "limit", validate.CheckLimitPointer, validationErrors)
	validationErrors = validate.Parameter(params.Statuses, "statuses", validate.CheckStatusesIntsPointer, validationErrors)
	validationErrors = validate.Parameter((*string)(params.SortBy), "sortBy", validate.CheckSortByPointer, validationErrors)
	validationErrors = validate.Parameter(params.TakerAsset, "takerAsset", validate.CheckEthereumAddressPointer, validationErrors)
	validationErrors = validate.Parameter(params.MakerAsset, "makerAsset", validate.CheckEthereumAddressPointer, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetCountParams struct {
	ChainId int
	LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams
}

func (params *GetCountParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainId, validationErrors)
	validationErrors = validate.Parameter(params.Statuses, "statuses", validate.CheckStatusesStrings, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetEventParams struct {
	ChainId   int
	OrderHash string
}

func (params *GetEventParams) Validate() error { // TODO Find validation criteria for OrderHash
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainId, validationErrors)
	validationErrors = validate.Parameter(params.OrderHash, "orderHash", validate.CheckOrderHash, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

type GetEventsParams struct {
	ChainId int
	LimitOrderV3SubscribedApiControllerGetEventsParams
}

func (params *GetEventsParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainId, validationErrors)
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
	validationErrors = validate.Parameter(params.ChainId, "chainId", validate.CheckChainId, validationErrors)
	validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckPrivateKey, validationErrors)
	validationErrors = validate.Parameter(params.Token, "token", validate.CheckEthereumAddress, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
