# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) starting with the *v1.0.0-beta.1* release.

## [Unreleased]

### Added
- Wire-surface diff in CI (`tools/wiredump` + `wire-surface` job): every PR's job summary shows the exact diff of all exported struct fields (name, Go type, json/url tags) against the base branch — catching struct-tag changes that compile fine but silently alter the wire format, which type-level API diffing cannot see. An A/B run against the pre-refactor tree verified all 73 wire-surface differences on this branch map to documented changes, with zero field removals and zero tag changes.
- Public API compatibility gate in CI: every PR gets a machine-generated `gorelease` diff against the latest stable release in its job summary, and incompatible changes fail the build unless BREAKING_CHANGES.md has an Unreleased section documenting them. The release workflow refuses patch/minor (and rc) releases that contain incompatible public API changes.
- Wire-format regression tests pinning the exact query-string and JSON encodings of the fusionplus quoter params, relayer order submission, and response decoding for fusionplus and tokens (`wire_test.go`)
- Live API smoke tests (`go test -tags liveapi ./tests/liveapi/`, requires `DEV_PORTAL_TOKEN`): read-only production requests validating that live responses still decode into the SDK's types; run weekly by the spec-drift workflow
- Codegen determinism test: two pipeline runs must produce byte-identical output
- Characterization test suite for the codegen pipeline (`codegen/codegen_test.go`): asserts that `generate_types.sh` reproduces every committed `*_types.gen.go` file byte-for-byte from the committed specs, that the specs are a fixed point of the in-place transforms, and documents each transform's behavior via synthetic fixtures. This makes the pipeline safe to restructure.

### Changed
- All `*Fixed` type names are now deprecated aliases (marked `// Deprecated:` for IDE/linter visibility); SDK method signatures, examples, and tests use the underlying generated names. `fusionplus.CreateAuctionDetailsPlus` now takes the generated `*Preset` directly, retiring the internal `PresetClassFixed` bridge, and `fusionorder.AuctionPointClassFixed`/`GasCostConfigClassFixed` are renamed to `AuctionPoint`/`GasCostConfig` with deprecated aliases.
- Known type bugs in the upstream OpenAPI specs are now corrected at generation time by a declarative override table (`codegen/overrides.go`) instead of hand-maintained `*Fixed` copies of generated types. All 12 `*Fixed` types in `fusion`, `fusionplus`, and `tokens` are now aliases of the corrected generated types. For `fusion` and `fusionplus` the field sets are unchanged (aliases are drop-in); `fusionplus.QuoterControllerGetQuoteWithCustomPresetsParamsFixed` additionally gains `Fee` and `IsPermit2` fields the endpoint accepts.
- New `make codegen-fetch-specs` (`go run ./codegen/cmd/fetch-specs`): refreshes the local OpenAPI spec copies from the Dev Portal (requires `DEV_PORTAL_TOKEN`), normalizing whitespace only and validating each fetched spec against the codegen override table so upstream drift fails loudly at fetch time. The committed specs are now treated as untouched upstream copies — all corrections live in the pipeline.
- New spec provenance lock (`codegen/specs.lock.json`): records each spec's source URL, upstream version, content hash, and fetch time, tying the generated code to a specific upstream snapshot; CI verifies the hashes, so hand-edited spec copies fail the build. A weekly `spec-drift.yml` workflow fetches the live specs and opens an automated PR with the diff, regenerated types, and updated lock when upstream changes.
- The type-generation pipeline was rewritten from bash/jq/sed (`generate_types.sh`) into a Go tool (`codegen` package, `go run ./codegen/cmd/generate-types`, `make codegen-types`). Output is byte-identical; spec files are no longer mutated in place; the pinned oapi-codegen is invoked hermetically via `go run module@version`. No user-facing impact.

### Fixed
- **`fusionplus` quote fee was silently dropped from requests**: the `Fee *big.Int` param could not be encoded by `go-querystring`, so a set fee never reached the API. The field is now `int` (bps) and a wire-format test pins the encoding.
- `codegen/generate_types.sh` now runs correctly on Linux and with jq 1.6: BSD-only `sed -i ''` calls are replaced with a portable helper, and jq error suppression that silently deleted entire `paths`/`components.schemas` sections on jq 1.6 is replaced with explicit existence guards (no user-facing impact; the SDK code itself is unchanged by regeneration)

### Breaking Changes
- **`tokens.ProviderTokenDtoFixed` / `tokens.TokenInfoDtoFixed` optional fields are now values instead of pointers** (`DisplayedSymbol`, `Eip2612`, `IsFoT`, `LogoURI`, `Extensions`), matching the SDK-wide no-pointer policy for optional fields; both types are now aliases of the corrected generated types, and `Tags` is now `[]TagDto` on the generated types as well.
- **Type generator upgraded from deprecated `deepmap/oapi-codegen` v1.16.2 to `oapi-codegen/oapi-codegen` v2.8.0.** Incompatible deltas (machine-verified with `gorelease`): bare enum constants in `spotprices` and `traces` are now prefixed with their type name (`spotprices.USD` → `spotprices.GetPricesForRequestedTokensParamsCurrencyUSD`, `traces.CALL` → `traces.CoreCustomRootTxEventCallstackTraceFullDtoTypeCALL`); optional object-typed fields `history.TransactionDetailsDto.Meta`, `nft.Asset.Collection`, and `nft.Asset.RarityData` are values instead of pointers; the unused `web3.ApiKeyAuthScopes` constant is gone. Additions: enums gain a `Valid()` method, response types gain doc comments, and `traces` gains typed union helpers.
- **Raw generated types now carry the spec corrections previously exclusive to the `*Fixed` types**: `QuoteId` string, `Amount` string, `Fee *big.Int`, `IsPermit2` bool, string USD prices, `Points` slice, `ExclusiveResolver` string (fusion/fusionplus). See BREAKING_CHANGES.md for the full list.
- **Chain-id fields changed from `float32` to `int`**: A `float32` cannot exactly represent integers above 2²⁴, so Aurora's chain id (1313161554) was silently rounded to 1313161600 — corrupting the EIP-712 signing domain and the encoded escrow `DstChainId`, and making Aurora impossible to pass through parameter validation. All chain-id fields in the `fusionplus` and `tokens` packages are now `int`. Affected exported types include `fusionplus.QuoterControllerGetQuoteParamsFixed` (`SrcChain`, `DstChain`), `fusionplus.QuoterControllerGetQuoteWithCustomPresetsParamsFixed`, `fusionplus.GetSettlementContractParams`, `fusionplus.EscrowExtension`/`EscrowExtensionParams`/`EscrowExtraData` (`DstChainId`), `fusionplus.SignedOrderInput` (`SrcChainId`), the generated fusionplus order/quoter/relayer types, and `tokens.ProviderTokenDtoFixed`/`TokenInfoDtoFixed` (`ChainId`). Callers using untyped constants (e.g. `constants.AuroraChainId` or literal chain ids) are unaffected; callers passing `float32`-typed variables must drop the conversion. The fix is applied at the codegen layer (`codegen/generate_types.sh` now rewrites chain-id fields typed as `number` to `integer` before type generation), and the internal `CheckChainIdFloat32`/`CheckChainIdFloat32Required` validators were removed.

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


