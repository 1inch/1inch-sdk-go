# CLAUDE.md - 1inch SDK Go

## Overview

This is the official Go SDK for interacting with 1inch Network APIs. It provides type-safe Go bindings for all 1inch DEX aggregation, limit orders, Fusion swaps, and various Web3 data services.

**Module**: `github.com/1inch/1inch-sdk-go`  
**Go Version**: 1.22+ (toolchain 1.23.0)  
**License**: MIT  
**Current Version**: v2.0.0

## Project Structure

```
/
├── sdk-clients/          # Main SDK client implementations (one per API)
│   ├── aggregation/      # DEX aggregation (swap) API client
│   ├── fusion/           # Fusion swap (gasless) API client
│   ├── fusionplus/       # Fusion+ cross-chain swap client
│   ├── orderbook/        # Limit order protocol client
│   ├── balances/         # Token balance/allowance queries
│   ├── gasprices/        # Gas price oracle
│   ├── history/          # Transaction history
│   ├── nft/              # NFT data queries
│   ├── portfolio/        # Portfolio tracking
│   ├── spotprices/       # Token spot prices
│   ├── tokens/           # Token information
│   ├── traces/           # Transaction traces
│   ├── txbroadcast/      # Transaction broadcasting
│   └── web3/             # Web3 RPC calls
├── common/               # Shared interfaces and types
│   └── fusionorder/      # Shared types/utilities for fusion and fusionplus
├── constants/            # Chain IDs, contract addresses, ABIs
├── internal/             # Internal utilities (not exported)
│   ├── http-executor/    # HTTP client implementation
│   ├── web3-provider/    # Ethereum wallet/provider implementation
│   ├── transaction-builder/  # Tx construction utilities
│   ├── validate/         # Parameter validation functions
│   ├── bigint/           # Big integer utilities
│   ├── keccak/           # Keccak hashing
│   └── multicall/        # Multicall contract support
├── codegen/              # OpenAPI spec files and type generation
│   ├── openapi/          # OpenAPI JSON specs for each API
│   ├── generate_types.sh # Type generation script
│   └── mapping.json      # Operation ID mappings
└── .github/workflows/    # CI/CD (pr.yml, release.yml)
```

## Key Commands

```bash
# Run all unit tests
make test

# Run linter (golangci-lint v1.54.1)
make lint

# Format code
make fmt

# Generate types from OpenAPI specs
make codegen-types   # Must run from codegen/ directory

# Get dependencies
make get
```

## Post-Change Verification

After making changes, run these checks to ensure code quality:

```bash
# 1. Build all packages (catch compile errors)
go build ./...

# 2. Run static analysis (catch common issues)
go vet ./...

# 3. Run linter with import formatting check
golangci-lint run --enable goimports

# 4. Run all tests
go test ./...

# Or run tests with race detection (slower but catches race conditions)
go test -race ./...
```

To auto-fix import formatting issues:
```bash
goimports -w .
```

## Architecture Patterns

### Client Pattern
Each SDK client follows a consistent pattern:

```go
// 1. Create configuration
config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
    NodeUrl:    nodeUrl,       // Ethereum RPC endpoint
    PrivateKey: privateKey,    // Wallet private key (64 hex chars, no 0x)
    ChainId:    constants.EthereumChainId,
    ApiUrl:     "https://api.1inch.dev",
    ApiKey:     devPortalToken,  // 1inch Dev Portal API key
})

// 2. Create client
client, err := aggregation.NewClient(config)

// 3. Use API methods
quote, err := client.GetQuote(ctx, aggregation.GetQuoteParams{...})
```

### Two Client Variants
- **Full Client** (`NewClient`): Includes wallet + API access for on-chain operations
- **API-Only Client** (`NewClientOnlyAPI`): API access without wallet (read-only operations)

### Core Interfaces (in `common/`)

```go
// HttpExecutor - HTTP request execution
type HttpExecutor interface {
    ExecuteRequest(ctx context.Context, payload RequestPayload, v interface{}) error
}

// Wallet - Ethereum wallet operations
type Wallet interface {
    Call(ctx context.Context, contractAddress Address, callData []byte) ([]byte, error)
    Nonce(ctx context.Context) (uint64, error)
    Address() Address
    Sign(tx *types.Transaction) (*types.Transaction, error)
    SignBytes(data []byte) ([]byte, error)
    BroadcastTransaction(ctx context.Context, tx *types.Transaction) error
    // ... permit methods, gas methods
}

// TransactionBuilder - Fluent transaction construction
type TransactionBuilder interface {
    SetData([]byte) TransactionBuilder
    SetNonce(uint64) TransactionBuilder
    SetGas(uint64) TransactionBuilder
    SetTo(*Address) TransactionBuilder
    SetValue(*big.Int) TransactionBuilder
    Build(context.Context) (*types.Transaction, error)
}
```

### Validation Pattern
Each API client has a `validation.go` with `Validate()` methods on params:

```go
func (params *GetSwapParams) Validate() error {
    var validationErrors []error
    validationErrors = validate.Parameter(params.Src, "src", validate.CheckEthereumAddressRequired, validationErrors)
    // ... more validations
    return validate.ConsolidateValidationErorrs(validationErrors)
}
```

Common validators in `internal/validate/`:
- `CheckEthereumAddressRequired` / `CheckEthereumAddress`
- `CheckBigIntRequired` / `CheckBigInt`
- `CheckSlippageRequired` (0.01-50 range)
- `CheckPrivateKeyRequired` (64 hex chars)
- `CheckChainIdIntRequired`

## Type Generation

Types are auto-generated from OpenAPI specs using `oapi-codegen`:

1. OpenAPI specs live in `codegen/openapi/*-openapi.json`
2. Run `./generate_types.sh` from `codegen/` directory
3. Generated files: `sdk-clients/{package}/{package}_types.gen.go`

**DO NOT manually edit `*_types.gen.go` files** - they are overwritten by codegen.

Extended types (manual additions) go in `*_types_extended.go` files.

## Supported Chains

Defined in `constants/chains.go`:
```go
EthereumChainId  = 1
PolygonChainId   = 137
BscChainId       = 56
ArbitrumChainId  = 42161
OptimismChainId  = 10
AvalancheChainId = 43114
GnosisChainId    = 100
FantomChainId    = 250
BaseChainId      = 8453
ZkSyncEraChainId = 324  // Limited support
AuroraChainId    = 1313161554
KlaytnChainId    = 8217
```

## Contract Addresses

Key addresses in `constants/contracts.go`:
- `AggregationRouterV6` = `0x111111125421cA6dc452d289314280a0f8842A65` (all chains except zkSync)
- `NativeToken` = `0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee`
- SeriesNonceManager addresses per chain

## Environment Variables for Examples

```bash
WALLET_KEY=<64-char-hex-private-key>  # No 0x prefix
WALLET_ADDRESS=<your-address>
NODE_URL=<ethereum-rpc-url>
DEV_PORTAL_TOKEN=<1inch-api-key>
```

## Testing

- Unit tests use `github.com/stretchr/testify`
- Tests are in `*_test.go` files alongside source
- Run with `make test` or `go test -race ./...`

## CI/CD

- **PR Validation** (`.github/workflows/pr.yml`): Runs tests + golangci-lint on PRs
- **Release** (`.github/workflows/release.yml`): Manual dispatch for versioned releases

## Key Dependencies

- `github.com/ethereum/go-ethereum v1.14.13` - Ethereum client
- `github.com/google/go-querystring v1.1.0` - URL query encoding
- `github.com/oapi-codegen/runtime v1.1.1` - OpenAPI runtime
- `github.com/stretchr/testify v1.9.0` - Testing

## API Clients Summary

| Client | Purpose | Key Methods |
|--------|---------|-------------|
| `aggregation` | DEX swap aggregation | `GetQuote`, `GetSwap`, `GetApproveAllowance` |
| `fusion` | Gasless swaps | `GetQuote`, `PlaceOrder`, `GetOrderStatus` |
| `fusionplus` | Cross-chain gasless | `GetQuote`, `PlaceOrder` |
| `orderbook` | Limit orders | `CreateOrder`, `GetAllOrders` |
| `balances` | Token balances | `GetBalancesAndAllowances` |
| `gasprices` | Gas oracle | `GetGasPrices` |
| `spotprices` | Token prices | `GetPricesForRequestedTokens` |
| `tokens` | Token info | `GetWhitelistedTokens`, `SearchToken` |
| `portfolio` | Portfolio tracking | `GetCurrentValue`, `GetProfitLoss` |
| `history` | Tx history | `GetHistoryEventsByAddress` |
| `web3` | RPC proxy | `PerformRpcCall` |
| `txbroadcast` | Tx broadcasting | `BroadcastPublicTransaction` |
| `traces` | Tx traces | `GetTxTrace` |
| `nft` | NFT data | `GetNftsByAddress` |

## Common Patterns

### Executing a Swap
```go
// 1. Get swap data
swapData, _ := client.GetSwap(ctx, aggregation.GetSwapParams{
    Src: "0x...", Dst: "0x...",
    Amount: "1000000", From: wallet.Address().Hex(),
    Slippage: 1,
})

// 2. Build transaction
tx, _ := client.TxBuilder.New().
    SetData(swapData.TxNormalized.Data).
    SetTo(&swapData.TxNormalized.To).
    SetGas(swapData.TxNormalized.Gas).
    SetValue(swapData.TxNormalized.Value).
    Build(ctx)

// 3. Sign and broadcast
signedTx, _ := client.Wallet.Sign(tx)
client.Wallet.BroadcastTransaction(ctx, signedTx)
```

### Fusion Order
```go
// 1. Get quote
quote, _ := client.GetQuote(ctx, fusion.QuoterControllerGetQuoteParamsFixed{...})

// 2. Place order (signs and submits)
orderHash, _ := client.PlaceOrder(ctx, *quote, fusion.OrderParams{
    Preset: fusion.Fast,
    ...
}, client.Wallet)

// 3. Monitor status
status, _ := client.GetOrderStatus(ctx, orderHash)
```

## File Naming Conventions

- `client.go` - Client struct and constructor
- `api.go` - API method implementations
- `configuration.go` - Configuration structs/constructors
- `validation.go` - Parameter validation
- `*_types.gen.go` - Auto-generated types (DO NOT EDIT)
- `*_types_extended.go` - Manual type extensions
- `examples/` - Usage examples per operation

## Fusion Package Architecture

The `fusion`, `fusionplus`, and `fusionorder` packages share a layered architecture:

### Package Hierarchy
```
common/fusionorder/   # Shared types and utilities (base layer)
├── bps.go            # Basis points type and operations
├── interaction.go    # Order interaction encoding/decoding
├── whitelist.go      # Whitelist generation (GenerateWhitelist, GenerateWhitelistFromItems)
├── whitelist_utils.go # Whitelist helpers (CanExecuteAt, IsExclusiveResolver)
├── auction.go        # Auction details encoding
├── nativetokenwrappers.go # Chain-specific wrapped token addresses
└── ...

fusion/               # Fusion (single-chain gasless swaps)
├── Uses fusionorder types via aliases
├── Extension, ExtensionParams (fusion-specific)
└── SettlementPostInteractionData (fusion-specific encoding)

fusionplus/           # Fusion+ (cross-chain swaps)
├── Uses fusionorder types via aliases
├── ExtensionPlus, ExtensionParamsPlus (fusionplus-specific)
├── EscrowExtension (cross-chain escrow data)
└── SettlementPostInteractionData (fusionplus-specific encoding)
```

### Naming Conventions

**Type Aliases**: Both `fusion` and `fusionplus` use type aliases to re-export shared types:
```go
// In fusion/fusion_types_extended.go
type Bps = fusionorder.Bps
type Interaction = fusionorder.Interaction
type WhitelistItem = fusionorder.WhitelistItem

// In fusionplus/fusionplus_types_extended.go  
type Interaction = fusionorder.Interaction
type WhitelistItem = fusionorder.WhitelistItem
```

**Plus Suffix**: Types in `fusionplus` that need to be distinguished from `fusion` equivalents use the `Plus` suffix:
- `ExtensionPlus` (not `Extension` - would conflict with fusion.Extension conceptually)
- `ExtensionParamsPlus`
- `NewExtensionPlus()`
- `CreateAuctionDetailsPlus()`

**Variable Naming**: Match the package context:
- In `fusion`: `fusionExtension`, `fusionOrder`
- In `fusionplus`: `extensionPlus`, `auctionDetailsPlus`, `presetPlus`

### Function Aliasing Pattern

When a function is shared, create an alias rather than duplicating code:
```go
// In fusion/settlementpostinteractiondata.go
var GenerateWhitelist = fusionorder.GenerateWhitelist
```

### Encoding Differences

The packages have different binary encoding formats:
- **Fusion**: Auction details include a point count byte (`Encode()`)
- **FusionPlus**: Auction details omit the point count byte (`EncodeWithoutPointCount()`)

### Common Pitfalls to Avoid

1. **Don't duplicate functions** - If logic is identical, put it in `fusionorder` and alias it
2. **Don't use `log.Fatalf`** - Always return errors properly
3. **Watch for case sensitivity** - Ethereum addresses are case-insensitive; use `strings.ToLower()` when comparing address halves
4. **Keep variable names consistent** - Use `*Plus` in fusionplus, not `*Fusion`

## Breaking Changes Documentation

When making breaking changes, update both files:
- `BREAKING_CHANGES.md` - Detailed migration guide with tables
- `CHANGELOG.md` - Summary for release notes

## Notes for Development

1. **ABIs are embedded** via `//go:embed` in `constants/abis.go`
2. **BigInt amounts** are passed as strings to avoid overflow
3. **Transaction builder** auto-fetches nonce/gas if not set
4. **EIP-1559** support is automatic based on chain
5. **Permit1** support for gasless approvals on supported tokens
6. **API rate limits** - Use a valid Dev Portal API key
7. **Error handling** - API errors are JSON with `statusCode`, `error`, `description`
8. **No `log.Fatalf`** - Always return errors; `log.Fatalf` terminates the program
9. **Check for duplicate code** - Before adding a function, search if it exists elsewhere
10. **Test files should match source** - If `foo.go` is deleted, check if `foo_test.go` should be too
11. **Always run post-change verification** - See "Post-Change Verification" section above
