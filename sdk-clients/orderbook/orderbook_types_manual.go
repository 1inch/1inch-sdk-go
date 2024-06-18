package orderbook

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type CreateOrderParams struct {
	SeriesNonce                    *big.Int
	MakerTraits                    *MakerTraits
	Extension                      Extension
	PrivateKey                     string
	ExpireAfter                    int64
	Maker                          string
	MakerAsset                     string
	TakerAsset                     string
	TakingAmount                   string
	MakingAmount                   string
	Taker                          string
	SkipWarnings                   bool
	EnableOnchainApprovalsIfNeeded bool
}

type GetOrdersByCreatorAddressParams struct {
	CreatorAddress string
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

type GetOrderParams struct {
	OrderHash string
}

type GetAllOrdersParams struct {
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

type GetCountParams struct {
	LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams
}

type GetEventParams struct {
	OrderHash string
}

type GetEventsParams struct {
	LimitOrderV3SubscribedApiControllerGetEventsParams
}

type GetActiveOrdersWithPermitParams struct {
	Wallet string
	Token  string
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
	MakerTraits   string `json:"makerTraits"`
	Extension     string `json:"extension"`
}

type CreateOrderResponse struct {
	Success bool `json:"success"`
}

type OrderResponse struct {
	Signature            string      `json:"signature"`
	OrderHash            string      `json:"orderHash"`
	CreateDateTime       time.Time   `json:"createDateTime"`
	RemainingMakerAmount string      `json:"remainingMakerAmount"`
	MakerBalance         string      `json:"makerBalance"`
	MakerAllowance       string      `json:"makerAllowance"`
	Data                 OrderData   `json:"data"`
	MakerRate            string      `json:"makerRate"`
	TakerRate            string      `json:"takerRate"`
	IsMakerContract      bool        `json:"isMakerContract"`
	OrderInvalidReason   interface{} `json:"orderInvalidReason"`
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
	Data                 OrderData   `json:"data"`
	MakerRate            string      `json:"makerRate"`
	TakerRate            string      `json:"takerRate"`
	TakerRateDoubled     float64     `json:"takerRateDoubled"`
	OrderHashSelector    int         `json:"orderHashSelector"`
	OrderInvalidReason   interface{} `json:"orderInvalidReason"`
	IsMakerContract      bool        `json:"isMakerContract"`
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

type TakerTraitsParams struct {
	Receiver        *common.Address
	Extension       string
	MakerAmount     bool
	UnwrapWETH      bool
	SkipOrderPermit bool
	UsePermit2      bool
	ArgsHasReceiver bool
}

type TakerTraitsCalldata struct {
	Trait *big.Int
	Args  string
}
