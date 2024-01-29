package orderbook

type OrderRequest struct {
	ChainId      int
	WalletKey    string
	SourceWallet string `json:"sourceWallet" validate:"required,eth_addr"`
	FromToken    string `json:"fromToken" validate:"required,eth_addr"`
	ToToken      string `json:"toToken" validate:"required,eth_addr"`
	TakingAmount string `json:"takingAmount"`
	MakingAmount string `json:"makingAmount"`
	Receiver     string `json:"receiver" validate:"omitempty,eth_addr"`
	SkipWarnings bool   `json:"skipWarnings"`
}

type Order struct {
	OrderHash string    `json:"orderHash"`
	Signature string    `json:"signature"`
	Data      OrderData `json:"data"`
}

type OrderData struct {
	MakerAsset    string `json:"makerAsset"`
	TakerAsset    string `json:"takerAsset"`
	MakingAmount  string `json:"makingAmount"`
	TakingAmount  string `json:"takingAmount"`
	Salt          string `json:"salt"`
	Maker         string `json:"maker"`
	AllowedSender string `json:"allowedSender"`
	Receiver      string `json:"receiver"`
	Offsets       string `json:"offsets"`
	Interactions  string `json:"interactions"`
}

type LimitOrderV3DomainData struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainId           int    `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}
