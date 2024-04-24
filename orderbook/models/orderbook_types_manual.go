package models

import (
	"math/big"
	"time"
)

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

type GetOrderByHashResponse struct {
	ID                   int         `json:"id"`
	OrderHash            string      `json:"orderHash"`
	CreateDateTime       time.Time   `json:"createDateTime"`
	LastChangedDateTime  time.Time   `json:"lastChangedDateTime"`
	TakerAsset           string      `json:"takerAsset"`
	MakerAsset           string      `json:"makerAsset"`
	OrderMaker           string      `json:"orderMaker"`
	OrderStatus          int         `json:"orderStatus"`
	Signature            string      `json:"signature"`
	MakerAmount          string      `json:"makerAmount"`
	RemainingMakerAmount string      `json:"remainingMakerAmount"`
	MakerBalance         string      `json:"makerBalance"`
	MakerAllowance       string      `json:"makerAllowance"`
	TakerAmount          string      `json:"takerAmount"`
	Data                 Data        `json:"data"`
	MakerRate            string      `json:"makerRate"`
	TakerRate            string      `json:"takerRate"`
	TakerRateDoubled     float64     `json:"takerRateDoubled"`
	OrderHashSelector    int         `json:"orderHashSelector"`
	OrderInvalidReason   interface{} `json:"orderInvalidReason"`
	IsMakerContract      bool        `json:"isMakerContract"`
}

type Data struct {
	Salt         string `json:"salt"`
	Maker        string `json:"maker"`
	Receiver     string `json:"receiver"`
	Extension    string `json:"extension,omitempty"`
	MakerAsset   string `json:"makerAsset"`
	TakerAsset   string `json:"takerAsset"`
	MakerTraits  string `json:"makerTraits"`
	MakingAmount string `json:"makingAmount"`
	TakingAmount string `json:"takingAmount"`
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

type GetOrderByHashResponseExtended struct {
	GetOrderByHashResponse

	LimitOrderDataNormalized NormalizedLimitOrderData
}

type NormalizedLimitOrderData struct {
	Salt         *big.Int
	MakerAsset   *big.Int
	TakerAsset   *big.Int
	Maker        *big.Int
	Receiver     *big.Int
	MakingAmount *big.Int
	TakingAmount *big.Int
	MakerTraits  *big.Int
}
