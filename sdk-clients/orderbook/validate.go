package orderbook

import (
	"strings"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *CreateOrderParams) Validate() error {
	var validationErrors []error
	// validationErrors = validate.Parameter(params.SeriesNonce, "seriesNonce", validate.CheckBigIntRequired, validationErrors) // TODO All other places expect a string value for raw request parameters, but this value will always come in as a big.Int because of the onchain nature of retrieving it
	//validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckWalletRequired, validationErrors) // TODO enable after parameter validation is fixed
	validationErrors = validate.Parameter(params.Maker, "maker", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.ExpireAfter, "expireAfter", validate.CheckExpireAfter, validationErrors)
	validationErrors = validate.Parameter(params.MakerAsset, "makerAsset", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.TakerAsset, "takerAsset", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.TakingAmount, "takingAmount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(params.MakingAmount, "makingAmount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(params.Taker, "taker", validate.CheckEthereumAddressRequired, validationErrors)
	if strings.EqualFold(params.MakerAsset, params.TakerAsset) && (params.MakerAsset != "" && params.TakerAsset != "") {
		validationErrors = append(validationErrors, validate.NewParameterCustomError("maker asset and taker asset cannot be the same"))
	}
	if strings.EqualFold(params.MakerAsset, constants.NativeToken) || strings.EqualFold(params.TakerAsset, constants.NativeToken) {
		validationErrors = append(validationErrors, validate.NewParameterCustomError("native gas token is not supported as maker or taker asset"))
	}

	//TODO if an extension is present, then MakerTraits must also be marked for an extension in the order

	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *GetOrdersByCreatorAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.CreatorAddress, "creatorAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Page, "page", validate.CheckPage, validationErrors)
	validationErrors = validate.Parameter(params.Limit, "limit", validate.CheckLimit, validationErrors)
	validationErrors = validate.Parameter(params.Statuses, "statuses", validate.CheckStatusesInts, validationErrors)
	validationErrors = validate.Parameter((string)(params.SortBy), "sortBy", validate.CheckSortBy, validationErrors)
	validationErrors = validate.Parameter(params.TakerAsset, "takerAsset", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.MakerAsset, "makerAsset", validate.CheckEthereumAddress, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *GetOrderParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.OrderHash, "orderHash", validate.CheckOrderHashRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *GetAllOrdersParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Page, "page", validate.CheckPage, validationErrors)
	validationErrors = validate.Parameter(params.Limit, "limit", validate.CheckLimit, validationErrors)
	validationErrors = validate.Parameter(params.Statuses, "statuses", validate.CheckStatusesInts, validationErrors)
	validationErrors = validate.Parameter((string)(params.SortBy), "sortBy", validate.CheckSortBy, validationErrors)
	validationErrors = validate.Parameter(params.TakerAsset, "takerAsset", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.MakerAsset, "makerAsset", validate.CheckEthereumAddress, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *GetCountParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Statuses, "statuses", validate.CheckStatusesStrings, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *GetEventParams) Validate() error { // TODO Find validation criteria for OrderHash
	var validationErrors []error
	validationErrors = validate.Parameter(params.OrderHash, "orderHash", validate.CheckOrderHashRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *GetEventsParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Limit, "limit", validate.CheckLimit, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *GetActiveOrdersWithPermitParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckPrivateKeyRequired, validationErrors)
	validationErrors = validate.Parameter(params.Token, "token", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
