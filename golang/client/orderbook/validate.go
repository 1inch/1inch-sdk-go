package orderbook

import (
	clienterrors "1inch-sdk-golang/client/errors"
	"1inch-sdk-golang/helpers"
)

func (params *LimitOrderV3Request) Validate() error {
	if params.Data == (LimitOrderV3Data{}) {
		return clienterrors.NewRequestValidationError("data is required")
	}
	return nil
}

func (params LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams) Validate() error {
	if params.Page != nil {
		if *params.Page < 1 {
			return clienterrors.NewRequestValidationError("page must be greater than 0")
		}
	}
	if params.Limit != nil {
		if *params.Limit < 1 {
			return clienterrors.NewRequestValidationError("limit must be greater than 0")
		}
	}
	if params.Statuses != nil {
		if helpers.HasDuplicates(*params.Statuses) {
			return clienterrors.NewRequestValidationError("statuses must not contain duplicates")
		}
		if !helpers.IsSubset(*params.Statuses, []float32{1, 2, 3}) {
			return clienterrors.NewRequestValidationError("statuses can only contain 1, 2, and/or 3")
		}
	}
	if params.SortBy != nil {
		if !helpers.Contains(string(*params.SortBy), []string{"createDateTime", "takerRate", "makerRate", "makerAmount", "takerAmount"}) {
			return clienterrors.NewRequestValidationError("sortBy can only contain createDateTime, takerRate, makerRate, makerAmount, or takerAmount")
		}
	}
	if params.TakerAsset != nil {
		if !helpers.IsEthereumAddress(*params.TakerAsset) {
			return clienterrors.NewRequestValidationError("takerAsset must be a valid Ethereum address")
		}
	}
	if params.MakerAsset != nil {
		if !helpers.IsEthereumAddress(*params.MakerAsset) {
			return clienterrors.NewRequestValidationError("makerAsset must be a valid Ethereum address")
		}
	}
	return nil
}

func (params *LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams) Validate() error {
	if helpers.HasDuplicates(params.Statuses) {
		return clienterrors.NewRequestValidationError("statuses must not contain duplicates")
	}
	if !helpers.IsSubset(params.Statuses, []string{"1", "2", "3"}) {
		return clienterrors.NewRequestValidationError("statuses can only contain 1, 2, and/or 3")
	}
	return nil
}

func (params *LimitOrderV3SubscribedApiControllerGetEventsParams) Validate() error {
	if params.Limit < 1 {
		return clienterrors.NewRequestValidationError("limit must be greater than 0")
	}
	return nil
}
