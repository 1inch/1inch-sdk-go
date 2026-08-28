package fusionplus

import (
	"math/big"

	"github.com/1inch/1inch-sdk-go/v4/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
)

type GetOrderByOrderHashParams struct {
	Hash string `url:"hash" json:"hash"`
}
type GetReadyToAcceptFillsParams struct {
	Hash string `url:"hash" json:"hash"`
}

// Deprecated: Use GetOrderFillsByHashOutput instead. The type bugs this type
// used to correct (string USD prices, points being an array) are now fixed at
// generation time (codegen/overrides.go); the alias is kept so existing
// integrations keep compiling.
type GetOrderFillsByHashOutputFixed = GetOrderFillsByHashOutput

// Deprecated: Use QuoteParams instead. The type bugs this
// type used to correct (string amount, integer fee, boolean isPermit2) are now
// fixed at generation time (codegen/overrides.go); the alias is kept so
// existing integrations keep compiling.
type QuoterControllerGetQuoteParamsFixed = QuoteParams

// Deprecated: Use CustomPresetQuoteParams instead;
// see QuoterControllerGetQuoteParamsFixed.
type QuoterControllerGetQuoteWithCustomPresetsParamsFixed = CustomPresetQuoteParams

// Deprecated: Use Quote instead. The quoteId type bug this type used
// to correct is now fixed at generation time (codegen/overrides.go); the alias
// is kept so existing integrations keep compiling.
type GetQuoteOutputFixed = Quote

type Order struct {
	EscExtension        *EscrowExtension
	Inner               orderbook.OrderData
	SettlementExtension common.Address
	OrderInfo           CrossChainOrderDto
	AuctionDetails      *fusionorder.AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	Extra               fusionorder.ExtraData
}

type EscrowExtensionParams struct {
	fusion.ExtensionParams
	ExtensionParamsPlus
	HashLock         *HashLock
	DstChainId       int
	DstToken         common.Address
	SrcSafetyDeposit string
	DstSafetyDeposit string
	TimeLocks        TimeLocks
}

type CrossChainOrderParams struct {
	HashLock                *HashLock
	Preset                  GetQuoteOutputRecommendedPreset
	Receiver                string
	Nonce                   *big.Int
	Permit                  string
	IsPermit2               bool
	TakingFeeReceiver       string
	DelayAuctionStartTimeBy float32
	/**
	 * Order will expire in `orderExpirationDelay` after auction ends
	 * Default 12s
	 */
	OrderExpirationDelay uint32
}

type OrderParams struct {
	HashLock          *HashLock
	SecretHashes      []string
	Permit            string
	Receiver          string
	Preset            GetQuoteOutputRecommendedPreset
	Nonce             *big.Int
	Fee               TakingFeeInfo
	Source            string
	IsPermit2         bool
	TakingFeeReceiver string
	CustomPreset      CustomPreset
}

// Deprecated: Use fusionorder.TakingFeeInfo directly instead.
type TakingFeeInfo = fusionorder.TakingFeeInfo

// Deprecated: Use fusionorder.CustomPreset directly instead.
type CustomPreset = fusionorder.CustomPreset

// Deprecated: Use fusionorder.CustomPresetPoint directly instead.
type CustomPresetPoint = fusionorder.CustomPresetPoint

type PreparedOrder struct {
	Order      Order  `json:"order"`
	Hash       string `json:"hash"`
	QuoteId    string `json:"quoteId"`
	LimitOrder *orderbook.Order
}

type AdditionalParams struct {
	NetworkId   int
	FromAddress string
	PrivateKey  string
}

type Details struct {
	Auction            *fusionorder.AuctionDetails `json:"auction"`
	Fees               Fees                        `json:"fees"`
	Whitelist          []fusionorder.AuctionWhitelistItem
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
	Whitelist          []fusionorder.AuctionWhitelistItem
	IntegratorFee      *IntegratorFee
	BankFee            *big.Int
	ResolvingStartTime *big.Int
	CustomReceiver     common.Address
}

// Deprecated: Use the generated Preset type instead. This legacy shape is no
// longer produced or consumed by the SDK (CreateAuctionDetailsPlus now takes
// *Preset); it is kept so existing integrations keep compiling.
type PresetClassFixed struct {
	AllowMultipleFills bool                `json:"allowMultipleFills"`
	AllowPartialFills  bool                `json:"allowPartialFills"`
	AuctionDuration    float32             `json:"auctionDuration"`
	AuctionEndAmount   string              `json:"auctionEndAmount"`
	AuctionStartAmount string              `json:"auctionStartAmount"`
	BankFee            string              `json:"bankFee"`
	EstP               float32             `json:"estP"`
	ExclusiveResolver  string              `json:"exclusiveResolver"` // This was changed to a string from a map[string]any
	GasCost            GasCostConfigClass  `json:"gasCost"`
	InitialRateBump    float32             `json:"initialRateBump"`
	Points             []AuctionPointClass `json:"points"`
	StartAuctionIn     float32             `json:"startAuctionIn"`
	TokenFee           string              `json:"tokenFee"`
}

// Deprecated: Use the generated GasCostConfig type instead; this legacy shape
// exists only as part of PresetClassFixed.
type GasCostConfigClass struct {
	GasBumpEstimate  float32 `json:"gasBumpEstimate"`
	GasPriceEstimate string  `json:"gasPriceEstimate"`
}

// Deprecated: Use the generated AuctionPoint type instead; this legacy shape
// exists only as part of PresetClassFixed.
type AuctionPointClass struct {
	Coefficient float32 `json:"coefficient"`
	Delay       float32 `json:"delay"`
}

// FusionOrderV4 defines model for FusionOrderV4.
type FusionOrderV4 struct {
	// Maker Address of the account creating the order (maker).
	Maker string `json:"maker"`

	// MakerAsset Identifier of the asset being offered by the maker.
	MakerAsset string `json:"makerAsset"`

	// MakerTraits Includes some flags like, allow multiple fills, is partial fill allowed or not, price improvement, nonce, deadline etc.
	MakerTraits string `json:"makerTraits"`

	// MakingAmount Amount of the makerAsset being offered by the maker.
	MakingAmount string `json:"makingAmount"`

	// Receiver Address of the account receiving the assets (receiver), if different from maker.
	Receiver string `json:"receiver"`

	// Salt Some unique value. It is necessary to be able to create limit orders with the same parameters (so that they have a different hash), Lowest 160 bits of the order salt must be equal to the lowest 160 bits of the extension hash
	Salt string `json:"salt"`

	// TakerAsset Identifier of the asset being requested by the maker in exchange.
	TakerAsset string `json:"takerAsset"`

	// TakingAmount Amount of the takerAsset being requested by the maker.
	TakingAmount string `json:"takingAmount"`
}

type ExtensionParamsPlus struct {
	SettlementContract  string
	AuctionDetails      *fusionorder.AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	Asset               string
	Permit              string

	MakerAssetSuffix string
	TakerAssetSuffix string
	Predicate        string
	PreInteraction   string
	CustomData       string
}

// ExtensionPlus represents the extension data for the FusionPlus order
// and should be only created using the NewExtensionPlus function
type ExtensionPlus struct {
	// Raw unencoded data
	SettlementContract  string
	AuctionDetails      *fusionorder.AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	Asset               string
	Permit              string

	// Data formatted for Limit Order Extension
	MakerAssetSuffix string
	TakerAssetSuffix string
	MakingAmountData string
	TakingAmountData string
	Predicate        string
	MakerPermit      string
	PreInteraction   string
	PostInteraction  string
	CustomData       string
}
