package models

import "time"

type CreateOrderResponse struct {
	Success bool `json:"success"`
}

type OrderResponse struct {
	Signature            string    `json:"signature"`
	OrderHash            string    `json:"orderHash"`
	CreateDateTime       time.Time `json:"createDateTime"`
	RemainingMakerAmount string    `json:"remainingMakerAmount"`
	MakerBalance         string    `json:"makerBalance"`
	MakerAllowance       string    `json:"makerAllowance"`
	Data                 struct {
		MakerAsset    string `json:"makerAsset"`
		TakerAsset    string `json:"takerAsset"`
		Salt          string `json:"salt"`
		Receiver      string `json:"receiver"`
		AllowedSender string `json:"allowedSender"`
		MakingAmount  string `json:"makingAmount"`
		TakingAmount  string `json:"takingAmount"`
		Maker         string `json:"maker"`
		Interactions  string `json:"interactions"`
		Offsets       string `json:"offsets"`
	} `json:"data"`
	MakerRate          string      `json:"makerRate"`
	TakerRate          string      `json:"takerRate"`
	IsMakerContract    bool        `json:"isMakerContract"`
	OrderInvalidReason interface{} `json:"orderInvalidReason"`
}

type CountResponse struct {
	Count int `json:"count"`
}

type EventResponse struct {
	Id                   int       `json:"id"`
	Network              int       `json:"network"`
	LogId                string    `json:"logId"`
	Version              int       `json:"version"`
	Action               string    `json:"action"`
	OrderHash            string    `json:"orderHash"`
	Taker                string    `json:"taker"`
	RemainingMakerAmount string    `json:"remainingMakerAmount"`
	TransactionHash      string    `json:"transactionHash"`
	BlockNumber          int       `json:"blockNumber"`
	CreateDateTime       time.Time `json:"createDateTime"`
}

type BuildMakerTraitsParams struct {
	AllowedSender      string
	ShouldCheckEpoch   bool
	UsePermit2         bool
	UnwrapWeth         bool
	HasExtension       bool
	HasPreInteraction  bool
	HasPostInteraction bool
	Expiry             int64
	Nonce              int64
	Series             int64
}
