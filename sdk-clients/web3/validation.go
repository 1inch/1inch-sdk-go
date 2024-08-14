package web3

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

//TODO add in a validation type for jsonrpc version numbering

func (params *PerformRpcCallAgainstFullNodeParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Jsonrpc, "Jsonrpc", validate.CheckString, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *PerformRpcCallParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Jsonrpc, "Jsonrpc", validate.CheckString, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
