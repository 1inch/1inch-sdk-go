package orderbook

import "github.com/1inch/1inch-sdk/golang/client/validate"

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
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
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
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
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
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
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
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
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
	if err := validate.Limit(params.Limit, "chainId"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
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
	if err := validate.EthereumAddress(params.Wallet, "wallet"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := validate.EthereumAddress(params.Token, "token"); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if len(validationErrors) > 0 {
		return validate.AggregateValidationErorrs(validationErrors)
	}
	return nil
}
