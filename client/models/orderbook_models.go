package models

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
	MakerTraits   string `json:"makerTraits"`
	Extension     string `json:"extension"`
}
