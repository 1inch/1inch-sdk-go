// Package fusionplus provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package fusionplus

// Defines values for EscrowEventDataOutputAction.
const (
	DstEscrowCreated EscrowEventDataOutputAction = "dst_escrow_created"
	EscrowCancelled  EscrowEventDataOutputAction = "escrow_cancelled"
	FundsRescued     EscrowEventDataOutputAction = "funds_rescued"
	SrcEscrowCreated EscrowEventDataOutputAction = "src_escrow_created"
	Withdrawn        EscrowEventDataOutputAction = "withdrawn"
)

// Defines values for EscrowEventDataOutputSide.
const (
	Dst EscrowEventDataOutputSide = "dst"
	Src EscrowEventDataOutputSide = "src"
)

// Defines values for FillOutputDtoStatus.
const (
	FillOutputDtoStatusExecuted  FillOutputDtoStatus = "executed"
	FillOutputDtoStatusPending   FillOutputDtoStatus = "pending"
	FillOutputDtoStatusRefunded  FillOutputDtoStatus = "refunded"
	FillOutputDtoStatusRefunding FillOutputDtoStatus = "refunding"
)

// Defines values for GetOrderFillsByHashOutputStatus.
const (
	GetOrderFillsByHashOutputStatusCancelled GetOrderFillsByHashOutputStatus = "cancelled"
	GetOrderFillsByHashOutputStatusExecuted  GetOrderFillsByHashOutputStatus = "executed"
	GetOrderFillsByHashOutputStatusExpired   GetOrderFillsByHashOutputStatus = "expired"
	GetOrderFillsByHashOutputStatusPending   GetOrderFillsByHashOutputStatus = "pending"
	GetOrderFillsByHashOutputStatusRefunded  GetOrderFillsByHashOutputStatus = "refunded"
	GetOrderFillsByHashOutputStatusRefunding GetOrderFillsByHashOutputStatus = "refunding"
)

// Defines values for GetOrderFillsByHashOutputValidation.
const (
	FailedToDecodeRemaining            GetOrderFillsByHashOutputValidation = "failed-to-decode-remaining"
	FailedToParsePermitDetails         GetOrderFillsByHashOutputValidation = "failed-to-parse-permit-details"
	InvalidPermitSignature             GetOrderFillsByHashOutputValidation = "invalid-permit-signature"
	InvalidPermitSigner                GetOrderFillsByHashOutputValidation = "invalid-permit-signer"
	InvalidPermitSpender               GetOrderFillsByHashOutputValidation = "invalid-permit-spender"
	InvalidSignature                   GetOrderFillsByHashOutputValidation = "invalid-signature"
	NotEnoughAllowance                 GetOrderFillsByHashOutputValidation = "not-enough-allowance"
	NotEnoughBalance                   GetOrderFillsByHashOutputValidation = "not-enough-balance"
	OrderPredicateReturnedFalse        GetOrderFillsByHashOutputValidation = "order-predicate-returned-false"
	UnknownFailure                     GetOrderFillsByHashOutputValidation = "unknown-failure"
	UnknownPermitVersion               GetOrderFillsByHashOutputValidation = "unknown-permit-version"
	Valid                              GetOrderFillsByHashOutputValidation = "valid"
	WrongEpochManagerAndBitInvalidator GetOrderFillsByHashOutputValidation = "wrong-epoch-manager-and-bit-invalidator"
)

// Defines values for ReadyToExecutePublicActionAction.
const (
	Cancel   ReadyToExecutePublicActionAction = "cancel"
	Withdraw ReadyToExecutePublicActionAction = "withdraw"
)

// Defines values for ResolverDataOutputOrderType.
const (
	MultipleFills ResolverDataOutputOrderType = "MultipleFills"
	SingleFill    ResolverDataOutputOrderType = "SingleFill"
)

// ActiveOrdersOutput defines model for ActiveOrdersOutput.
type ActiveOrdersOutput struct {
	// AuctionEndDate End date of the auction for this order.
	AuctionEndDate float32 `json:"auctionEndDate"`

	// AuctionStartDate Start date of the auction for this order.
	AuctionStartDate float32 `json:"auctionStartDate"`

	// Deadline Deadline by which the order must be filled.
	Deadline float32 `json:"deadline"`

	// DstChainId Identifier of the chain where the taker asset is located.
	DstChainId float32 `json:"dstChainId"`

	// Extension An interaction call data. ABI encoded set of makerAssetSuffix, takerAssetSuffix, makingAmountGetter, takingAmountGetter, predicate, permit, preInteraction, postInteraction.If extension exists then lowest 160 bits of the order salt must be equal to the lowest 160 bits of the extension hash
	Extension string `json:"extension"`

	// Fills Array of fills.
	Fills []string `json:"fills"`

	// IsMakerContract True if order signed by contract (GnosisSafe, etc.)
	IsMakerContract bool `json:"isMakerContract"`

	// MakerAllowance Amount of the maker asset allowance.
	MakerAllowance string `json:"makerAllowance"`

	// MakerBalance Amount of the maker asset balance.
	MakerBalance string             `json:"makerBalance"`
	Order        CrossChainOrderDto `json:"order"`

	// OrderHash Unique identifier of the order.
	OrderHash string `json:"orderHash"`

	// QuoteId Identifier of the quote associated with this order.
	QuoteId string `json:"quoteId"`

	// RemainingMakerAmount Remaining amount of the maker asset that can still be filled.
	RemainingMakerAmount string `json:"remainingMakerAmount"`

	// SecretHashes Array of secret hashes.
	SecretHashes [][]interface{} `json:"secretHashes,omitempty"`

	// Signature Signature of the order.
	Signature string `json:"signature"`

	// SrcChainId Identifier of the chain where the maker asset is located.
	SrcChainId float32 `json:"srcChainId"`
}

// AuctionPointOutput defines model for AuctionPointOutput.
type AuctionPointOutput struct {
	// Coefficient The rate bump from the order min taker amount
	Coefficient float32 `json:"coefficient"`

	// Delay The delay in seconds from the previous point or auction start time
	Delay float32 `json:"delay"`
}

// CrossChainOrderDto defines model for CrossChainOrderDto.
type CrossChainOrderDto struct {
	// Maker Address of the account creating the order (maker) in src chain.
	Maker string `json:"maker"`

	// MakerAsset Identifier of the asset being offered by the maker in src chain.
	MakerAsset string `json:"makerAsset"`

	// MakerTraits Includes some flags like, allow multiple fills, is partial fill allowed or not, price improvement, nonce, deadline etc.
	MakerTraits string `json:"makerTraits"`

	// MakingAmount Amount of the makerAsset being offered by the maker in src chain.
	MakingAmount string `json:"makingAmount"`

	// Receiver Address of the account receiving the assets (receiver), if different from maker in dst chain.
	Receiver string `json:"receiver"`

	// Salt Some unique value. It is necessary to be able to create cross chain orders with the same parameters (so that they have a different hash), Lowest 160 bits of the order salt must be equal to the lowest 160 bits of the extension hash
	Salt string `json:"salt"`

	// TakerAsset Identifier of the asset being requested by the maker in exchange in dst chain.
	TakerAsset string `json:"takerAsset"`

	// TakingAmount Amount of the takerAsset being requested by the maker in dst chain.
	TakingAmount string `json:"takingAmount"`
}

// EscrowEventDataOutput defines model for EscrowEventDataOutput.
type EscrowEventDataOutput struct {
	// Action Action of the escrow event
	Action EscrowEventDataOutputAction `json:"action"`

	// BlockTimestamp Unix timestamp in milliseconds
	BlockTimestamp float32 `json:"blockTimestamp"`

	// Side Side of the escrow event SRC or DST
	Side EscrowEventDataOutputSide `json:"side"`

	// TransactionHash Transaction hash
	TransactionHash string `json:"transactionHash"`
}

// EscrowEventDataOutputAction Action of the escrow event
type EscrowEventDataOutputAction string

// EscrowEventDataOutputSide Side of the escrow event SRC or DST
type EscrowEventDataOutputSide string

// EscrowFactory defines model for EscrowFactory.
type EscrowFactory struct {
	// Address actual escrow factory contract address
	Address string `json:"address"`
}

// FillOutputDto defines model for FillOutputDto.
type FillOutputDto struct {
	EscrowEvents []EscrowEventDataOutput `json:"escrowEvents"`

	// FilledAuctionTakerAmount Amount of the takerAsset filled in dst chain.
	FilledAuctionTakerAmount string `json:"filledAuctionTakerAmount"`

	// FilledMakerAmount Amount of the makerAsset filled in src chain.
	FilledMakerAmount string `json:"filledMakerAmount"`

	// Status Fill status
	Status FillOutputDtoStatus `json:"status"`

	// TxHash Transaction hash
	TxHash string `json:"txHash"`
}

// FillOutputDtoStatus Fill status
type FillOutputDtoStatus string

// GetActiveOrdersOutput defines model for GetActiveOrdersOutput.
type GetActiveOrdersOutput struct {
	Items []ActiveOrdersOutput `json:"items"`
	Meta  Meta                 `json:"meta"`
}

// GetOrderByMakerOutput defines model for GetOrderByMakerOutput.
type GetOrderByMakerOutput struct {
	Items []ActiveOrdersOutput `json:"items"`
	Meta  Meta                 `json:"meta"`
}

// GetOrderFillsByHashOutput defines model for GetOrderFillsByHashOutput.
type GetOrderFillsByHashOutput struct {
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
	DstChainId       float32                `json:"dstChainId"`
	DstTokenPriceUsd map[string]interface{} `json:"dstTokenPriceUsd"`

	// Extension An interaction call data. ABI encoded set of makerAssetSuffix, takerAssetSuffix, makingAmountGetter, takingAmountGetter, predicate, permit, preInteraction, postInteraction.If extension exists then lowest 160 bits of the order salt must be equal to the lowest 160 bits of the extension hash
	Extension string `json:"extension"`

	// Fills Fills
	Fills []FillOutputDto `json:"fills"`

	// InitialRateBump Initial rate bump
	InitialRateBump float32                  `json:"initialRateBump"`
	Order           LimitOrderV4StructOutput `json:"order"`

	// OrderHash Order hash
	OrderHash string             `json:"orderHash"`
	Points    AuctionPointOutput `json:"points"`

	// SrcChainId Identifier of the chain where the maker asset is located.
	SrcChainId       float32                `json:"srcChainId"`
	SrcTokenPriceUsd map[string]interface{} `json:"srcTokenPriceUsd"`

	// Status Order status
	Status GetOrderFillsByHashOutputStatus `json:"status"`

	// TakerAsset Identifier of the asset being requested by the maker in exchange in dst chain.
	TakerAsset string `json:"takerAsset"`

	// TimeLocks TimeLocks without deployedAt
	TimeLocks string `json:"timeLocks"`

	// Validation Order validation status
	Validation GetOrderFillsByHashOutputValidation `json:"validation"`
}

// GetOrderFillsByHashOutputStatus Order status
type GetOrderFillsByHashOutputStatus string

// GetOrderFillsByHashOutputValidation Order validation status
type GetOrderFillsByHashOutputValidation string

// Immutables defines model for Immutables.
type Immutables struct {
	// Amount Amount of token to receive
	Amount string `json:"amount"`

	// Hashlock keccak256(secret(idx))
	Hashlock string `json:"hashlock"`

	// Maker Maker's address which will receive tokens
	Maker string `json:"maker"`

	// OrderHash Order's hash 32 bytes hex sting
	OrderHash string `json:"orderHash"`

	// SafetyDeposit Security deposit in chain's native currency
	SafetyDeposit string `json:"safetyDeposit"`

	// Taker Escrow creation initiator address
	Taker string `json:"taker"`

	// Timelocks Encoded timelocks. To decode use: https://github.com/1inch/cross-chain-sdk/blob/master/src/cross-chain-order/time-locks/time-locks.ts
	Timelocks string `json:"timelocks"`

	// Token Token to receive on specific chain
	Token string `json:"token"`
}

// LimitOrderV4StructOutput defines model for LimitOrderV4StructOutput.
type LimitOrderV4StructOutput struct {
	// Maker Maker address
	Maker string `json:"maker"`

	// MakerAsset Maker asset address
	MakerAsset  string `json:"makerAsset"`
	MakerTraits string `json:"makerTraits"`

	// MakingAmount Amount of the maker asset
	MakingAmount string `json:"makingAmount"`

	// Receiver Receiver address
	Receiver string `json:"receiver"`
	Salt     string `json:"salt"`

	// TakerAsset Taker asset address
	TakerAsset string `json:"takerAsset"`

	// TakingAmount Amount of the taker asset
	TakingAmount string `json:"takingAmount"`
}

// Meta defines model for Meta.
type Meta struct {
	CurrentPage  float32 `json:"currentPage"`
	ItemsPerPage float32 `json:"itemsPerPage"`
	TotalItems   float32 `json:"totalItems"`
	TotalPages   float32 `json:"totalPages"`
}

// OrdersByHashesInput defines model for OrdersByHashesInput.
type OrdersByHashesInput struct {
	OrderHashes []string `json:"orderHashes"`
}

// PublicSecret defines model for PublicSecret.
type PublicSecret struct {
	DstImmutables Immutables `json:"dstImmutables"`

	// Idx Sequence number of secrets
	Idx float32 `json:"idx"`

	// Secret Public secret to perform a withdrawal
	Secret        string     `json:"secret"`
	SrcImmutables Immutables `json:"srcImmutables"`
}

// ReadyToAcceptSecretFill defines model for ReadyToAcceptSecretFill.
type ReadyToAcceptSecretFill struct {
	// DstEscrowDeployTxHash Transaction hash where the destination chain escrow was deployed
	DstEscrowDeployTxHash string `json:"dstEscrowDeployTxHash"`

	// Idx Sequence number of secrets for submission
	Idx float32 `json:"idx"`

	// SrcEscrowDeployTxHash Transaction hash where the source chain escrow was deployed
	SrcEscrowDeployTxHash string `json:"srcEscrowDeployTxHash"`
}

// ReadyToAcceptSecretFills defines model for ReadyToAcceptSecretFills.
type ReadyToAcceptSecretFills struct {
	// Fills Fills that are ready to accept secrets from the client
	Fills []ReadyToAcceptSecretFill `json:"fills"`
}

// ReadyToAcceptSecretFillsForAllOrders defines model for ReadyToAcceptSecretFillsForAllOrders.
type ReadyToAcceptSecretFillsForAllOrders struct {
	// Orders Fills that are ready to accept secrets from the client for all orders
	Orders []ReadyToAcceptSecretFillsForOrder `json:"orders"`
}

// ReadyToAcceptSecretFillsForOrder defines model for ReadyToAcceptSecretFillsForOrder.
type ReadyToAcceptSecretFillsForOrder struct {
	// Fills Fills that are ready to accept secrets from the client
	Fills []ReadyToAcceptSecretFill `json:"fills"`

	// MakerAddress Maker address
	MakerAddress string `json:"makerAddress"`

	// OrderHash Order hash
	OrderHash string `json:"orderHash"`
}

// ReadyToExecutePublicAction defines model for ReadyToExecutePublicAction.
type ReadyToExecutePublicAction struct {
	Action ReadyToExecutePublicActionAction `json:"action"`

	// ChainId Execute action on this chain
	ChainId float32 `json:"chainId"`

	// Escrow Escrow's address to perform public action
	Escrow     string     `json:"escrow"`
	Immutables Immutables `json:"immutables"`

	// Secret Presented only for withdraw action
	Secret string `json:"secret,omitempty"`
}

// ReadyToExecutePublicActionAction defines model for ReadyToExecutePublicAction.Action.
type ReadyToExecutePublicActionAction string

// ReadyToExecutePublicActionsOutput defines model for ReadyToExecutePublicActionsOutput.
type ReadyToExecutePublicActionsOutput struct {
	// Actions Actions allowed to be performed on public timelock periods
	Actions []ReadyToExecutePublicAction `json:"actions"`
}

// ResolverDataOutput defines model for ResolverDataOutput.
type ResolverDataOutput struct {
	// OrderType Type of the order: enabled or disabled partial fills
	OrderType ResolverDataOutputOrderType `json:"orderType"`

	// SecretHashes keccak256(secret(idx))[]
	SecretHashes [][]interface{} `json:"secretHashes,omitempty"`

	// Secrets The data required for order withdraw and cancel
	Secrets []PublicSecret `json:"secrets"`
}

// ResolverDataOutputOrderType Type of the order: enabled or disabled partial fills
type ResolverDataOutputOrderType string

// OrderApiControllerGetActiveOrdersParams defines parameters for OrderApiControllerGetActiveOrders.
type OrderApiControllerGetActiveOrdersParams struct {
	// Page Pagination step, default: 1 (page = offset / limit)
	Page float32 `url:"page,omitempty" json:"page,omitempty"`

	// Limit Number of active orders to receive (default: 100, max: 500)
	Limit float32 `url:"limit,omitempty" json:"limit,omitempty"`

	// SrcChain Source chain of cross chain
	SrcChain float32 `url:"srcChain,omitempty" json:"srcChain,omitempty"`

	// DstChain Destination chain of cross chain
	DstChain float32 `url:"dstChain,omitempty" json:"dstChain,omitempty"`
}

// OrderApiControllerGetSettlementContractParams defines parameters for OrderApiControllerGetSettlementContract.
type OrderApiControllerGetSettlementContractParams struct {
	// ChainId Chain ID
	ChainId float32 `url:"chainId,omitempty" json:"chainId,omitempty"`
}

// OrderApiControllerGetOrdersByMakerParams defines parameters for OrderApiControllerGetOrdersByMaker.
type OrderApiControllerGetOrdersByMakerParams struct {
	// Page Pagination step, default: 1 (page = offset / limit)
	Page float32 `url:"page,omitempty" json:"page,omitempty"`

	// Limit Number of active orders to receive (default: 100, max: 500)
	Limit float32 `url:"limit,omitempty" json:"limit,omitempty"`

	// TimestampFrom timestampFrom in milliseconds for interval [timestampFrom, timestampTo)
	TimestampFrom float32 `url:"timestampFrom,omitempty" json:"timestampFrom,omitempty"`

	// TimestampTo timestampTo in milliseconds for interval [timestampFrom, timestampTo)
	TimestampTo float32 `url:"timestampTo,omitempty" json:"timestampTo,omitempty"`

	// SrcToken Find history by the given source token
	SrcToken string `url:"srcToken,omitempty" json:"srcToken,omitempty"`

	// DstToken Find history by the given destination token
	DstToken string `url:"dstToken,omitempty" json:"dstToken,omitempty"`

	// WithToken Find history items by source or destination token
	WithToken string `url:"withToken,omitempty" json:"withToken,omitempty"`

	// DstChainId Destination chain of cross chain
	DstChainId float32 `url:"dstChainId,omitempty" json:"dstChainId,omitempty"`

	// SrcChainId Source chain of cross chain
	SrcChainId float32 `url:"srcChainId,omitempty" json:"srcChainId,omitempty"`

	// ChainId chainId for looking by dstChainId == chainId OR srcChainId == chainId
	ChainId float32 `url:"chainId,omitempty" json:"chainId,omitempty"`
}

// OrderApiControllerGetOrdersByOrderHashesJSONRequestBody defines body for OrderApiControllerGetOrdersByOrderHashes for application/json ContentType.
type OrderApiControllerGetOrdersByOrderHashesJSONRequestBody = OrdersByHashesInput