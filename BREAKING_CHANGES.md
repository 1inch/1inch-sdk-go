# Breaking Changes

This document tracks breaking changes between major versions of the SDK that affect users importing and integrating the library.

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
