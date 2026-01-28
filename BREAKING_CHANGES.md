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
| `fusion.NativeToken` | `fusionorder.NativeToken` |
| `fusion.NetworkEnum` | `fusionorder.NetworkEnum` |
| `fusion.ETHEREUM`, `fusion.POLYGON`, etc. | `fusionorder.ETHEREUM`, `fusionorder.POLYGON`, etc. |
| `fusion.GetWrappedToken()` | `fusionorder.GetWrappedToken()` |
| `fusion.ChainToWrapper` | `fusionorder.ChainToWrapper` |
| `fusion.NewAuctionDetails()` | `fusionorder.NewAuctionDetails()` |
| `fusion.DecodeAuctionDetails()` | `fusionorder.DecodeLegacyAuctionDetails()` |
| `fusion.CalcAuctionStartTime()` | `fusionorder.CalcAuctionStartTime()` |

#### Removed Exports from `fusionplus` Package

The following types and functions are no longer exported from the `fusionplus` package. Import them from `fusionorder` instead:

| Old Import | New Import |
|------------|------------|
| `fusionplus.NewInteraction()` | `fusionorder.NewInteraction()` |
| `fusionplus.DecodeInteraction()` | `fusionorder.DecodeInteraction()` |
| `fusionplus.NativeToken` | `fusionorder.NativeToken` |
| `fusionplus.NetworkEnum` | `fusionorder.NetworkEnum` |
| `fusionplus.ETHEREUM`, `fusionplus.POLYGON`, etc. | `fusionorder.ETHEREUM`, `fusionorder.POLYGON`, etc. |
| `fusionplus.GetWrappedToken()` | `fusionorder.GetWrappedToken()` |
| `fusionplus.ChainToWrapper` | `fusionorder.ChainToWrapper` |
| `fusionplus.NewAuctionDetails()` | `fusionorder.NewAuctionDetails()` |
| `fusionplus.DecodeAuctionDetails()` | `fusionorder.DecodeAuctionDetails()` |
| `fusionplus.CalcAuctionStartTime()` | `fusionorder.CalcAuctionStartTime()` |

#### Migration Example

Before:
```go
import "github.com/1inch/1inch-sdk-go/sdk-clients/fusion"

bps := fusion.NewBps(big.NewInt(100))
token, ok := fusion.GetWrappedToken(fusion.ETHEREUM)
```

After:
```go
import "github.com/1inch/1inch-sdk-go/sdk-clients/fusionorder"

bps := fusionorder.NewBps(big.NewInt(100))
token, ok := fusionorder.GetWrappedToken(fusionorder.ETHEREUM)
```

#### Type Aliases Remain (Internal Use)

The following type aliases remain in `fusion` and `fusionplus` packages. These are used internally and also available to external users through the existing packages:

**In `fusion`:**
- `fusion.Bps` = `fusionorder.Bps`
- `fusion.Interaction` = `fusionorder.Interaction`
- `fusion.AuctionDetails` = `fusionorder.AuctionDetails`
- `fusion.AuctionPointClassFixed` = `fusionorder.AuctionPointClassFixed`
- `fusion.GasCostConfigClassFixed` = `fusionorder.GasCostConfigClassFixed`
- `fusion.WhitelistItem` = `fusionorder.WhitelistItem`
- `fusion.ExtraData` = `fusionorder.ExtraData`
- `fusion.CustomPreset` = `fusionorder.CustomPreset`
- `fusion.TakingFeeInfo` = `fusionorder.TakingFeeInfo`
- `fusion.AuctionWhitelistItem` = `fusionorder.AuctionWhitelistItem`

**In `fusionplus`:**
- `fusionplus.Interaction` = `fusionorder.Interaction`
- `fusionplus.AuctionDetails` = `fusionorder.AuctionDetails`
- `fusionplus.AuctionPointClassFixed` = `fusionorder.AuctionPointClassFixed`
- `fusionplus.GasCostConfigClassFixed` = `fusionorder.GasCostConfigClassFixed`
- `fusionplus.WhitelistItem` = `fusionorder.WhitelistItem`
- `fusionplus.ExtraData` = `fusionorder.ExtraData`
- `fusionplus.CustomPreset` = `fusionorder.CustomPreset`
- `fusionplus.TakingFeeInfo` = `fusionorder.TakingFeeInfo`
- `fusionplus.AuctionWhitelistItem` = `fusionorder.AuctionWhitelistItem`

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

**Removed duplicate code:**
- `fusion.GenerateWhitelist` now aliases `fusionorder.GenerateWhitelist` instead of duplicating the logic

**Fixed error handling:**
- `fusionplus/escrowextension.go`: Replaced `log.Fatalf` calls with proper error returns in `decodeExtraData()`

**Fixed typos and error messages:**
- `fusionplus/extension_plus.go`: Fixed typo `postInteractoinDataEncoded` → `postInteractionDataEncoded`
- `fusionplus/extension_plus.go`: Fixed error message "fusion extension" → "fusionplus extension"

### New `fusionorder` Package

A new shared package has been created at `sdk-clients/fusionorder/` containing:

- **Types:** `Bps`, `Interaction`, `AuctionDetails`, `WhitelistItem`, `ExtraData`, `CustomPreset`, `TakingFeeInfo`, `AuctionWhitelistItem`, `NetworkEnum`, and more
- **Functions:** `NewBps`, `FromPercent`, `FromFraction`, `GetDefaultBase`, `NewInteraction`, `DecodeInteraction`, `NewAuctionDetails`, `DecodeAuctionDetails`, `DecodeLegacyAuctionDetails`, `CalcAuctionStartTime`, `IsNonceRequired`, `BpsToRatioFormat`, `GetWrappedToken`, `GenerateWhitelist`, `CanExecuteAt`, `IsExclusiveResolver`, and more
- **Constants:** `NativeToken`, `ETHEREUM`, `POLYGON`, `BINANCE`, `ARBITRUM`, `AVALANCHE`, `OPTIMISM`, `FANTOM`, `GNOSIS`, `COINBASE`, `ChainToWrapper`

This package is the single source of truth for shared order-related types and functions used by both `fusion` and `fusionplus` packages.
