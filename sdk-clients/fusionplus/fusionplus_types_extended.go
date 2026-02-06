package fusionplus

import (
	"math/big"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
)

type GetOrderByOrderHashParams struct {
	Hash string `url:"hash" json:"hash"`
}
type GetReadyToAcceptFillsParams struct {
	Hash string `url:"hash" json:"hash"`
}

// GetOrderFillsByHashOutputFixed replaces the DstTokenPriceUsd and SrcTokenPriceUsd fields with string and changes Points to be an array
type GetOrderFillsByHashOutputFixed struct {
	// ApproximateTakingAmount Approximate amount of the takerAsset being requested by the maker in dst chain.
	ApproximateTakingAmount string `json:"approximateTakingAmount"`

	// AuctionDuration Unix timestamp in milliseconds
	AuctionDuration float32 `json:"auctionDuration"`

	// AuctionStartDate Unix timestamp in milliseconds
	AuctionStartDate float32                `json:"auctionStartDate"`
	CancelTx         map[string]any `json:"cancelTx"`

	// Cancelable Is order cancelable
	Cancelable bool `json:"cancelable"`

	// CreatedAt Unix timestamp in milliseconds
	CreatedAt float32 `json:"createdAt"`

	// DstChainId Identifier of the chain where the taker asset is located.
	DstChainId       float32 `json:"dstChainId"`
	DstTokenPriceUsd string  `json:"dstTokenPriceUsd"`

	// Extension An interaction call data. ABI encoded set of makerAssetSuffix, takerAssetSuffix, makingAmountGetter, takingAmountGetter, predicate, permit, preInteraction, postInteraction.If extension exists then lowest 160 bits of the order salt must be equal to the lowest 160 bits of the extension hash
	Extension string `json:"extension"`

	// Fills Fills
	Fills []FillOutputDto `json:"fills"`

	// InitialRateBump Initial rate bump
	InitialRateBump float32                  `json:"initialRateBump"`
	Order           LimitOrderV4StructOutput `json:"order"`

	// OrderHash Order hash
	OrderHash string               `json:"orderHash"`
	Points    []AuctionPointOutput `json:"points"`

	// SrcChainId Identifier of the chain where the maker asset is located.
	SrcChainId       float32 `json:"srcChainId"`
	SrcTokenPriceUsd string  `json:"srcTokenPriceUsd"`

	// Status Order status
	Status GetOrderFillsByHashOutputStatus `json:"status"`

	// TakerAsset Identifier of the asset being requested by the maker in exchange in dst chain.
	TakerAsset string `json:"takerAsset"`

	// TimeLocks TimeLocks without deployedAt
	TimeLocks string `json:"timeLocks"`

	// Validation Order validation status
	Validation GetOrderFillsByHashOutputValidation `json:"validation"`
}

// QuoterControllerGetQuoteParamsFixed defines parameters for QuoterControllerGetQuote.
type QuoterControllerGetQuoteParamsFixed struct {
	// SrcChain Id of source chain
	SrcChain float32 `url:"srcChain" json:"srcChain"`

	// DstChain Id of destination chain
	DstChain float32 `url:"dstChain" json:"dstChain"`

	// SrcTokenAddress Address of "SOURCE" token in source chain
	SrcTokenAddress string `url:"srcTokenAddress" json:"srcTokenAddress"`

	// DstTokenAddress Address of "DESTINATION" token in destination chain
	DstTokenAddress string `url:"dstTokenAddress" json:"dstTokenAddress"`

	// Amount to take from "SOURCE" token to get "DESTINATION" token
	Amount string `url:"amount" json:"amount"`

	// WalletAddress An address of the wallet or contract in source chain who will create Fusion order
	WalletAddress string `url:"walletAddress" json:"walletAddress"`

	// EnableEstimate if enabled then get estimation from 1inch swap builder and generates quoteId, by default is false
	EnableEstimate bool `url:"enableEstimate" json:"enableEstimate"`

	// Fee in bps format, 1% is equal to 100bps
	Fee *big.Int `url:"fee,omitempty" json:"fee,omitempty"` // This is changed from float32 to *big.Int

	// IsPermit2 permit2 allowance transfer encoded call
	IsPermit2 bool `url:"isPermit2,omitempty" json:"isPermit2,omitempty"` // This is changed from string to bool

	// Permit permit, user approval sign
	Permit string `url:"permit,omitempty" json:"permit,omitempty"`
}

// QuoterControllerGetQuoteWithCustomPresetsParamsFixed defines parameters for QuoterControllerGetQuoteWithCustomPresets.
// This is a fixed version with Amount as string instead of float32 for proper BigInt validation.
type QuoterControllerGetQuoteWithCustomPresetsParamsFixed struct {
	// SrcChain Id of source chain
	SrcChain float32 `url:"srcChain" json:"srcChain"`

	// DstChain Id of destination chain
	DstChain float32 `url:"dstChain" json:"dstChain"`

	// SrcTokenAddress Address of "SOURCE" token
	SrcTokenAddress string `url:"srcTokenAddress" json:"srcTokenAddress"`

	// DstTokenAddress Address of "DESTINATION" token
	DstTokenAddress string `url:"dstTokenAddress" json:"dstTokenAddress"`

	// Amount Amount to take from "SOURCE" token to get "DESTINATION" token
	Amount string `url:"amount" json:"amount"`

	// WalletAddress An address of the wallet or contract who will create Fusion order
	WalletAddress string `url:"walletAddress" json:"walletAddress"`

	// EnableEstimate if enabled then get estimation from 1inch swap builder and generates quoteId, by default is false
	EnableEstimate bool `url:"enableEstimate" json:"enableEstimate"`

	// Permit permit, user approval sign
	Permit string `url:"permit,omitempty" json:"permit,omitempty"`
}

// GetQuoteOutputFixed defines model for GetQuoteOutput. QuoteId, DstSafetyDeposit, and SrcSafetyDeposit have been fixed
type GetQuoteOutputFixed struct {
	// DstEscrowFactory Escrow factory contract address at destination chain
	DstEscrowFactory string       `json:"dstEscrowFactory"`
	DstSafetyDeposit string       `json:"dstSafetyDeposit"` // This is changed from string to *big.Int
	DstTokenAmount   string       `json:"dstTokenAmount"`
	Presets          QuotePresets `json:"presets"`
	Prices           PairCurrency `json:"prices"`

	// QuoteId Current generated quote id, should be passed with order
	QuoteId string `json:"quoteId"` // This is changed from map[string]any to string

	// RecommendedPreset suggested preset
	RecommendedPreset GetQuoteOutputRecommendedPreset `json:"recommendedPreset"`

	// SrcEscrowFactory Escrow factory contract address at source chain
	SrcEscrowFactory string       `json:"srcEscrowFactory"`
	SrcSafetyDeposit string       `json:"srcSafetyDeposit"` // This is changed from string to *big.Int
	SrcTokenAmount   string       `json:"srcTokenAmount"`
	TimeLocks        TimeLocks    `json:"timeLocks"`
	Volume           PairCurrency `json:"volume"`

	// Whitelist current executors whitelist addresses
	Whitelist []string `json:"whitelist"`
}

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
	DstChainId       float32
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

// TakingFeeInfo is an alias for the shared fusionorder.TakingFeeInfo type
type TakingFeeInfo = fusionorder.TakingFeeInfo

// CustomPreset is an alias for the shared fusionorder.CustomPreset type
type CustomPreset = fusionorder.CustomPreset

// CustomPresetPoint is an alias for the shared fusionorder.CustomPresetPoint type
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

// PresetClassFixed defines model for PresetClass.
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

// GasCostConfigClass defines model for GasCostConfigClass.
type GasCostConfigClass struct {
	GasBumpEstimate  float32 `json:"gasBumpEstimate"`
	GasPriceEstimate string  `json:"gasPriceEstimate"`
}

// AuctionPointClass defines model for AuctionPointClass.
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
