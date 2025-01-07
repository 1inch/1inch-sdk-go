package nft

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *GetNftsByAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Address, "Address", validate.CheckEthereumAddressRequired, validationErrors)
	for _, v := range params.ChainIds {
		validationErrors = validate.Parameter(int(v), "ChainId", validate.CheckChainIdIntRequired, validationErrors)
	}
	return validate.ConsolidateValidationErorrs(validationErrors)
}
