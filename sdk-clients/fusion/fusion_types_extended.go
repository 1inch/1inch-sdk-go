package fusion

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

type GetQuoteOutputFixed struct {
	// FeeToken Destination token address
	FeeToken        string            `json:"feeToken"`
	FromTokenAmount string            `json:"fromTokenAmount"`
	Presets         QuotePresetsClass `json:"presets"`
	Prices          TokenPairValue    `json:"prices"`

	// QuoteId Current generated quote id, should be passed with order
	QuoteId string `json:"quoteId"` // TODO This field is marked as "object" instead of "string" in the swagger file. This is an easy fix from the Fusion team

	// RecommendedPreset suggested to use this preset
	RecommendedPreset GetQuoteOutputRecommendedPreset `json:"recommended_preset"`

	// SettlementAddress settlement contract address
	SettlementAddress string `json:"settlementAddress"`

	// Suggested is it suggested to use Fusion
	Suggested     bool           `json:"suggested"`
	ToTokenAmount string         `json:"toTokenAmount"`
	Volume        TokenPairValue `json:"volume"`

	// Whitelist current executors whitelist addresses
	Whitelist []string `json:"whitelist"`
}

type PlaceOrderBody struct {
	Maker        string
	MakerAsset   string
	MakerTraits  string
	MakingAmount string
	Receiver     string
	TakerAsset   string
	TakingAmount string
}

type Order struct {
	FusionExtension     *Extension
	Inner               orderbook.OrderData
	SettlementExtension common.Address
	OrderInfo           FusionOrderV4
	AuctionDetails      *AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	Extra               ExtraData
}

type OrderParams struct {
	FromTokenAddress        string                          `json:"fromTokenAddress"`
	ToTokenAddress          string                          `json:"toTokenAddress"`
	Amount                  string                          `json:"amount"`
	WalletAddress           string                          `json:"walletAddress"`
	Permit                  string                          `json:"permit,omitempty"` // without the first 20 bytes of token address
	Receiver                string                          `json:"receiver,omitempty"`
	Preset                  GetQuoteOutputRecommendedPreset `json:"preset,omitempty"`
	Nonce                   *big.Int                        `json:"nonce,omitempty"`
	Fee                     TakingFeeInfo                   `json:"fee,omitempty"`
	Source                  string                          `json:"source,omitempty"`
	IsPermit2               bool                            `json:"isPermit2,omitempty"`
	CustomPreset            *CustomPreset                   `json:"customPreset,omitempty"`
	AllowPartialFills       bool                            `json:"allowPartialFills,omitempty"`
	AllowMultipleFills      bool                            `json:"allowMultipleFills,omitempty"`
	DelayAuctionStartTimeBy float32
	OrderExpirationDelay    uint32 // TODO this field is inaccessible in the typescript SDK
}

type TakingFeeInfo struct {
	TakingFeeBps      *big.Int // 100 == 1%
	TakingFeeReceiver common.Address
}

type CustomPreset struct {
	AuctionDuration    int                 `json:"auctionDuration"`
	AuctionStartAmount string              `json:"auctionStartAmount"`
	AuctionEndAmount   string              `json:"auctionEndAmount"`
	Points             []CustomPresetPoint `json:"points,omitempty"`
}

type CustomPresetPoint struct {
	ToTokenAmount string `json:"toTokenAmount"`
	Delay         int    `json:"delay"`
}

type AuctionDetails struct {
	StartTime       uint32                   `json:"startTime"`
	Duration        uint32                   `json:"duration"`
	InitialRateBump uint32                   `json:"initialRateBump"`
	Points          []AuctionPointClassFixed `json:"points"`
	GasCost         GasCostConfigClassFixed  `json:"gasCost"`
}

type AuctionPointClassFixed struct {
	Coefficient uint32 `json:"coefficient"`
	Delay       uint16 `json:"delay"`
}

type GasCostConfigClassFixed struct {
	GasBumpEstimate  uint32 `json:"gasBumpEstimate"`
	GasPriceEstimate uint32 `json:"gasPriceEstimate"`
}

type Preset struct {
	AuctionDuration    *big.Int            `json:"auctionDuration"`
	StartAuctionIn     *big.Int            `json:"startAuctionIn"`
	BankFee            *big.Int            `json:"bankFee"`
	InitialRateBump    *big.Int            `json:"initialRateBump"`
	AuctionStartAmount *big.Int            `json:"auctionStartAmount"`
	AuctionEndAmount   *big.Int            `json:"auctionEndAmount"`
	TokenFee           *big.Int            `json:"tokenFee"`
	Points             []AuctionPointClass `json:"points"`
	GasCostInfo        GasCostConfigClass  `json:"gasCostInfo"`
	ExclusiveResolver  *common.Address     `json:"exclusiveResolver,omitempty"`
	AllowPartialFills  bool                `json:"allowPartialFills"`
	AllowMultipleFills bool                `json:"allowMultipleFills"`
}

type PreparedOrder struct {
	Order   Order  `json:"order"`
	Hash    string `json:"hash"`
	QuoteId string `json:"quoteId"`
}

type AdditionalParams struct {
	NetworkId   int
	FromAddress string
	PrivateKey  string
}

type FusionOrderConstructor struct {
	SettlementExtension common.Address
	OrderInfo           FusionOrderV4
}

type Details struct {
	Auction            *AuctionDetails `json:"auction"`
	Fees               Fees            `json:"fees"`
	Whitelist          []AuctionWhitelistItem
	ResolvingStartTime *big.Int
}

type Fees struct {
	IntFee  IntegratorFee
	BankFee *big.Int
}

type IntegratorFee struct {
	Ratio    *big.Int
	Receiver common.Address
}

type AuctionWhitelistItem struct {
	Address common.Address
	/**
	 * Timestamp in sec at which address can start resolving
	 */
	AllowFrom *big.Int
}

type ExtraParams struct {
	Nonce                *big.Int
	Permit               string
	AllowPartialFills    bool
	AllowMultipleFills   bool
	OrderExpirationDelay uint32
	EnablePermit2        bool
	Source               string
	unwrapWeth           bool
}

type SettlementSuffixData struct {
	Whitelist          []AuctionWhitelistItem
	IntegratorFee      *IntegratorFee
	BankFee            *big.Int
	ResolvingStartTime *big.Int
	CustomReceiver     common.Address
}

type WhitelistItem struct {
	/**
	 * last 10 bytes of address, no 0x prefix
	 */
	AddressHalf string
	/**
	 * Delay from previous resolver in seconds
	 * For first resolver delay from `resolvingStartTime`
	 */
	Delay *big.Int
}

type ExtraData struct {
	UnwrapWETH           bool
	Nonce                *big.Int
	Permit               string
	AllowPartialFills    bool
	AllowMultipleFills   bool
	OrderExpirationDelay uint32
	EnablePermit2        bool
	Source               string
}
