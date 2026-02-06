# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) starting with the *v1.0.0-beta.1* release.

## [v3.0.0] - 2026-02-06

### Breaking Changes
- **New `fusionorder` package**: Common types and functions from `fusion` and `fusionplus` have been consolidated into `common/fusionorder/`. Types like `Bps`, `Interaction`, `AuctionDetails`, `WhitelistItem` and functions like `NewBps()`, `CalcAuctionStartTime()`, `GenerateWhitelist()` are now in `fusionorder`.
- **Constants consolidated**: `NativeToken`, `NetworkEnum`, `ETHEREUM`/`POLYGON`/etc. moved from `fusion`/`fusionplus` to `constants`. Network constants renamed to Go conventions (e.g., `ETHEREUM` → `NetworkEthereum`).
- **Renamed types in `fusionplus`**: Types with `Fusion` suffix renamed to `Plus` suffix (e.g., `ExtensionFusion` → `ExtensionPlus`). Redundant types like `FeesFusion`, `DetailsFusion` removed.
- **Merged types in `fusionplus`**: `SettlementPostInteractionDataFusion` merged into `SettlementPostInteractionData`. `DecodeFusion()` → `DecodeSettlementPostInteractionData()`.
- **`interface{}` replaced with `any`**: Affects public types including `common.RequestPayload`, `common.HttpExecutor`, and `web3` return types.
- **Constant renamed**: `constants.ERC20_APPROVE_GAS` → `constants.Erc20ApproveGas`.
- **Removed V5 router constants**: `AggregationRouterV5`, `AggregationV5RouterZkSyncEra`, and related V5 constants/ABI removed.
- **`BitMask.ToString()`** renamed to `BitMask.String()` in `orderbook`.
- **Signature changes**: Several functions now return errors — `Extension.Keccak256()`, `FromPercent()`, `FromFraction()`, `orderbook.NewBitMask()`, `orderbook.TakerTraits.Encode()`. `Must*` panic variants added where appropriate.
- **Deprecated type aliases**: `fusion.TakingFeeInfo`, `fusion.CustomPreset`, `fusion.CustomPresetPoint` (and `fusionplus` equivalents) still work but are deprecated in favor of `fusionorder.*`.
- See `BREAKING_CHANGES.md` for full migration guide with tables.

### Added
- New `fusionorder` package with shared types and functions for `fusion` and `fusionplus`
- New public constants: `constants.ChainToWrapper`, `constants.GetWrappedToken()`, `constants.ZeroAddress`, `constants.Uint16Max`/`Uint24Max`/`Uint32Max`/`Uint40Max`/`Uint256Max`
- `Must*` panic variants: `MustNewBps()`, `MustFromPercent()`, `MustFromFraction()`, `MustNewBitMask()`

### Fixed
- **`fusion.PlaceOrders`**: HTTP method changed from `GET` to `POST` (sending a body with `GET` is semantically incorrect).
- **`fusion` validation**: Removed duplicate `WalletAddress` check in `QuoterControllerGetQuoteWithCustomPresetsParamsFixed.Validate()`.
- **`fusionplus`**: Replaced `log.Fatalf` calls with proper error returns (library no longer terminates the calling process on decode errors).
- **Error wrapping**: Standardized `%v` to `%w` in `fmt.Errorf` calls for proper `errors.Is`/`errors.As` support.

### Changed
- Eliminated code duplication between `fusion` and `fusionplus` packages

## [v2.0.0] - 2025-11-05
[v2.0.0 release page](https://github.com/1inch/1inch-sdk-go/releases/tag/v2.0.0)

### Changed
- Fusion Plus updated to use v1.1 API

## [v2.0.0-preview.2] - 2025-10-30
[v2.0.0-preview.2 release page](https://github.com/1inch/1inch-sdk-go/releases/tag/v2.0.0-preview.2)

### Breaking Changes
- Limit Orders have been refactored. Order creation now uses a different flow. See the examples for more details.

### Changed
- Limit Order SDK updated to support v4.1 API

## [v2.0.0-preview] - 2025-1-22
[v2.0.0-preview release page](https://github.com/1inch/1inch-sdk-go/releases/tag/v2.0.0-preview)

### Breaking Changes
- a new `surplus=true` query parameter must be added to Fusion quote requests

### Changed
- Fusion implementation updated to support new Fusion backend
- Fusion+ is disabled until refactor is complete

## [v1.0.0-beta.3] - 2025-1-22
[v1.0.0-beta.3 release page](https://github.com/1inch/1inch-sdk-go/releases/tag/v1.0.0-beta.3)
### Changed
- Fusion Plus support added
- Pending Fusion orders can now be tracked using the SDK
- Orderbook client updated to support new API schema

## [v1.0.0-beta.2] - 2024-10-23
[v1.0.0-beta.2 release page](https://github.com/1inch/1inch-sdk-go/releases/tag/v1.0.0-beta.2)
### Changed
- Classic Swap updated to use V6 API
- Added examples for all Classic Swap endpoints
- When using TransactionBuilder, if no `gas` value is specified in the transaction config, `eth_estimateGas` will be used by default

## [v1.0.0-beta.1] - 2024-8-22
[v1.0.0-beta.1 release page](https://github.com/1inch/1inch-sdk-go/releases/tag/v1.0.0-beta.1)

Note: This changelog summarizes all changes since the last *changelog* version of v0.0.3-developer-preview

### Added
- Web3 API added
- Fusion SDK added
- Portfolio API added
- Permit1 support added for Orderbook orders and Aggregator Swaps

### Changed
- Readme updated to link to all API docs and examples
- Updating Geth version
- Types generation script updated to handle Web3 API spec design
- Normalized and improved SDK examples
- Improved code generation to make optional parameters pointers

# [v0.0.3-developer-preview] - 2024-3-9
[v0.0.3-developer-preview](https://github.com/1inch/1inch-sdk/releases/tag/v0.0.3-developer-preview)

### New Features and Enhancements:

- All non-global query configurations have been moved to the request-level
  params ([PR](https://github.com/1inch/1inch-sdk/pull/6))
    - RPC providers for all chains will now be defined/set at SDK startup
- Query parameters now use concrete types instead of pointers ([PR](https://github.com/1inch/1inch-sdk/pull/16))
- Limit orders created within the SDK now support auto-expiration ([PR](https://github.com/1inch/1inch-sdk/pull/23))
- Permit1 properly supported for limit orders when possible (fallback to Approval if Permit1 does not
  work) ([commit](https://github.com/1inch/1inch-sdk/commit/f2e79e5f0e81503bfeeff076e41455e86e5a5120))
- When creating a limit order, integrators can error out when an approval is needed. This is useful for integrators who
  want all onchain actions to be performed manually by their users ([PR](https://github.com/1inch/1inch-sdk/pull/26))

### Optimizations and Bug Fixes:

- Tenderly forks are cleaned up automatically at the beginning of each test
  run ([PR](https://github.com/1inch/1inch-sdk/pull/6))
- Validation pattern for swagger-generated input params is now fully handled on all
  endpoints ([PR](https://github.com/1inch/1inch-sdk/pull/8))
- Project-wide validation scripts added to verify validation logic
  standards ([PR](https://github.com/1inch/1inch-sdk/pull/11))

# [v0.0.2-developer-preview] 2024-1-23
Tag: [v0.0.2-developer-preview](https://github.com/1inch/1inch-sdk/releases/tag/v0.0.2-developer-preview)

### New Features and Enhancements:

- **Added Tenderly support for e2e swap tests**
    - e2e tests will now create forks, apply state overrides, and run simulations when a Tenderly API key is provided.
- **Added approval type selection**
    - Users can choose between `Approve` and `Permit1` (`Permit2` currently unsupported)
- **Implemented nonce cache to address RPC lag**
    - Once a wallet has posted a transaction, the nonce of that transaction is tracked and incremented internally by the
      SDK.

### Optimizations and Bug Fixes:

- Updated orderbook to use string inputs instead of integers to support all of uint256.
- Increased gas limit and reduced permit duration to improve transactions success and debugging.
- Moved Actions into a service within the main client to consolidate SDK structure.
- Simplified tests and refactored onchain actions to have more uniformity across the library.

# Release (January 15, 2024)

Tag: [v0.0.1-developer-preview](https://github.com/1inch/1inch-sdk/releases/tag/v0.0.1-developer-preview)

### New Features and Enhancements:

### Limit Order support

- Enables posting orders to the 1inch Limit Order Protocol
- Enables reading orders from the 1inch Limit Order Protocol
- Most endpoints from the Limit Order API supported
    - `has-active-orders-with-permit` REST endpoint still untested

### Aggregator Protocol support

- All REST endpoints supported
- Get quotes and swap data from the Aggregator Protocol

### Onchain execution support

- Execute swaps onchain from within the SDK


