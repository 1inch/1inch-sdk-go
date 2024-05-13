package history

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *HistoryEventsByAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.TokenAddress, "TokenAddress", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(params.Address, "Address", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainId, validationErrors)
	//validationErrors = validate.Parameter(params.FromTimestampMs, "FromTimestampMs", validate.CheckFloat32NonNegativeWhole, validationErrors)
	//validationErrors = validate.Parameter(params.ToTimestampMs, "ToTimestampMs", validate.CheckFloat32NonNegativeWhole, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
