package fusion

import (
	"math/big"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
)

// Type aliases for internal use - these types are now in fusionorder
// Users should import from fusionorder directly for new code
type (
	Bps         = fusionorder.Bps
	Interaction = fusionorder.Interaction
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

	SurplusFee   float32 `json:"surplusFee"`
	MarketAmount string  `json:"marketAmount"`
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
	AuctionDetails      *fusionorder.AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	Extra               fusionorder.ExtraData
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

// TakingFeeInfo is an alias for the shared fusionorder.TakingFeeInfo type
type TakingFeeInfo = fusionorder.TakingFeeInfo

// CustomPreset is an alias for the shared fusionorder.CustomPreset type
type CustomPreset = fusionorder.CustomPreset

// CustomPresetPoint is an alias for the shared fusionorder.CustomPresetPoint type
type CustomPresetPoint = fusionorder.CustomPresetPoint


type Preset struct {
	AuctionDuration    float32             `json:"auctionDuration"`
	StartAuctionIn     float32             `json:"startAuctionIn"`
	BankFee            *big.Int            `json:"bankFee"`
	InitialRateBump    float32             `json:"initialRateBump"`
	AuctionStartAmount string              `json:"auctionStartAmount"`
	AuctionEndAmount   string              `json:"auctionEndAmount"`
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
	Auction            *fusionorder.AuctionDetails `json:"auction"`
	Whitelist          []fusionorder.AuctionWhitelistItem
	ResolvingStartTime *big.Int
	FeesIntAndRes      *FeesIntegratorAndResolver
}

type FeesIntegratorAndResolver struct {
	Resolver   ResolverFee
	Integrator IntegratorFee
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
	Permit  string `url:"permit,omitempty" json:"permit,omitempty"`
	Surplus bool   `url:"surplus,omitempty" json:"surplus,omitempty"`
}

type QuoterControllerGetQuoteWithCustomPresetsParamsFixed struct {
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
	Permit  string `url:"permit,omitempty" json:"permit,omitempty"`
	Surplus bool   `url:"surplus,omitempty" json:"surplus,omitempty"`
}

type OrderResponse struct {
	ApproximateTakingAmount string  `json:"approximateTakingAmount"`
	AuctionDuration         int     `json:"auctionDuration"`
	AuctionStartDate        int64   `json:"auctionStartDate"`
	CancelTx                *string `json:"cancelTx"`
	CreatedAt               string  `json:"createdAt"`
	Extension               string  `json:"extension"`
	Fills                   []struct {
		FilledAuctionTakerAmount string `json:"filledAuctionTakerAmount"`
		FilledMakerAmount        string `json:"filledMakerAmount"`
		TxHash                   string `json:"txHash"`
	} `json:"fills"`
	FromTokenToUsdPrice string `json:"fromTokenToUsdPrice"`
	InitialRateBump     int    `json:"initialRateBump"`
	IsNativeCurrency    bool   `json:"isNativeCurrency"`
	Order               struct {
		Maker        string `json:"maker"`
		MakerAsset   string `json:"makerAsset"`
		MakerTraits  string `json:"makerTraits"`
		MakingAmount string `json:"makingAmount"`
		Receiver     string `json:"receiver"`
		Salt         string `json:"salt"`
		TakerAsset   string `json:"takerAsset"`
		TakingAmount string `json:"takingAmount"`
	} `json:"order"`
	OrderHash         string                   `json:"orderHash"`
	Points            []fusionorder.AuctionPointClassFixed `json:"points"`
	Status            string                   `json:"status"`
	ToTokenToUsdPrice string                   `json:"toTokenToUsdPrice"`
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
