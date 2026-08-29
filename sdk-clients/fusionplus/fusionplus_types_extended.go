package fusionplus

import (
	"math/big"

	"github.com/1inch/1inch-sdk-go/v5/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/v5/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/v5/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
)

type GetOrderFillsByHashParams struct {
	Hash string `url:"hash" json:"hash"`
}

// Deprecated: Use GetOrderFillsByHashParams instead (renamed to match the
// GetOrderFillsByHash method and its GetOrderFillsByHashOutput return type).
type GetOrderByOrderHashParams = GetOrderFillsByHashParams
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
	EscrowExtension     *EscrowExtension
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

// Ergonomic re-exports of fusionorder types, so callers of this package need not
// import common/fusionorder directly.
type TakingFeeInfo = fusionorder.TakingFeeInfo
type CustomPreset = fusionorder.CustomPreset
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

// Deprecated: Use the generated GasCostConfig type instead. Now an alias of it
// (the shapes are identical); kept so existing integrations keep compiling.
type GasCostConfigClass = GasCostConfig

// Deprecated: Use the generated AuctionPoint type instead. Now an alias of it
// (the shapes are identical); kept so existing integrations keep compiling.
type AuctionPointClass = AuctionPoint

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

// Extension represents the extension data for the FusionPlus order
// and should be only created using the NewExtension function
type Extension struct {
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

// Deprecated: Use Extension. The Plus suffix is redundant inside the fusionplus package.
type ExtensionPlus = Extension

// Deprecated: Use NewExtension. The Plus suffix is redundant inside the fusionplus package.
func NewExtensionPlus(params ExtensionParamsPlus) (*Extension, error) { return NewExtension(params) }
