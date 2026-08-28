# Breaking Changes

This document tracks breaking changes between major versions of the SDK that affect users importing and integrating the library.

For a step-by-step upgrade checklist, see [MIGRATION.md](MIGRATION.md); this document is the per-change reference it links back to.

## Unreleased

### Generated Type Names No Longer Leak Upstream Controller Naming

Generated type names previously mirrored 1inch's internal operation ids (`QuoterControllerGetQuoteWithCustomPresetsParams`, `TxProcessorApiControllerBroadcastTransactionJSONRequestBody`, …). Every operation now has a friendly name in `codegen/mapping.json`, and key schemas are renamed via `x-go-name` overrides, so the SDK reads as intended:

```go
// Before
quote, err := client.GetQuote(ctx, fusionplus.QuoterControllerGetQuoteParams{...}) // returns *GetQuoteOutput

// After
quote, err := client.GetQuote(ctx, fusionplus.QuoteParams{...}) // returns *fusionplus.Quote
```

**Compatibility:** every renamed identifier (55 types, 14 constants) keeps its old name as a `// Deprecated:` alias in each package's `deprecated_names.go` — aliases are identical types, so existing code compiles and behaves identically. The one exception: the `orderbook` wrapper params (`GetAllOrdersParams`, `GetCountParams`, `GetEventsParams`, `GetOrdersByCreatorAddressParams`) embed the generated query types, and an embedded field's name is its type name — struct literals that named the embedded field (e.g. `GetAllOrdersParams{LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams: q}`) must use the new field name (`LimitOrdersQueryParams: q`). Promoted field access (`params.Page`) and assignment through a variable are unaffected.

### `*Fixed` Types Are Now Aliases of Corrected Generated Types

Known type bugs in the upstream OpenAPI specs are now corrected at generation time (`codegen/overrides.go`) instead of via hand-maintained `*Fixed` shadow copies. All 12 `*Fixed` types in `fusion`, `fusionplus`, and `tokens` are now type aliases of the corrected generated types.

**Impact:**
- `fusion` and `fusionplus` aliases keep the same field sets and types — no code changes needed. (`fusionplus.QuoterControllerGetQuoteWithCustomPresetsParamsFixed` gains `Fee *big.Int` and `IsPermit2 bool` fields; additive.)
- `tokens.ProviderTokenDtoFixed` / `TokenInfoDtoFixed`: optional fields (`DisplayedSymbol`, `Eip2612`, `IsFoT`, `LogoURI`, `Extensions`) changed from pointers to values, matching the SDK-wide policy of generating optional fields without pointers. `Tags` is `[]TagDto` on the generated types too.
- **The raw (non-`Fixed`) generated types now carry the corrections**, so code using them directly sees type changes: `fusion.GetQuoteOutput.QuoteId` and `fusionplus.GetQuoteOutput.QuoteId` (`map[string]interface{}` → `string`), `fusion.PresetClass.ExclusiveResolver` (`map[string]interface{}` → `string`), `Amount` on all fusion/fusionplus quoter params (`float32` → `string`), `fusionplus` quoter params `IsPermit2` (`string` → `bool`) — including `QuoterControllerBuildQuoteTypedDataParams` — and `fusionplus.GetOrderFillsByHashOutput` (`DstTokenPriceUsd`/`SrcTokenPriceUsd` `map[string]interface{}` → `string`, `Points` → `[]AuctionPointOutput`). These now match what the `Fixed` types (the types the SDK's own methods use) have always had.
- **`Fee` on the fusionplus quoter params changed from `*big.Int` to `int`** (including the `Fixed` alias, where it was `*big.Int` since the type was introduced). This also fixes a silent bug: `go-querystring` cannot encode `*big.Int` (a struct with only unexported fields), so a fee set on a quote request was **silently omitted from the query string** and never reached the API. Fees are integral basis points (1% = 100 bps), so `int` is both correct and wire-safe; a wire-format regression test now pins the encoding.

**Migration (tokens only):**

```go
// Before
if token.DisplayedSymbol != nil { use(*token.DisplayedSymbol) }
eip2612 := token.Eip2612 != nil && *token.Eip2612

// After
if token.DisplayedSymbol != "" { use(token.DisplayedSymbol) }
eip2612 := token.Eip2612
```

The `*Fixed` names remain available as aliases, now marked `// Deprecated:` (IDEs and linters flag usages; pkg.go.dev strikes them through). All SDK method signatures, examples, and tests use the underlying generated names; aliases are identical types, so this changes nothing for callers.

Two related changes:

- **`fusionplus.CreateAuctionDetailsPlus` now takes `*Preset`** (the generated quoter preset) instead of `*PresetClassFixed`. The legacy `PresetClassFixed`/`GasCostConfigClass`/`AuctionPointClass` bridge shape is no longer produced or consumed by the SDK and is kept only as deprecated types.
- **`fusionorder.AuctionPointClassFixed` and `fusionorder.GasCostConfigClassFixed` are renamed to `fusionorder.AuctionPoint` and `fusionorder.GasCostConfig`**, with deprecated aliases under the old names.
- Minor encoding nuance: `fusion.Quote.SurplusFee` and `fusion.Quote.MarketAmount` (previously hand-added on `GetQuoteOutputFixed`) now carry `omitempty`, so re-marshaling a quote omits them when zero. Decoding is unchanged.

### Type Generator Upgraded to oapi-codegen v2

The deprecated `deepmap/oapi-codegen` v1.16.2 has been replaced by the maintained `oapi-codegen/oapi-codegen` v2.8.0. The full public-API delta was machine-verified with `gorelease -base=v4.1.0`; the incompatible changes are:

- **`spotprices` and `traces` bare enum constants are now prefixed with their enum type name** (v1 only prefixed on collision; v2 is consistent).

  ```go
  // Before
  params.Currency = spotprices.GetPricesRequestDtoCurrency(spotprices.USD)
  trace.Type == traces.CALL

  // After
  params.Currency = spotprices.GetPricesRequestDtoCurrencyUSD
  trace.Type == traces.CoreCustomRootTxEventCallstackTraceFullDtoTypeCALL
  ```

- **Optional object-typed response fields lost their pointers**: `history.TransactionDetailsDto.Meta` (`*TransactionDetailsMetaDto` → `TransactionDetailsMetaDto`), `nft.Asset.Collection` (`*Collection` → `Collection`), and `nft.Asset.RarityData` (`*RarityData` → `RarityData`). v1 ignored the SDK's skip-optional-pointer annotation on object-typed properties; v2 honors it. Nil checks become zero-value checks; when a response omits the field you now get a zero struct instead of nil.

- **`web3.ApiKeyAuthScopes` constant removed** (a security-scheme artifact v2 no longer emits for types-only generation; it referenced nothing usable).

Additions (non-breaking): every generated enum type gains a `Valid() bool` method, response types gain doc comments from the spec, and `traces` gains typed union response helpers (`As*/From*/Merge*`).

### Chain-ID Fields Changed from `float32` to `int`

The OpenAPI specs type chain ids as `number`, which oapi-codegen generated as `float32`. A `float32` cannot exactly represent integers above 2²⁴ (16,777,216), so Aurora's chain id (1313161554) was silently rounded to 1313161600 the moment it was assigned. This corrupted the EIP-712 signing domain and the escrow extension's encoded `DstChainId`, and made Aurora impossible to pass through parameter validation (which correctly rejected the rounded value with a misleading error). All chain-id fields in the `fusionplus` and `tokens` packages are now `int`. The codegen pipeline (`codegen/transforms.go`) now rewrites chain-id fields typed as `number` to `integer` before generation, so regenerated types stay correct when specs are refreshed.

**Affected exported types:**

| Package | Type | Field(s) |
|---------|------|----------|
| `fusionplus` | `QuoterControllerGetQuoteParamsFixed` | `SrcChain`, `DstChain` |
| `fusionplus` | `QuoterControllerGetQuoteWithCustomPresetsParamsFixed` | `SrcChain`, `DstChain` |
| `fusionplus` | `GetSettlementContractParams` | `ChainId` |
| `fusionplus` | `EscrowExtension`, `EscrowExtensionParams`, `EscrowExtraData` | `DstChainId` |
| `fusionplus` | `SignedOrderInput` | `SrcChainId` |
| `fusionplus` | `GetOrderFillsByHashOutputFixed`, `GetOrderFillsByHashOutput`, `ActiveOrdersOutput`, `ReadyToExecutePublicAction`, `OrderApiControllerGetActiveOrdersParams`, `OrderApiControllerGetOrdersByMakerParams`, `OrderApiControllerGetSettlementContractParams`, `QuoterControllerGetQuoteParams`, `QuoterControllerGetQuoteWithCustomPresetsParams`, `QuoterControllerBuildQuoteTypedDataParams` | chain-id fields |
| `tokens` | `ProviderTokenDtoFixed`, `TokenInfoDtoFixed`, `ProviderTokenDto`, `TokenDto`, `TokenInfoDto` | `ChainId` |

**Impact:** Code that assigns a `float32`-typed variable or an explicit `float32(...)` conversion to these fields will no longer compile. Code using untyped constants — including all `constants.*ChainId` values and integer literals — compiles unchanged.

**Migration:** Remove `float32` conversions:

```go
// Before
params := fusionplus.QuoterControllerGetQuoteParamsFixed{
    SrcChain: float32(constants.ArbitrumChainId),
    DstChain: float32(constants.BaseChainId),
}

// After
params := fusionplus.QuoterControllerGetQuoteParamsFixed{
    SrcChain: constants.ArbitrumChainId,
    DstChain: constants.BaseChainId,
}
```

The internal validators `CheckChainIdFloat32` and `CheckChainIdFloat32Required` were removed (internal package; not importable by consumers).

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
