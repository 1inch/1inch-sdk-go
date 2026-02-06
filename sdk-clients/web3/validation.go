package web3

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *PerformRpcCallAgainstFullNodeParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Jsonrpc, "Jsonrpc", validate.CheckJsonRpcVersionRequired, validationErrors)
	validationErrors = validate.Parameter(params.Method, "Method", validate.CheckStringRequired, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *PerformRpcCallParams) Validate() error {
	var validationErrors []error
	//validationErrors = validate.Parameter(params.PostChainIdNodeTypeParamsNodeType, "NodeType", validate.CheckNodeType, validationErrors) // TODO Cannot validate these types due to the way oapi-codegen types them
	validationErrors = validate.Parameter(params.Jsonrpc, "Jsonrpc", validate.CheckJsonRpcVersionRequired, validationErrors)
	validationErrors = validate.Parameter(params.Method, "Method", validate.CheckStringRequired, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}
