// Package fusion provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package fusion

import (
	"time"
)

// ActiveOrdersOutput defines model for ActiveOrdersOutput.
type ActiveOrdersOutput struct {
	// AuctionEndDate End date of the auction for this order.
	AuctionEndDate time.Time `json:"auctionEndDate"`

	// AuctionStartDate Start date of the auction for this order.
	AuctionStartDate time.Time `json:"auctionStartDate"`

	// Deadline Deadline by which the order must be filled.
	Deadline time.Time `json:"deadline"`

	// Extension An interaction call data. ABI encoded set of makerAssetSuffix, takerAssetSuffix, makingAmountGetter, takingAmountGetter, predicate, permit, preInteraction, postInteraction.If extension exists then lowest 160 bits of the order salt must be equal to the lowest 160 bits of the extension hash
	Extension string        `json:"extension"`
	Order     FusionOrderV4 `json:"order"`

	// OrderHash i.e 0x806039f5149065924ad52de616b50abff488c986716d052e9c160887bc09e559
	OrderHash string `json:"orderHash"`

	// QuoteId Identifier of the quote associated with this order.
	QuoteId string `json:"quoteId"`

	// RemainingMakerAmount Remaining amount of the maker asset that can still be filled.
	RemainingMakerAmount string `json:"remainingMakerAmount"`

	// Signature i.e 0x38de7c8c406c8668eec947d59679028c068735e56c8a41bcc5b3dc2d2229dec258424e0f06b189d2b87f9f3d9cdd9edcb7b3be4108bd8605d052c20c84e65ad61c
	Signature string `json:"signature"`
}

// FusionOrderV4 defines model for FusionOrderV4.
type FusionOrderV4 struct {
	// Maker Address of the account creating the order (maker).
	Maker string `json:"maker"`

	// MakerAsset Identifier of the asset being offered by the maker.
	MakerAsset string `json:"makerAsset"`

	// MakerTraits Includes some flags like, allow multiple fills, is partial fill allowed or not, price improvement, nonce, deadline etc.
	MakerTraits string `json:"makerTraits"`

	// MakingAmount Amount of the makerAsset being offered by the maker.
	MakingAmount string `json:"makingAmount"`

	// Receiver Address of the account receiving the assets (receiver), if different from maker.
	Receiver string `json:"receiver"`

	// Salt Some unique value. It is necessary to be able to create limit orders with the same parameters (so that they have a different hash), Lowest 160 bits of the order salt must be equal to the lowest 160 bits of the extension hash
	Salt string `json:"salt"`

	// TakerAsset Identifier of the asset being requested by the maker in exchange.
	TakerAsset string `json:"takerAsset"`

	// TakingAmount Amount of the takerAsset being requested by the maker.
	TakingAmount string `json:"takingAmount"`
}

// GetActiveOrdersOutput defines model for GetActiveOrdersOutput.
type GetActiveOrdersOutput struct {
	Items []ActiveOrdersOutput `json:"items"`
	Meta  Meta                 `json:"meta"`
}

// Meta defines model for Meta.
type Meta struct {
	CurrentPage  float32 `json:"currentPage"`
	ItemsPerPage float32 `json:"itemsPerPage"`
	TotalItems   float32 `json:"totalItems"`
	TotalPages   float32 `json:"totalPages"`
}

// SettlementAddressOutput defines model for SettlementAddressOutput.
type SettlementAddressOutput struct {
	// Address actual settlement contract address
	Address string `json:"address"`
}

// OrderApiControllerGetActiveOrdersParams defines parameters for OrderApiControllerGetActiveOrders.
type OrderApiControllerGetActiveOrdersParams struct {
	// Page Pagination step, default: 1 (page = offset / limit)
	Page float32 `url:"page,omitempty" json:"page,omitempty"`

	// Limit Number of active orders to receive (default: 100, max: 500)
	Limit float32 `url:"limit,omitempty" json:"limit,omitempty"`
}