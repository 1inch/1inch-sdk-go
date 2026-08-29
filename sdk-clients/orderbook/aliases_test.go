package orderbook

// Compile-time assertions that the deprecated names still resolve.
var (
	_ GetSaltParams       = GenerateSaltWithFeesParams{}
	_ WalletConfiguration = ConfigurationWallet{}
)
