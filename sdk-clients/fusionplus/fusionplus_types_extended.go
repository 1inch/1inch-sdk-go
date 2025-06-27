package fusionplus

import (
	"math/big"

	"github.com/1inch/1inch-sdk-go/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
)

type GetOrderByOrderHashParams struct {
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
	CancelTx         map[string]interface{} `json:"cancelTx"`

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

// GetQuoteOutputFixed defines model for GetQuoteOutput. QuoteId, DstSafetyDeposit, and SrcSafetyDeposit have been fixed
type GetQuoteOutputFixed struct {
	// DstEscrowFactory Escrow factory contract address at destination chain
	DstEscrowFactory string       `json:"dstEscrowFactory"`
	DstSafetyDeposit string       `json:"dstSafetyDeposit"` // This is changed from string to *big.Int
	DstTokenAmount   string       `json:"dstTokenAmount"`
	Presets          QuotePresets `json:"presets"`
	Prices           PairCurrency `json:"prices"`

	// QuoteId Current generated quote id, should be passed with order
	QuoteId string `json:"quoteId"` // This is changed from map[string]interface{} to string

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
	AuctionDetails      *AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	Extra               ExtraData
}

type EscrowExtensionParams struct {
	fusion.ExtensionParams
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
	Auction            *AuctionDetails        `json:"auction"`
	Whitelist          []AuctionWhitelistItem `json:"whitelist"`
	ResolvingStartTime *big.Int               `json:"resolvingStartTime"`
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
