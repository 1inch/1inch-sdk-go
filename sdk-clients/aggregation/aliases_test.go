package aggregation

// Compile-time assertions that the clean param names alias the generated ones.
var (
	_ GetApproveAllowanceParams   = GetAllowanceParams{}
	_ GetApproveTransactionParams = GetApproveParams{}
)
