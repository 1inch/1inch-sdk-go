package web3_provider

type TransactionExecutor interface {
}

type Wallet interface {
	Nonce()
	Address()
	Balance()

	Sign(transaction string) (string, error)
	BroadcastTransaction(transaction string) error

	// will generate the data for transaction or transaction itself
	TokenPermit()
	TokenApprove()

	// view functions
	TokenBalance()
	TokenAllowance()
}
