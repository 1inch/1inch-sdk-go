# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) starting with the *v1.0.0-beta.1* release.

## [Unreleased]

### Added
- MIGRATION.md: a v4 → v5 upgrade guide (checklist + compiler-error fixes), complementing the per-change reference in BREAKING_CHANGES.md.
- Fee validation on fusion/fusionplus quote params: fees must be basis points in [0, 10000].

### Changed
- All `*Fixed` types are now deprecated aliases; SDK signatures use the generated names.
- Codegen pipeline hardened: rewritten from bash/jq into a Go tool with a declarative override table (replacing hand-maintained `*Fixed` copies), a spec-provenance lock, and a weekly upstream-drift workflow. Output is byte-for-byte reproducible.
- CI expanded: public-API compatibility gate (`gorelease`), wire-surface (struct-tag) diff, decoder fuzzing, `govulncheck`, and weekly live-API smoke tests.

### Fixed
- **Funds-safety:** `fusionplus.DecodeEscrowExtension` no longer corrupts escrow terms on decode → re-encode (swapped safety deposits; hex/decimal mismatch on deposits and hashlock). A round-trip test pins it.
- **Funds-safety:** `EscrowExtension.ConvertToOrderbookExtension` is now pure (it previously mutated its receiver, corrupting the extension/salt on a second call).
- `fusionplus.DecodeEscrowExtension` returns an error instead of panicking on malformed input.
- HTTP client now applies a 60s default timeout and caps response bodies at 64 MiB (previously could hang / read unbounded).
- `fusionplus` quote `Fee` is now transmitted — it was silently dropped (`*big.Int` cannot be query-encoded); the field is now `int` basis points.

### Deprecated
- Generated type names are now intent-based (e.g. `fusionplus.QuoteParams`/`Quote`), and the v5 surface cleanup renamed types/methods across all clients for consistency. **These are source-compatible:** every old name is kept as a `type Old = New` alias or a forwarding method, so existing code compiles and behaves identically — it just gets a deprecation warning. Migrate before a future major removes them. Full list in MIGRATION.md §4.

### Breaking Changes
_Only genuine incompatibilities — code updates required. Renames with aliases are under Deprecated above, not here._
- Chain-id fields are now `int` instead of `float32` (which silently corrupted Aurora's chain id in the signing domain); drop any `float32(...)` conversions.
- Type/shape corrections now on the generated types (previously only on the `*Fixed` copies): `Amount` → `string`, `Fee` → `int`, `QuoteId` → `string`, `tokens` optional fields are values not pointers, `Tags` → `[]TagDto`. The `*Fixed` names still resolve (aliases), but field access on the old types must be updated (e.g. drop `*token.Eip2612`).
- oapi-codegen v1 → v2: bare enum constants are now type-prefixed (`spotprices.USD` → `spotprices.GetPricesForRequestedTokensParamsCurrencyUSD`) — the bare names are gone (no alias); `web3.ApiKeyAuthScopes` removed.
- No-alias-possible breaks: `common.RequestPayload.U` → `.Path`, `fusionplus.Order.EscExtension` → `.EscrowExtension`, `traces.NewConfiguration` now takes `ConfigurationParams`, `constants` maps keyed by `int` (`NetworkEnum` deprecated), removed dead types, and struct literals that named the renamed embedded `orderbook` query field.
- See BREAKING_CHANGES.md for per-change detail and MIGRATION.md for upgrade steps.

## [v4.2.0] - 2026-08-28

### Added
- `constants.ChainToTrueERC20` and `constants.GetTrueERC20` give the Fusion+ TRUE_ERC20 sentinel token for each supported chain.

### Fixed
- **Fusion+ orders had the wrong taker asset on the source chain.** `CreateFusionPlusOrderData` put the destination token into the taker asset of the source-chain order. For a native destination, it used the wrapped token of the source chain. The taker asset is now the TRUE_ERC20 sentinel of the chain. This sentinel is a token that does nothing when a resolver transfers it. The escrow extension now holds the real destination token in its `DstToken` field. This includes the native sentinel `0xEeee…`. The destination token can also be a live ERC-20 on the source chain. In this condition, the old behavior could transfer a real and unwanted token on the source chain during a fill. The SDK now stops with an error when the source chain has no TRUE_ERC20 sentinel. It does not sign an order that has the wrong asset.

## [v4.1.0] - 2026-07-25

### Added
- New method `fusion.Client.PlaceOrderFromParams`: fetches a quote and places the order in one call, so settings like `Permit` and `IsPermit2` are supplied once and propagate to both the quote request and the order
- New function `orderbook.DecodeMakerTraits`: parses an encoded maker traits value back into a `MakerTraits` struct, enabling flag reads like `ShouldUsePermit2` from on-chain orders
- New functions `orderbook.BuildPermit2Calldata` and `orderbook.GetPermit2Allowance`: sign a Permit2 AllowanceTransfer PermitSingle (the 352-byte form the protocol fills) and read Permit2 allowance state. The protocol's 96-byte compact permit2 form is deliberately not offered and is rejected by order validation, because fills through the deployed Aggregation Router v6 revert on it
- Mainnet-fork integration tests under `tests/integration` (build tag `integration`, `make test-integration`)
- New methods `fusionplus.Client.GetActiveOrders` and `fusionplus.Client.GetSettlementContract`: list open cross-chain orders and fetch the escrow factory address

### Fixed
- **`fusion.CreateFusionOrderData`**: `OrderParams.Permit` and `OrderParams.IsPermit2` are now honored
- **`fusionplus.CreateFusionPlusOrderData`**: `OrderParams.IsPermit2` now sets the `USE_PERMIT2` maker-traits bit.
- **`fusionplus.FromLimitOrderExtension`**: post-interaction data now decodes correctly
- **Permit input validation**: odd-length permit hex is now rejected by `CheckPermitHash`, `fusion.NewExtension`, and `fusionplus.NewExtensionPlus`; it previously corrupted the encoded extension and produced orders that could never fill.
- **User-Agent header**: API requests now report the SDK release version from `internal/version.Version`
- **`common.Wallet.Call`**: wallets created without a node URL now return an error instead of panicking on on-chain calls
- **`orderbook.BuildOrderExtensionBytes`**: a hex string cast to `[]byte` in `MakerPermit` is now rejected with an error
- **Transaction fee cap**: EIP-1559 transactions built without an explicit `SetGasFeeCap` now default the fee cap to twice the node's suggested gas price instead of the bare suggestion, and gas estimation no longer sends a gas price

### Changed
- The maker permit token field is encoded in lowercase hex in `fusion` and `fusionplus` extensions.
- The "Release new version" workflow now updates `internal/version/version.go` and rolls the CHANGELOG `Unreleased` section into a dated version heading as part of each stable release, and uses that CHANGELOG section as the GitHub release notes (falling back to auto-generated notes when no entries were written). PR CI blocks manual edits to the version constant and warns when code changes ship without a CHANGELOG entry.

## [v4.0.0] - 2026-07-14

### Breaking Changes
- **Module path now includes `/v4` suffix**: The module path is now `github.com/1inch/1inch-sdk-go/v4`, as required by Go for major versions >= 2. All imports and `go get` commands must include the `/v4` suffix (e.g. `github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation`).
- **Minimum Go version raised to 1.25**: The `go` directive in `go.mod` is now `go 1.25.0` (previously `go 1.22`). This is required by `golang.org/x/crypto` v0.52.0. Consumers must build with Go 1.25 or newer.

### Changed
- **Dependency security upgrades** (resolves open Dependabot alerts):
  - `golang.org/x/crypto` v0.31.0 → v0.52.0
  - `github.com/ethereum/go-ethereum` v1.14.13 → v1.17.0
  - `github.com/consensys/gnark-crypto` v0.12.1 → v0.18.1 (transitive)

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


