package fusion

import (
	"math/big"

	"github.com/1inch/1inch-sdk-go/v5/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/v5/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
)

// Type aliases for internal use - these types are now in fusionorder
// Ergonomic re-exports of fusionorder types, so callers of this package need not
// import common/fusionorder directly.
type (
	Bps         = fusionorder.Bps
	Interaction = fusionorder.Interaction
)

// Deprecated: Use Quote instead. The spec bugs this type used to
// correct are now fixed at generation time (codegen/overrides.go); the alias
// is kept so existing integrations keep compiling.
type GetQuoteOutputFixed = Quote

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

// Ergonomic re-exports of fusionorder types, so callers of this package need not
// import common/fusionorder directly.
type TakingFeeInfo = fusionorder.TakingFeeInfo
type CustomPreset = fusionorder.CustomPreset
type CustomPresetPoint = fusionorder.CustomPresetPoint

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

type OrderConstructor struct {
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

// Deprecated: Use QuoteParams instead. The spec bugs this
// type used to correct are now fixed at generation time
// (codegen/overrides.go); the alias is kept so existing integrations keep
// compiling.
type QuoterControllerGetQuoteParamsFixed = QuoteParams

// Deprecated: Use CustomPresetQuoteParams instead;
// see QuoterControllerGetQuoteParamsFixed.
type QuoterControllerGetQuoteWithCustomPresetsParamsFixed = CustomPresetQuoteParams

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
	OrderHash         string                     `json:"orderHash"`
	Points            []fusionorder.AuctionPoint `json:"points"`
	Status            string                     `json:"status"`
	ToTokenToUsdPrice string                     `json:"toTokenToUsdPrice"`
}

// Deprecated: Use Preset instead. The exclusiveResolver type bug this type used
// to correct is now fixed at generation time (codegen/overrides.go); the alias
// is kept so existing integrations keep compiling.
type PresetClassFixed = PresetClass

// Deprecated: Use QuotePresets instead. The alias is kept so existing
// integrations keep compiling.
type QuotePresetsClassFixed = QuotePresetsClass

// Clean, fusionplus-consistent names for the fusion quoter types. The generated
// types keep the upstream "Class" suffix (renaming them via x-go-name would also
// rename every struct field that references them), so these aliases give callers
// the clean names — prefer these over the *Class names.
type (
	Preset        = PresetClass
	QuotePresets  = QuotePresetsClass
	AuctionPoint  = AuctionPointClass
	GasCostConfig = GasCostConfigClass
)

// Deprecated: Use OrderConstructor. The Fusion prefix stutters with the package name.
type FusionOrderConstructor = OrderConstructor
