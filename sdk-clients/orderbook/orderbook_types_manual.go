package orderbook

import (
	"math/big"
	"time"

	"github.com/1inch/1inch-sdk-go/common"
	geth_common "github.com/ethereum/go-ethereum/common"
)

type CreateOrderParams struct {
	Wallet                         common.Wallet
	SeriesNonce                    *big.Int
	MakerTraits                    *MakerTraits
	Extension                      Extension
	ExtensionEncoded               string
	Salt                           string
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
	OrderHash               string
	SleepBetweenSubrequests bool // For free accounts, this should be set to true to avoid 429 errors when using the GetOrderWithSignature method
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
	AllowedSender string `json:"allowedSender,omitempty"`
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

type OrderResponseExtended struct {
	OrderResponse
	LimitOrderDataNormalized NormalizedLimitOrderData
}

type GetOrderByHashResponse struct {
	OrderHash            string    `json:"orderHash"`
	CreateDateTime       time.Time `json:"createDateTime"`
	Signature            string    `json:"signature"`
	OrderStatus          int       `json:"orderStatus"`
	RemainingMakerAmount string    `json:"remainingMakerAmount"`
	MakerBalance         string    `json:"makerBalance"`
	MakerAllowance       string    `json:"makerAllowance"`
	Data                 OrderData `json:"data"`
	MakerRate            string    `json:"makerRate"`
	TakerRate            string    `json:"takerRate"`
	OrderInvalidReason   string    `json:"orderInvalidReason"`
	IsMakerContract      bool      `json:"isMakerContract"`
	Events               string    `json:"events"`
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

type Orders struct {
	Meta struct {
		HasMore    bool   `json:"hasMore"`
		NextCursor string `json:"nextCursor"`
		Count      int    `json:"count"`
	} `json:"meta"`
	Items []GetOrderByHashResponse `json:"items"`
}

type OrderExtendedWithSignature struct {
	GetOrderByHashResponse
	LimitOrderDataNormalized NormalizedLimitOrderData
	Signature                string
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
	Receiver        *geth_common.Address
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

type GetFeeInfoParams struct {
	MakerAsset  string `url:"makerAsset"`
	TakerAsset  string `url:"takerAsset"`
	MakerAmount string `url:"makerAmount"`
	TakerAmount string `url:"takerAmount"`
}

type FeeInfoResponse struct {
	Whitelist                map[string]string `json:"whitelist"`
	FeeBps                   int               `json:"feeBps"`                   // Fee in basis points (e.g., 50 = 0.5%)
	WhitelistDiscountPercent int               `json:"whitelistDiscountPercent"` // Discount percentage for whitelisted resolvers (e.g., 50 = 50% off)
	ProtocolFeeReceiver      string            `json:"protocolFeeReceiver"`
	ExtensionAddress         string            `json:"extensionAddress"`
}

type OrderStatus int

const (
	ValidOrders              OrderStatus = 1
	TemporarilyInvalidOrders OrderStatus = 2
	InvalidOrders            OrderStatus = 3
)

type GetOrderCountParams struct {
	Statuses   []OrderStatus `url:"Statuses"`
	TakerAsset string        `url:"TakerAsset"`
	MakerAsset string        `url:"MakerAsset"`
}

type GetOrderCountResponse struct {
	Count int `json:"count"`
}

type IntegratorFee struct {
	Integrator string
	Protocol   string
	Fee        int // Fee in basis points (e.g., 1 = 0.01%, 100 = 1%)
	Share      int // Integrator's share in basis points (e.g., 1 = 0.01%, 100 = 1%)
}

type ResolverFee struct {
	Receiver          string
	Fee               int // Fee in basis points (e.g., 1 = 0.01%, 100 = 1%)
	WhitelistDiscount int // Discount percentage for whitelisted addresses (0-100)
}

type buildFeePostInteractionDataParams struct {
	CustomReceiver         bool
	CustomReceiverAddress  string
	IntegratorFee          *IntegratorFee
	ResolverFee            *ResolverFee
	Whitelist              []string
	ExtraInteractionTarget string
	ExtraInteractionData   []byte
}

type BuildOrderExtensionBytesParams struct {
	ExtensionTarget  string
	IntegratorFee    *IntegratorFee
	ResolverFee      *ResolverFee
	Whitelist        map[string]string
	MakerPermit      []byte
	CustomReceiver   string
	ExtraInteraction []byte
	CustomData       []byte
}
