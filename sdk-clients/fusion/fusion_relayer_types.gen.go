// Package fusion provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package fusion

// OrderInput defines model for OrderInput.
type OrderInput struct {
	// Maker An address of the maker (wallet or contract address)
	Maker string `json:"maker"`

	// MakerAsset Address of the maker asset
	MakerAsset string `json:"makerAsset"`

	// MakerTraits Includes some flags like, allow multiple fills, is partial fill allowed or not, price improvement, nonce, deadline etc
	MakerTraits string `json:"makerTraits,omitempty"`

	// MakingAmount Order maker's token amount
	MakingAmount string `json:"makingAmount"`

	// Receiver An address of the wallet or contract who will receive filled amount
	Receiver string `json:"receiver,omitempty"`
	Salt     string `json:"salt"`

	// TakerAsset Address of the taker asset
	TakerAsset string `json:"takerAsset"`

	// TakingAmount Order taker's token amount
	TakingAmount string `json:"takingAmount"`
}

// SignedOrderInput defines model for SignedOrderInput.
type SignedOrderInput struct {
	// Extension An interaction call data. ABI encoded set of makerAssetSuffix, takerAssetSuffix, makingAmountGetter, takingAmountGetter, predicate, permit, preInteraction, postInteraction.Lowest 160 bits of the order salt must be equal to the lowest 160 bits of the extension hash
	Extension string     `json:"extension,omitempty"`
	Order     OrderInput `json:"order"`

	// QuoteId Quote id of the quote with presets
	QuoteId string `json:"quoteId"`

	// Signature Signature of the gasless order typed data (using signTypedData_v4)
	Signature string `json:"signature"`
}

// RelayerControllerSubmitManyJSONBody defines parameters for RelayerControllerSubmitMany.
type RelayerControllerSubmitManyJSONBody = []SignedOrderInput

// RelayerControllerSubmitJSONRequestBody defines body for RelayerControllerSubmit for application/json ContentType.
type RelayerControllerSubmitJSONRequestBody = SignedOrderInput

// RelayerControllerSubmitManyJSONRequestBody defines body for RelayerControllerSubmitMany for application/json ContentType.
type RelayerControllerSubmitManyJSONRequestBody = RelayerControllerSubmitManyJSONBody
