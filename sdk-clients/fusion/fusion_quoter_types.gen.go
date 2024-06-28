// Package fusion provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package fusion

// Defines values for GetQuoteOutputRecommendedPreset.
const (
	Custom GetQuoteOutputRecommendedPreset = "custom"
	Fast   GetQuoteOutputRecommendedPreset = "fast"
	Medium GetQuoteOutputRecommendedPreset = "medium"
	Slow   GetQuoteOutputRecommendedPreset = "slow"
)

// AuctionPointClass defines model for AuctionPointClass.
type AuctionPointClass struct {
	Coefficient float32 `json:"coefficient"`
	Delay       float32 `json:"delay"`
}

// CustomPresetInput defines model for CustomPresetInput.
type CustomPresetInput struct {
	AuctionDuration    float32  `json:"auctionDuration"`
	AuctionEndAmount   int64    `json:"auctionEndAmount"`
	AuctionStartAmount int64    `json:"auctionStartAmount"`
	Points             []string `json:"points,omitempty"`
}

// GasCostConfigClass defines model for GasCostConfigClass.
type GasCostConfigClass struct {
	GasBumpEstimate  float32 `json:"gasBumpEstimate"`
	GasPriceEstimate string  `json:"gasPriceEstimate"`
}

// GetQuoteOutput defines model for GetQuoteOutput.
type GetQuoteOutput struct {
	// FeeToken Destination token address
	FeeToken        string            `json:"feeToken"`
	FromTokenAmount string            `json:"fromTokenAmount"`
	Presets         QuotePresetsClass `json:"presets"`
	Prices          TokenPairValue    `json:"prices"`

	// QuoteId Current generated quote id, should be passed with order
	QuoteId map[string]interface{} `json:"quoteId"`

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

// GetQuoteOutputRecommendedPreset suggested to use this preset
type GetQuoteOutputRecommendedPreset string

// PairCurrencyValue defines model for PairCurrencyValue.
type PairCurrencyValue struct {
	FromToken string `json:"fromToken"`
	ToToken   string `json:"toToken"`
}

// PresetClass defines model for PresetClass.
type PresetClass struct {
	AllowMultipleFills bool                   `json:"allowMultipleFills"`
	AllowPartialFills  bool                   `json:"allowPartialFills"`
	AuctionDuration    float32                `json:"auctionDuration"`
	AuctionEndAmount   string                 `json:"auctionEndAmount"`
	AuctionStartAmount string                 `json:"auctionStartAmount"`
	BankFee            string                 `json:"bankFee"`
	EstP               float32                `json:"estP"`
	ExclusiveResolver  map[string]interface{} `json:"exclusiveResolver"`
	GasCost            GasCostConfigClass     `json:"gasCost"`
	InitialRateBump    float32                `json:"initialRateBump"`
	Points             []AuctionPointClass    `json:"points"`
	StartAuctionIn     float32                `json:"startAuctionIn"`
	TokenFee           string                 `json:"tokenFee"`
}

// QuotePresetsClass defines model for QuotePresetsClass.
type QuotePresetsClass struct {
	Custom *PresetClass `json:"custom,omitempty"`
	Fast   PresetClass  `json:"fast"`
	Medium PresetClass  `json:"medium"`
	Slow   PresetClass  `json:"slow"`
}

// TokenPairValue defines model for TokenPairValue.
type TokenPairValue struct {
	Usd PairCurrencyValue `json:"usd"`
}

// QuoterControllerGetQuoteParams defines parameters for QuoterControllerGetQuote.
type QuoterControllerGetQuoteParams struct {
	// FromTokenAddress Address of "FROM" token
	FromTokenAddress string `url:"fromTokenAddress" json:"fromTokenAddress"`

	// ToTokenAddress Address of "TO" token
	ToTokenAddress string `url:"toTokenAddress" json:"toTokenAddress"`

	// Amount Amount to take from "FROM" token to get "TO" token
	Amount float32 `url:"amount" json:"amount"`

	// WalletAddress An address of the wallet or contract who will create Fusion order
	WalletAddress string `url:"walletAddress" json:"walletAddress"`

	// EnableEstimate if enabled then get estimation from 1inch swap builder and generates quoteId, by default is false
	EnableEstimate bool `url:"enableEstimate" json:"enableEstimate"`

	// Fee fee in bps format, 1% is equal to 100bps
	Fee float32 `url:"fee,omitempty" json:"fee,omitempty"`

	// IsPermit2 permit2 allowance transfer encoded call
	IsPermit2    string `url:"isPermit2,omitempty" json:"isPermit2,omitempty"`
	IsLedgerLive bool   `url:"isLedgerLive" json:"isLedgerLive"`

	// Permit permit, user approval sign
	Permit string `url:"permit,omitempty" json:"permit,omitempty"`
}

// QuoterControllerGetQuoteWithCustomPresetsParams defines parameters for QuoterControllerGetQuoteWithCustomPresets.
type QuoterControllerGetQuoteWithCustomPresetsParams struct {
	// FromTokenAddress Address of "FROM" token
	FromTokenAddress string `url:"fromTokenAddress" json:"fromTokenAddress"`

	// ToTokenAddress Address of "TO" token
	ToTokenAddress string `url:"toTokenAddress" json:"toTokenAddress"`

	// Amount Amount to take from "FROM" token to get "TO" token
	Amount float32 `url:"amount" json:"amount"`

	// WalletAddress An address of the wallet or contract who will create Fusion order
	WalletAddress string `url:"walletAddress" json:"walletAddress"`

	// EnableEstimate if enabled then get estimation from 1inch swap builder and generates quoteId, by default is false
	EnableEstimate bool `url:"enableEstimate" json:"enableEstimate"`

	// Fee fee in bps format, 1% is equal to 100bps
	Fee float32 `url:"fee,omitempty" json:"fee,omitempty"`

	// IsPermit2 permit2 allowance transfer encoded call
	IsPermit2 string `url:"isPermit2,omitempty" json:"isPermit2,omitempty"`

	// Permit permit, user approval sign
	Permit string `url:"permit,omitempty" json:"permit,omitempty"`
}

// QuoterControllerGetQuoteWithCustomPresetsJSONRequestBody defines body for QuoterControllerGetQuoteWithCustomPresets for application/json ContentType.
type QuoterControllerGetQuoteWithCustomPresetsJSONRequestBody = CustomPresetInput
