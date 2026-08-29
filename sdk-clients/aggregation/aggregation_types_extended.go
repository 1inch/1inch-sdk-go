package aggregation

// Clean names matching their consumer methods (GetApproveAllowance /
// GetApproveTransaction). The generated GetAllowanceParams / GetApproveParams
// remain valid.
type (
	GetApproveAllowanceParams   = GetAllowanceParams
	GetApproveTransactionParams = GetApproveParams
)
