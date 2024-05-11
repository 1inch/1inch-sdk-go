package txbroadcast

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *BroadcastRequest) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.RawTransaction, "RawTransaction", validate.CheckString, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
