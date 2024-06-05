package fusion

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *OrderApiControllerGetActiveOrdersParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Page, "Page", validate.CheckPage, validationErrors)
	validationErrors = validate.Parameter(params.Limit, "Limit", validate.CheckLimit, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
