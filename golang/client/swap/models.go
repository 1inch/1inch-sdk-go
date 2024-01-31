package swap

type ExecuteSwapConfig struct {
	WalletKey          string
	ChainId            int
	PublicAddress      string
	FromToken          string
	ToToken            string
	Amount             string
	Slippage           float32
	EstimatedAmountOut string
	TransactionData    string
	IsPermitSwap       bool
	SkipWarnings       bool
}

type ApprovalType int

const (
	PermitIfPossible ApprovalType = iota
	PermitAlways
	ApprovalAlways
)