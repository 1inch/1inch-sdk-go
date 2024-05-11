// Package txbroadcast provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package txbroadcast

// BroadcastRequest defines model for BroadcastRequest.
type BroadcastRequest struct {
	// RawTransaction The raw transaction data
	RawTransaction string `json:"rawTransaction"`
}

// BroadcastResponse defines model for BroadcastResponse.
type BroadcastResponse struct {
	// TransactionHash The transaction hash
	TransactionHash string `json:"transactionHash"`
}

// TxProcessorApiControllerBroadcastTransactionJSONRequestBody defines body for TxProcessorApiControllerBroadcastTransaction for application/json ContentType.
type TxProcessorApiControllerBroadcastTransactionJSONRequestBody = BroadcastRequest

// TxProcessorApiControllerBroadcastFlashbotsTransactionJSONRequestBody defines body for TxProcessorApiControllerBroadcastFlashbotsTransaction for application/json ContentType.
type TxProcessorApiControllerBroadcastFlashbotsTransactionJSONRequestBody = BroadcastRequest
