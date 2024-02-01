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
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.PrivateKey(params.WalletKey, "walletKey"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.SourceWallet, "sourceWallet"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.FromToken, "fromToken"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.ToToken, "toToken"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.BigInt(params.TakingAmount, "takingAmount"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.BigInt(params.MakingAmount, "makingAmount"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.Receiver, "receiver"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validate.AggregateValidationErorrs(validationErrors)
}

type GetOrdersByCreatorAddressParams struct {
	ChainId        int
	CreatorAddress string
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

func (params *GetOrdersByCreatorAddressParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.CreatorAddress, "creatorAddress"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.PagePointer(params.Page, "page"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.LimitPointer(params.Limit, "limit"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.StatusesIntsPointer(params.Statuses, "statuses"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.SortByPointer((*string)(params.SortBy), "sortBy"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddressPointer(params.TakerAsset, "takerAsset"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddressPointer(params.MakerAsset, "makerAsset"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validate.AggregateValidationErorrs(validationErrors)
}

type GetAllOrdersParams struct {
	ChainId int
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

func (params *GetAllOrdersParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.PagePointer(params.Page, "page"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.LimitPointer(params.Limit, "limit"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.StatusesIntsPointer(params.Statuses, "statuses"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.SortByPointer((*string)(params.SortBy), "sortBy"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddressPointer(params.TakerAsset, "takerAsset"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddressPointer(params.MakerAsset, "makerAsset"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validate.AggregateValidationErorrs(validationErrors)
}

type GetCountParams struct {
	ChainId int
	LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams
}

func (params *GetCountParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if err := validate.StatusesStrings(params.Statuses, "statuses"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validate.AggregateValidationErorrs(validationErrors)
}

type GetEventParams struct {
	ChainId   int
	OrderHash string
}

func (params *GetEventParams) Validate() error { // TODO Find validation criteria for OrderHash
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.OrderHash(params.OrderHash, "orderHash"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validate.AggregateValidationErorrs(validationErrors)
}

type GetEventsParams struct {
	ChainId int
	LimitOrderV3SubscribedApiControllerGetEventsParams
}

func (params *GetEventsParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.Limit(params.Limit, "limit"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validate.AggregateValidationErorrs(validationErrors)
}

type GetActiveOrdersWithPermitParams struct {
	ChainId int
	Wallet  string
	Token   string
}

func (params *GetActiveOrdersWithPermitParams) Validate() error {
	var validationErrors []error
	if err := validate.ChainId(params.ChainId, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.PrivateKey(params.Wallet, "wallet"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.Token, "token"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validate.AggregateValidationErorrs(validationErrors)
}
