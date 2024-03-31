package aggregation

import "net/url"

type Configuration struct {
	WalletConfig *WalletConfiguration
	ChainId      uint64

	ApiKey string
	ApiURL *url.URL
}

type WalletConfiguration struct {
	PrivateKey string
	NodeURL    string
}
