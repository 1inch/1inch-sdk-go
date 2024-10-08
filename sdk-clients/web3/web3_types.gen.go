// Package web3 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package web3

const (
	ApiKeyAuthScopes = "ApiKeyAuth.Scopes"
)

// Defines values for PostChainIdNodeTypeParamsNodeType.
const (
	Archive PostChainIdNodeTypeParamsNodeType = "archive"
	Full    PostChainIdNodeTypeParamsNodeType = "full"
)

// PostChainIdJSONBody defines parameters for PostChainId.
type PostChainIdJSONBody struct {
	// Id An identifier established by the client
	Id string `json:"id,omitempty"`

	// Jsonrpc JSON-RPC version, typically "2.0"
	Jsonrpc string `json:"jsonrpc"`

	// Method The name of the RPC method to be invoked
	Method string `json:"method"`

	// Params Parameters for the RPC method
	Params []string `json:"params,omitempty"`
}

// PostChainIdNodeTypeJSONBody defines parameters for PostChainIdNodeType.
type PostChainIdNodeTypeJSONBody struct {
	// Id An identifier established by the client
	Id int `json:"id,omitempty"`

	// Jsonrpc JSON-RPC version, typically "2.0"
	Jsonrpc string `json:"jsonrpc,omitempty"`

	// Method The name of the RPC method to be invoked
	Method string `json:"method,omitempty"`

	// Params Parameters for the RPC method
	Params []string `json:"params,omitempty"`
}

// PostChainIdNodeTypeParamsNodeType defines parameters for PostChainIdNodeType.
type PostChainIdNodeTypeParamsNodeType string

// PostChainIdJSONRequestBody defines body for PostChainId for application/json ContentType.
type PostChainIdJSONRequestBody PostChainIdJSONBody

// PostChainIdNodeTypeJSONRequestBody defines body for PostChainIdNodeType for application/json ContentType.
type PostChainIdNodeTypeJSONRequestBody PostChainIdNodeTypeJSONBody
