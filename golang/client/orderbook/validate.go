package orderbook

import (
	"fmt"

	"1inch-sdk-golang/helpers"
)

func (params *LimitOrderV3Request) Validate() error {
	if params.Data == (LimitOrderV3Data{}) {
		return fmt.Errorf("data is required")
	}
	return nil
}

func (params LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams) Validate() error {
	if params.Page != nil {
		if *params.Page < 1 {
			return fmt.Errorf("page must be greater than 0")
		}
	}
	if params.Limit != nil {
		if *params.Limit < 1 {
			return fmt.Errorf("limit must be greater than 0")
		}
	}
	if params.Statuses != nil {
		if helpers.HasDuplicates(*params.Statuses) {
			return fmt.Errorf("statuses must not contain duplicates")
		}
		if !helpers.IsSubset(*params.Statuses, []float32{1, 2, 3}) {
			return fmt.Errorf("statuses can only contain 1, 2, and/or 3")
		}
	}
	if params.SortBy != nil {
		if !helpers.Contains(string(*params.SortBy), []string{"createDateTime", "takerRate", "makerRate", "makerAmount", "takerAmount"}) {
			return fmt.Errorf("sortBy can only contain createDateTime, takerRate, makerRate, makerAmount, or takerAmount")
		}
	}
	if params.TakerAsset != nil {
		if !helpers.IsEthereumAddress(*params.TakerAsset) {
			return fmt.Errorf("takerAsset must be a valid Ethereum address")
		}
	}
	if params.MakerAsset != nil {
		if !helpers.IsEthereumAddress(*params.MakerAsset) {
			return fmt.Errorf("makerAsset must be a valid Ethereum address")
		}
	}
	return nil
}

func (params *LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams) Validate() error {
	if helpers.HasDuplicates(params.Statuses) {
		return fmt.Errorf("statuses must not contain duplicates")
	}
	if !helpers.IsSubset(params.Statuses, []string{"1", "2", "3"}) {
		return fmt.Errorf("statuses can only contain 1, 2, and/or 3")
	}
	return nil
}

func (params *LimitOrderV3SubscribedApiControllerGetEventsParams) Validate() error {
	if params.Limit < 1 {
		return fmt.Errorf("limit must be greater than 0")
	}
	return nil
}
