# Breaking Changes

This document tracks breaking changes between major versions of the SDK that affect users importing and integrating the library.

For a step-by-step upgrade checklist, see [MIGRATION.md](MIGRATION.md); this document is the per-change reference it links back to.

## Unreleased

Step-by-step upgrade instructions for everything below: [MIGRATION.md](MIGRATION.md).

### Compile-Time Breaking Changes

- **Chain-id fields are `int` instead of `float32`** across `fusionplus` and `tokens` (quoter/orders params, response DTOs, `EscrowExtension`/`EscrowExtensionParams`/`EscrowExtraData`, `SignedOrderInput`, `GetSettlementContractParams`). `float32` cannot represent Aurora's chain id (1313161554), which silently rounded and corrupted the EIP-712 signing domain. Untyped constants and literals compile unchanged; explicit `float32(...)` conversions do not.
- **Corrected field types on the generated quoter/orders types** (previously exclusive to the `*Fixed` copies): `Amount` `float32` → `string` on all fusion/fusionplus quoter params; `fusionplus` `IsPermit2` `string` → `bool`; `fusion.Quote.QuoteId` and `fusionplus.Quote.QuoteId` `map[string]interface{}` → `string`; `fusion.PresetClass.ExclusiveResolver` `map[string]interface{}` → `string`; `fusionplus.GetOrderFillsByHashOutput` USD prices `map[string]interface{}` → `string` and `Points` → `[]AuctionPointOutput`.
- **`Fee` on the fusionplus quoter params is `int` (basis points)** — previously `*big.Int` on the `Fixed` types and `float32` on the raw generated types.
- **`tokens` optional fields are values instead of pointers** (`DisplayedSymbol`, `Eip2612`, `IsFoT`, `LogoURI`, `Extensions`), and `Tags` is `[]TagDto` (was `[]string` on the generated types).
- **Optional object-typed response fields lost their pointers**: `history.TransactionDetailsDto.Meta`, `nft.Asset.Collection`, `nft.Asset.RarityData`.
- **Bare enum constants in `spotprices` and `traces` are prefixed with their type name** (`spotprices.USD` → `spotprices.GetPricesForRequestedTokensParamsCurrencyUSD`, `traces.CALL` → `traces.CoreCustomRootTxEventCallstackTraceFullDtoTypeCALL`).
- **`web3.ApiKeyAuthScopes` removed** (unused generated artifact).
- **`fusionplus.CreateAuctionDetailsPlus` takes `*Preset`** instead of `*PresetClassFixed`.
- **Embedded field names changed in four `orderbook` wrapper params** (`GetAllOrdersParams`, `GetCountParams`, `GetEventsParams`, `GetOrdersByCreatorAddressParams`): struct literals naming the embedded query field must use the new name (e.g. `LimitOrdersQueryParams`). Promoted field access (`params.Page`) is unaffected.

Not breaking, for the avoidance of doubt: ~60 generated types and constants were renamed to intent-based names (`fusionplus.QuoteParams`, `fusionplus.Quote`, …), `fusionorder.AuctionPointClassFixed`/`GasCostConfigClassFixed` became `AuctionPoint`/`GasCostConfig`, and all 12 `*Fixed` types became aliases — every old name remains available as a `// Deprecated:` alias of the identical type, so existing code compiles and behaves unchanged.

### Behavior Changes (compile clean, act differently)

- **`fusionplus` quote fees are now transmitted.** Previously a configured `Fee` was silently omitted from the quote request (`*big.Int` cannot be encoded by go-querystring), so quotes were priced without the fee the order then carried.
- **`fusionplus.DecodeEscrowExtension` output is corrected**: source/destination safety deposits are no longer swapped, deposits are decimal strings, and the hashlock is 0x-prefixed hex — decode → re-encode is now lossless. Malformed input returns an error instead of panicking.
- **`EscrowExtension.ConvertToOrderbookExtension` no longer mutates the receiver** (a second call previously double-appended the escrow extra data, corrupting the extension).
- **API requests time out after 60 seconds** (previously they could hang forever) and response bodies are capped at 64 MiB.
- `fusion.Quote.SurplusFee` and `fusion.Quote.MarketAmount` marshal with `omitempty`; fee params are validated as basis points in [0, 10000].

## Version 4.0.0

### Module Path Now Includes the `/v4` Major-Version Suffix

Per [Go module rules](https://go.dev/ref/mod#major-version-suffixes), modules at major version 2 or higher must include the major-version suffix in their module path. The module path is now `github.com/1inch/1inch-sdk-go/v4` (previously `github.com/1inch/1inch-sdk-go`). Without this suffix the Go toolchain rejects `v4.x.x` tags as invalid and the module proxy cannot resolve them.

**Impact:** All imports must include the `/v4` suffix.

**Migration:** Update your import paths and `go get` commands:

```go
// Before
import "github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"

// After
import "github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation"
```

```bash
go get github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation
```

### Minimum Go Version Raised to 1.25

The module's `go` directive has been bumped from `go 1.22` to `go 1.25.0`, and the explicit `toolchain go1.23.0` line has been removed. This change is required by the upgrade to `golang.org/x/crypto` v0.52.0 (part of a batch of security dependency upgrades that clear open Dependabot advisories).

**Impact:** Downstream projects must build with Go 1.25 or newer. Projects pinned to an older Go toolchain will fail to compile against this version of the SDK.

**Migration:** Update your toolchain to Go 1.25+ (e.g. bump the `go` directive in your own `go.mod` and your CI Go version).

## Version 3.0.0

### New Shared `fusionorder` Package

Common types and functions previously in `fusion` and `fusionplus` have been consolidated into a new shared package at `common/fusionorder/`. This is the single source of truth for order-related types used by both packages.

#### Moved Exports from `fusion`

| Old Import | New Import |
|------------|------------|
| `fusion.BpsZero` | `fusionorder.BpsZero` |
| `fusion.NewBps()` | `fusionorder.NewBps()` (signature changed: now returns `(*Bps, error)`) |
| `fusion.FromPercent()` | `fusionorder.FromPercent()` (now returns `(*Bps, error)`) |
| `fusion.FromFraction()` | `fusionorder.FromFraction()` (now returns `(*Bps, error)`) |
| `fusion.GetDefaultBase()` | `fusionorder.GetDefaultBase()` |
| `fusion.NewInteraction()` | `fusionorder.NewInteraction()` |
| `fusion.DecodeInteraction()` | `fusionorder.DecodeInteraction()` |
| `fusion.NewAuctionDetails()` | `fusionorder.NewAuctionDetails()` |
| `fusion.DecodeAuctionDetails()` | `fusionorder.DecodeLegacyAuctionDetails()` |
| `fusion.CalcAuctionStartTime()` | `fusionorder.CalcAuctionStartTime()` |
| `fusion.CalcAuctionStartTimeFunc` | `fusionorder.CalcAuctionStartTimeFunc` |
| `fusion.GenerateWhitelist()` | `fusionorder.GenerateWhitelist()` |
| `fusion.BpsToRatioFormat()` | `fusionorder.BpsToRatioFormat()` |
| `fusion.NativeToken` | `constants.NativeToken` |
| `fusion.NetworkEnum` | `constants.NetworkEnum` |
| `fusion.ETHEREUM`, `fusion.POLYGON`, etc. | `constants.NetworkEthereum`, `constants.NetworkPolygon`, etc. |
| `fusion.Bps` | `fusionorder.Bps` |
| `fusion.Interaction` | `fusionorder.Interaction` |
| `fusion.AuctionDetails` | `fusionorder.AuctionDetails` |
| `fusion.WhitelistItem` | `fusionorder.WhitelistItem` |
| `fusion.AuctionWhitelistItem` | `fusionorder.AuctionWhitelistItem` |
| `fusion.ExtraData` | `fusionorder.ExtraData` |

#### Moved Exports from `fusionplus`

| Old Import | New Import |
|------------|------------|
| `fusionplus.NewInteraction()` | `fusionorder.NewInteraction()` |
| `fusionplus.DecodeInteraction()` | `fusionorder.DecodeInteraction()` |
| `fusionplus.NewAuctionDetails()` | `fusionorder.NewAuctionDetails()` |
| `fusionplus.DecodeAuctionDetails()` | `fusionorder.DecodeAuctionDetails()` |
| `fusionplus.CalcAuctionStartTime()` | `fusionorder.CalcAuctionStartTime()` |
| `fusionplus.CalcAuctionStartTimeFunc` | `fusionorder.CalcAuctionStartTimeFunc` |
| `fusionplus.CreateMakerTraitsFusion()` | `fusionplus.CreateMakerTraits()` (param types changed: `Details` replaces `DetailsFusion`) |
| `fusionplus.NativeToken` | `constants.NativeToken` |
| `fusionplus.NetworkEnum` | `constants.NetworkEnum` |
| `fusionplus.ETHEREUM`, `fusionplus.POLYGON`, etc. | `constants.NetworkEthereum`, `constants.NetworkPolygon`, etc. |
| `fusionplus.Interaction` | `fusionorder.Interaction` |
| `fusionplus.AuctionDetails` | `fusionorder.AuctionDetails` |
| `fusionplus.WhitelistItem` | `fusionorder.WhitelistItem` |
| `fusionplus.AuctionWhitelistItem` | `fusionorder.AuctionWhitelistItem` |
| `fusionplus.ExtraData` | `fusionorder.ExtraData` |

#### Migration Example

Before:
```go
import "github.com/1inch/1inch-sdk-go/sdk-clients/fusion"

bps := fusion.NewBps(big.NewInt(100))          // v2: returned *Bps (no error)
details := fusion.NewAuctionDetails(...)
```

After:
```go
import (
    "github.com/1inch/1inch-sdk-go/common/fusionorder"
    "github.com/1inch/1inch-sdk-go/constants"
)

bps, err := fusionorder.NewBps(big.NewInt(100)) // v3: now returns (*Bps, error)
details := fusionorder.NewAuctionDetails(...)
token := constants.NetworkEthereum              // was fusion.ETHEREUM
```

#### Deprecated Type Aliases

The following types are preserved as **deprecated type aliases** for backward compatibility. IDEs with `gopls` support will show these with strikethrough. Migrate to importing from `fusionorder` directly:

**In `fusion`:**
- `fusion.TakingFeeInfo` → use `fusionorder.TakingFeeInfo`
- `fusion.CustomPreset` → use `fusionorder.CustomPreset`
- `fusion.CustomPresetPoint` → use `fusionorder.CustomPresetPoint`

**In `fusionplus`:**
- `fusionplus.TakingFeeInfo` → use `fusionorder.TakingFeeInfo`
- `fusionplus.CustomPreset` → use `fusionorder.CustomPreset`
- `fusionplus.CustomPresetPoint` → use `fusionorder.CustomPresetPoint`

### Renamed Types and Functions in `fusionplus`

Types with the `Fusion` suffix have been renamed. Types that need to be distinguished from `fusion` equivalents now use the `Plus` suffix:

| Old Name | New Name |
|----------|----------|
| `ExtensionFusion` | `ExtensionPlus` |
| `ExtensionParamsFusion` | `ExtensionParamsPlus` |
| `NewExtensionFusion()` | `NewExtensionPlus()` |
| `CreateAuctionDetailsFusion()` | `CreateAuctionDetailsPlus()` |
| `PresetClassFixedFusion` | `PresetClassFixed` |
| `GasCostConfigClassFusion` | `GasCostConfigClass` |
| `AuctionPointClassFusion` | `AuctionPointClass` |
| `SettlementSuffixDataFusion` | Removed (use `SettlementSuffixData`) |
| `FeesFusion` | Removed (use `Fees`) |
| `IntegratorFeeFusion` | Removed (use `IntegratorFee`) |
| `DetailsFusion` | Removed (use `Details`) |
| `SettlementPostInteractionDataFusion` | Merged into `SettlementPostInteractionData` |
| `NewSettlementPostInteractionDataFusion()` | `NewSettlementPostInteractionDataWithFees()` |
| `CreateSettlementPostInteractionDataFusion()` | `CreateSettlementPostInteractionDataWithFees()` |
| `DecodeFusion()` | `DecodeSettlementPostInteractionData()` |

### Signature Changes (Now Return Errors)

Several functions that previously could not fail now return errors for proper validation:

| Function | Old Return | New Return |
|----------|-----------|------------|
| `fusion.Extension.Keccak256()` | `*big.Int` | `(*big.Int, error)` |
| `fusionorder.FromPercent()` | `*Bps` | `(*Bps, error)` |
| `fusionorder.FromFraction()` | `*Bps` | `(*Bps, error)` |
| `orderbook.NewBitMask()` | `*BitMask` | `(*BitMask, error)` |
| `orderbook.TakerTraits.Encode()` | `*TakerTraitsEncoded` | `(*TakerTraitsEncoded, error)` |

`Must*` panic variants are provided for cases where failure is not expected:
- `fusionorder.MustNewBps()`, `fusionorder.MustFromPercent()`, `fusionorder.MustFromFraction()`
- `orderbook.MustNewBitMask()`

### Constants Changes

**Moved to `constants` package:**

| Old Location | New Location |
|-------------|-------------|
| `fusion.NativeToken` / `fusionplus.NativeToken` | `constants.NativeToken` |
| `fusion.ETHEREUM`, `fusion.POLYGON`, etc. | `constants.NetworkEthereum`, `constants.NetworkPolygon`, etc. |
| `fusion.NetworkEnum` / `fusionplus.NetworkEnum` | `constants.NetworkEnum` |

**Renamed:**

| Old Name | New Name |
|----------|----------|
| `constants.ERC20_APPROVE_GAS` | `constants.Erc20ApproveGas` |

**New:**
- `constants.ChainToWrapper` - wrapped native token addresses per chain
- `constants.GetWrappedToken()` - lookup convenience function
- `constants.ZeroAddress` - was previously in unexported `internal/addresses`
- `constants.Uint16Max`, `Uint24Max`, `Uint32Max`, `Uint40Max`, `Uint256Max`

**Removed:**
- `constants.AggregationRouterV5`, `AggregationV5RouterZkSyncEra`, `AggregationRouterV5Name`, `AggregationRouterV5VersionNumber` - unused V5 router constants
- `constants.AggregationRouterV5ABI` and `aggregationRouterV5.abi.json`

### `interface{}` Replaced with `any`

All public types now use `any` instead of `interface{}` (Go 1.18+):

| Type | Change |
|------|--------|
| `common.RequestPayload` | `Params any` |
| `common.HttpExecutor` | `ExecuteRequest(ctx, payload, v any) error` |
| `fusionorder.Keccak256Hash` | `data any` |
| `fusionplus` types | `CancelTx map[string]any`, etc. |
| `web3.PerformRpcCall*` | Returns `map[string]any` |

### Renamed Exported Symbols

| Old Name | New Name | Package |
|----------|----------|---------|
| `BitMask.ToString()` | `BitMask.String()` | `orderbook` |

### Bug Fixes

- **`fusion.PlaceOrders`**: Fixed HTTP method from `GET` to `POST`.
- **`fusion` validation**: Removed duplicate `WalletAddress` validation in `QuoterControllerGetQuoteWithCustomPresetsParamsFixed.Validate()`.
