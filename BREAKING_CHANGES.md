# Breaking Changes

This document tracks breaking changes introduced in each version of the SDK.

## Version 2.0.0 (Upcoming)

### Architectural Refactoring: Shared `fusionorder` Package

The `fusion` and `fusionplus` packages have been refactored to share common types and utilities through a new `fusionorder` package. This eliminates code duplication and provides a cleaner architecture.

#### Removed Exports from `fusion` Package

The following types and functions are no longer exported from the `fusion` package. Import them from `fusionorder` instead:

| Old Import | New Import |
|------------|------------|
| `fusion.BpsZero` | `fusionorder.BpsZero` |
| `fusion.NewBps()` | `fusionorder.NewBps()` |
| `fusion.FromPercent()` | `fusionorder.FromPercent()` |
| `fusion.FromFraction()` | `fusionorder.FromFraction()` |
| `fusion.GetDefaultBase()` | `fusionorder.GetDefaultBase()` |
| `fusion.NewInteraction()` | `fusionorder.NewInteraction()` |
| `fusion.DecodeInteraction()` | `fusionorder.DecodeInteraction()` |
| `fusion.NativeToken` | `constants.NativeToken` |
| `fusion.NetworkEnum` | `constants.NetworkEnum` |
| `fusion.ETHEREUM`, `fusion.POLYGON`, etc. | `constants.NetworkEthereum`, `constants.NetworkPolygon`, etc. |
| `fusion.GetWrappedToken()` | `constants.GetWrappedToken()` |
| `fusion.ChainToWrapper` | `constants.ChainToWrapper` |
| `fusion.NewAuctionDetails()` | `fusionorder.NewAuctionDetails()` |
| `fusion.DecodeAuctionDetails()` | `fusionorder.DecodeLegacyAuctionDetails()` |
| `fusion.CalcAuctionStartTime()` | `fusionorder.CalcAuctionStartTime()` |
| `fusion.CalcAuctionStartTimeFunc` | `fusionorder.CalcAuctionStartTimeFunc` |
| `fusion.GenerateWhitelist()` | `fusionorder.GenerateWhitelist()` |
| `fusion.BpsToRatioFormat()` | `fusionorder.BpsToRatioFormat()` |

#### Removed Exports from `fusionplus` Package

The following types and functions are no longer exported from the `fusionplus` package. Import them from `fusionorder` instead:

| Old Import | New Import |
|------------|------------|
| `fusionplus.NewInteraction()` | `fusionorder.NewInteraction()` |
| `fusionplus.DecodeInteraction()` | `fusionorder.DecodeInteraction()` |
| `fusionplus.NativeToken` | `constants.NativeToken` |
| `fusionplus.NetworkEnum` | `constants.NetworkEnum` |
| `fusionplus.ETHEREUM`, `fusionplus.POLYGON`, etc. | `constants.NetworkEthereum`, `constants.NetworkPolygon`, etc. |
| `fusionplus.GetWrappedToken()` | `constants.GetWrappedToken()` |
| `fusionplus.ChainToWrapper` | `constants.ChainToWrapper` |
| `fusionplus.NewAuctionDetails()` | `fusionorder.NewAuctionDetails()` |
| `fusionplus.DecodeAuctionDetails()` | `fusionorder.DecodeAuctionDetails()` |
| `fusionplus.CalcAuctionStartTime()` | `fusionorder.CalcAuctionStartTime()` |
| `fusionplus.CalcAuctionStartTimeFunc` | `fusionorder.CalcAuctionStartTimeFunc` |

#### Migration Example

Before:
```go
import "github.com/1inch/1inch-sdk-go/sdk-clients/fusion"

bps := fusion.NewBps(big.NewInt(100))
token, ok := fusion.GetWrappedToken(fusion.ETHEREUM)
```

After:
```go
import (
    "github.com/1inch/1inch-sdk-go/common/fusionorder"
    "github.com/1inch/1inch-sdk-go/constants"
)

bps := fusionorder.NewBps(big.NewInt(100))
token, ok := constants.GetWrappedToken(constants.NetworkEthereum)
```

#### Type Aliases Kept (Required for Public API)

Only type aliases needed by users at the public API level remain:

**In `fusion`:**
- `fusion.TakingFeeInfo` = `fusionorder.TakingFeeInfo` (used in `OrderParams.Fee`)
- `fusion.CustomPreset` = `fusionorder.CustomPreset` (used in `GetQuoteWithCustomPreset()`)
- `fusion.CustomPresetPoint` = `fusionorder.CustomPresetPoint` (used in `CustomPreset.Points`)

**In `fusionplus`:**
- `fusionplus.TakingFeeInfo` = `fusionorder.TakingFeeInfo`
- `fusionplus.CustomPreset` = `fusionorder.CustomPreset`
- `fusionplus.CustomPresetPoint` = `fusionorder.CustomPresetPoint`
- `fusionplus.Interaction` = `fusionorder.Interaction`
- `fusionplus.AuctionDetails` = `fusionorder.AuctionDetails`
- `fusionplus.AuctionPointClassFixed` = `fusionorder.AuctionPointClassFixed`
- `fusionplus.GasCostConfigClassFixed` = `fusionorder.GasCostConfigClassFixed`
- `fusionplus.WhitelistItem` = `fusionorder.WhitelistItem`
- `fusionplus.ExtraData` = `fusionorder.ExtraData`
- `fusionplus.AuctionWhitelistItem` = `fusionorder.AuctionWhitelistItem`

#### Type Aliases Removed (Internal-only types)

The following type aliases have been removed from `fusion` - use `fusionorder.*` directly:

| Removed Alias | Use Instead |
|---------------|-------------|
| `fusion.AuctionDetails` | `fusionorder.AuctionDetails` |
| `fusion.AuctionPointClassFixed` | `fusionorder.AuctionPointClassFixed` |
| `fusion.GasCostConfigClassFixed` | `fusionorder.GasCostConfigClassFixed` |
| `fusion.WhitelistItem` | `fusionorder.WhitelistItem` |
| `fusion.AuctionWhitelistItem` | `fusionorder.AuctionWhitelistItem` |
| `fusion.ExtraData` | `fusionorder.ExtraData` |
| `fusion.Bps` | `fusionorder.Bps` |
| `fusion.Interaction` | `fusionorder.Interaction` |

These types are only used internally by the SDK and users do not need to construct them directly.

### Deleted Files

The following wrapper files have been deleted:

**From `fusion/`:**
- `bps.go` - Re-exports moved to `fusionorder`
- `interaction.go` - Re-exports moved to `fusionorder`
- `nativetokenwrappers.go` - Re-exports moved to `fusionorder`
- `auctiondetails.go` - Re-exports moved to `fusionorder`
- `bytesiter.go` - Unused dead code (use `internal/bytesiterator` instead)
- `custompreset.go` - Validation logic moved to `fusionorder`

**From `fusionplus/`:**
- `interaction.go` - Re-exports moved to `fusionorder`
- `nativetokenwrappers.go` - Re-exports moved to `fusionorder`
- `auctiondetails.go` - Re-exports moved to `fusionorder`

### Consolidated Duplicate Functions

The following functions that existed in multiple packages are now only in `fusionorder`:

- `CalcAuctionStartTime()` - Now only in `fusionorder`
- `IsNonceRequired()` - Now only in `fusionorder`
- `BpsToRatioFormat()` - Now only in `fusionorder`
- `CanExecuteAt()` - Shared helper for whitelist execution checks
- `IsExclusiveResolver()` - Shared helper for whitelist resolver checks
- `GenerateWhitelistFromItems()` - Shared whitelist generation with sorting

### Consolidated Types in `fusionplus`

- `SettlementPostInteractionData` and `SettlementPostInteractionDataFusion` merged into single `SettlementPostInteractionData` type with optional fee fields
- `settlementpostinteractiondatafusion.go` deleted - functionality merged into `settlementpostinteractiondata.go`
- `NewSettlementPostInteractionDataFusion()` renamed to `NewSettlementPostInteractionDataWithFees()`
- `CreateSettlementPostInteractionDataFusion()` renamed to `CreateSettlementPostInteractionDataWithFees()`
- `DecodeFusion()` renamed to `DecodeSettlementPostInteractionData()`

### Renamed Types and Functions in `fusionplus`

Types with the `Fusion` suffix have been renamed for consistency. Types that need to be distinguished from `fusion` package equivalents now use the `Plus` suffix:

| Old Name | New Name |
|----------|----------|
| `ExtensionFusion` | `ExtensionPlus` |
| `ExtensionParamsFusion` | `ExtensionParamsPlus` |
| `NewExtensionFusion()` | `NewExtensionPlus()` |
| `CreateAuctionDetailsFusion()` | `CreateAuctionDetailsPlus()` |
| `PresetClassFixedFusion` | `PresetClassFixed` |
| `GasCostConfigClassFusion` | `GasCostConfigClass` |
| `AuctionPointClassFusion` | `AuctionPointClass` |
| `SettlementSuffixDataFusion` | Removed (duplicate of `SettlementSuffixData`) |

**Files renamed/deleted:**
- `extension_fusion.go` → `extension_plus.go`
- `interaction_test.go` deleted (was testing `fusionorder` functions)

**Variable names updated** (internal, for consistency):
- `fusionExtension` → `extensionPlus`
- `auctionPointsFusion` → `auctionPointsPlus`
- `gasCostsFusion` → `gasCostsPlus`
- `presetFusion` → `presetPlus`
- `auctionDetailsFusion` → `auctionDetailsPlus`

**Removed unused type aliases:**
- `FeesFusion` (was alias for `Fees`)
- `IntegratorFeeFusion` (was alias for `IntegratorFee`)
- `DetailsFusion` (was alias for `Details`)

### Code Quality Improvements

**Removed re-exports and aliases:**
- `fusion.GenerateWhitelist` removed - use `fusionorder.GenerateWhitelist` directly
- `fusion.CalcAuctionStartTimeFunc` removed - use `fusionorder.CalcAuctionStartTimeFunc` directly
- `fusion.CalcAuctionStartTime` removed - use `fusionorder.CalcAuctionStartTime` directly
- `fusion.BpsToRatioFormat` removed - use `fusionorder.BpsToRatioFormat` directly
- `fusionplus.CalcAuctionStartTimeFunc` removed - use `fusionorder.CalcAuctionStartTimeFunc` directly
- `fusionplus.CalcAuctionStartTime` removed - use `fusionorder.CalcAuctionStartTime` directly

**Fixed error handling:**
- `fusionplus/escrowextension.go`: Replaced `log.Fatalf` calls with proper error returns in `decodeExtraData()`

**Fixed typos and error messages:**
- `fusionplus/extension_plus.go`: Fixed typo `postInteractoinDataEncoded` → `postInteractionDataEncoded`
- `fusionplus/extension_plus.go`: Fixed error message "fusion extension" → "fusionplus extension"

### New `fusionorder` Package

A new shared package has been created at `common/fusionorder/` containing:

- **Types:** `Bps`, `Interaction`, `AuctionDetails`, `WhitelistItem`, `ExtraData`, `CustomPreset`, `TakingFeeInfo`, `AuctionWhitelistItem`, and more
- **Functions:** `NewBps`, `FromPercent`, `FromFraction`, `GetDefaultBase`, `NewInteraction`, `DecodeInteraction`, `NewAuctionDetails`, `DecodeAuctionDetails`, `DecodeLegacyAuctionDetails`, `CalcAuctionStartTime`, `IsNonceRequired`, `BpsToRatioFormat`, `GenerateWhitelist`, `CanExecuteAt`, `IsExclusiveResolver`, `ToOrderbookExtension`, and more

This package is the single source of truth for shared order-related types and functions used by both `fusion` and `fusionplus` packages.

### Constants Consolidation

Token and network constants have been moved to the top-level `constants` package:

| Constant | Location |
|----------|----------|
| `NativeToken` | `constants.NativeToken` |
| `ZeroAddress` | `constants.ZeroAddress` |
| `NetworkEnum` | `constants.NetworkEnum` |
| `NetworkEthereum`, `NetworkPolygon`, etc. | `constants.NetworkEthereum`, `constants.NetworkPolygon`, etc. |
| `ChainToWrapper` | `constants.ChainToWrapper` |
| `GetWrappedToken()` | `constants.GetWrappedToken()` |
| `Uint16Max`, `Uint24Max`, `Uint32Max`, `Uint40Max`, `Uint256Max` | `constants.Uint16Max`, etc. |

The `internal/addresses` package has been removed; use `constants.ZeroAddress` instead.
