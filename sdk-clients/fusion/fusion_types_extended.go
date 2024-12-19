package fusion

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

type GetQuoteOutputFixed struct {
	// FeeToken Destination token address
	FeeToken        string                 `json:"feeToken"`
	FromTokenAmount string                 `json:"fromTokenAmount"`
	Presets         QuotePresetsClassFixed `json:"presets"`
	Prices          TokenPairValue         `json:"prices"`

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
	Permit                  string                          `json:"permit,omitempty"`   // without the first 20 bytes of token address
	Receiver                string                          `json:"receiver,omitempty"` // Should be set to the full zero address if this order should be filled by anyone
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

type QuoterControllerGetQuoteParamsFixed struct {
	// FromTokenAddress Address of "FROM" token
	FromTokenAddress string `url:"fromTokenAddress" json:"fromTokenAddress"`

	// ToTokenAddress Address of "TO" token
	ToTokenAddress string `url:"toTokenAddress" json:"toTokenAddress"`

	// Amount to take from "FROM" token to get "TO" token
	Amount string `url:"amount" json:"amount"`

	// WalletAddress An address of the wallet or contract who will create Fusion order
	WalletAddress string `url:"walletAddress" json:"walletAddress"`

	// EnableEstimate if enabled then get estimation from 1inch swap builder and generates quoteId, by default is false
	EnableEstimate bool `url:"enableEstimate" json:"enableEstimate"`

	// Fee in bps format, 1% is equal to 100bps
	Fee float32 `url:"fee,omitempty" json:"fee,omitempty"`

	// IsPermit2 permit2 allowance transfer encoded call
	IsPermit2    string `url:"isPermit2,omitempty" json:"isPermit2,omitempty"`
	IsLedgerLive bool   `url:"isLedgerLive" json:"isLedgerLive"`

	// Permit permit, user approval sign
	Permit string `url:"permit,omitempty" json:"permit,omitempty"`
}

// PresetClassFixed defines model for PresetClass.
type PresetClassFixed struct {
	AllowMultipleFills bool                `json:"allowMultipleFills"`
	AllowPartialFills  bool                `json:"allowPartialFills"`
	AuctionDuration    float32             `json:"auctionDuration"`
	AuctionEndAmount   string              `json:"auctionEndAmount"`
	AuctionStartAmount string              `json:"auctionStartAmount"`
	BankFee            string              `json:"bankFee"`
	EstP               float32             `json:"estP"`
	ExclusiveResolver  string              `json:"exclusiveResolver"` // This was changed to a string from a map[string]interface{}
	GasCost            GasCostConfigClass  `json:"gasCost"`
	InitialRateBump    float32             `json:"initialRateBump"`
	Points             []AuctionPointClass `json:"points"`
	StartAuctionIn     float32             `json:"startAuctionIn"`
	TokenFee           string              `json:"tokenFee"`
}

// QuotePresetsClassFixed defines model for QuotePresetsClass.
type QuotePresetsClassFixed struct {
	Custom *PresetClassFixed `json:"custom,omitempty"`
	Fast   PresetClassFixed  `json:"fast"`
	Medium PresetClassFixed  `json:"medium"`
	Slow   PresetClassFixed  `json:"slow"`
}
