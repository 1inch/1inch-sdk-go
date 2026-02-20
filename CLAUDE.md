# CLAUDE.md - 1inch SDK Go

**Module**: `github.com/1inch/1inch-sdk-go` | **Go**: 1.22+ | **Version**: v3.0.0

## Project Structure

```
sdk-clients/          # Per-API client packages (aggregation, fusion, fusionplus, orderbook, balances, gasprices, history, nft, portfolio, spotprices, tokens, traces, txbroadcast, web3)
common/fusionorder/   # Shared types/utilities for fusion and fusionplus
constants/            # Chain IDs, contract addresses, embedded ABIs
internal/             # Unexported utilities (http-executor, web3-provider, validate, bigint, bytesbuilder, bytesiterator, hexadecimal, keccak, etc.)
codegen/              # OpenAPI specs (codegen/openapi/) and generate_types.sh
```

## Commands

```bash
make test             # go test -race ./...
make lint             # golangci-lint v1.54.1
make fmt              # go fmt ./...
make codegen-types    # Generate types from OpenAPI specs
```

## Post-Change Verification

```bash
go build ./...                      # Compile check
go vet ./...                        # Static analysis
golangci-lint run --enable goimports # Lint + import formatting (goimports -w . to auto-fix)
go test ./...                       # Tests
```

## Architecture

### Client Pattern
Each package in `sdk-clients/` has: `client.go`, `api.go`, `configuration.go`, `validation.go`, `*_types.gen.go` (auto-generated, DO NOT EDIT), `*_types_extended.go` (manual fixes), `examples/`.

Two variants: `NewClient` (wallet + API) and `NewClientOnlyAPI` (read-only).

### Validation
Params have `Validate()` methods using generic `validate.Parameter[T]()` with typed check functions (`CheckEthereumAddressRequired`, `CheckBigIntRequired`, `CheckSlippageRequired`, `CheckPrivateKeyRequired`, `CheckChainIdIntRequired`, etc.).

### Type Generation
- Generator: `oapi-codegen/oapi-codegen@v2.5.1` (types-only mode)
- Specs: `codegen/openapi/*-openapi.json` (18 specs; `fusion`/`fusionplus` each consume 3)
- `generate_types.sh` copies specs to a staging directory, applies jq transforms + spec patches from `codegen/patches/`, runs codegen, then replaces `form:` tags with `url:` tags. Checked-in specs are never mutated.
- Spec patches (`codegen/patches/*.jq`) fix known upstream type errors (e.g., QuoteId as object→string, ExclusiveResolver as object→string). Add new patches here when upstream specs have incorrect types.
- CI enforces codegen freshness: PRs fail if generated types are out of date.
- Multi-spec packages produce: `{pkg}_orders_types.gen.go`, `{pkg}_quoter_types.gen.go`, `{pkg}_relayer_types.gen.go`

### The `*Fixed` Type Pattern (Remaining)
Some OpenAPI specs still have issues that can't be fixed with simple patches (e.g., fusion's entire GetQuoteOutput response schema differs from the actual API, fusionplus Fee needs `*big.Int`). These corrected versions live in `*_types_extended.go` with `*Fixed` suffix. **Check for a `*Fixed` variant before using a generated type.** When adding new API methods, test against the live API to verify generated types are correct. Prefer adding spec patches over new `*Fixed` types when possible.

## Fusion/FusionPlus Architecture

```
common/fusionorder/   # Base layer: Bps, Interaction, AuctionDetails, WhitelistItem, salt, presets
fusion/               # Single-chain: Extension, SettlementPostInteractionData (uses Encode() with point count)
fusionplus/           # Cross-chain: ExtensionPlus, EscrowExtension, SettlementPostInteractionData (uses EncodeWithoutPointCount())
```

**Naming**: fusionplus types use `Plus` suffix (`ExtensionPlus`, `ExtensionParamsPlus`, `CreateAuctionDetailsPlus()`). Variables: `extensionPlus`, `auctionDetailsPlus` in fusionplus; `fusionExtension`, `fusionOrder` in fusion.

**Shared code**: Functions go in `fusionorder` and are imported directly. Type aliases are acceptable for re-exporting commonly-used types in leaf packages.

## Breaking Changes Documentation

**Always** update both `CHANGELOG.md` and `BREAKING_CHANGES.md` as part of any change. Add entries to the `[Unreleased]` section as you work — do not defer this to a later step. Only document changes affecting the **public API surface** (compile errors, behavior changes, required code updates for downstream consumers). Exclude internal changes, `internal/` package changes, file renames.

## Testing

Table-driven tests required. Use `tests` as slice name, `tc` as loop variable, `t.Run(tc.name, ...)`. Use `require` for fatal assertions, `assert` for non-fatal. Both from `testify`.

## Development Rules

1. No `log.Fatalf` — always return errors
2. BigInt amounts passed as strings
3. Use `*Fixed` types when they exist; prefer adding spec patches (`codegen/patches/`) over new `*Fixed` types
4. Struct tags: `url:` not `form:`
5. Check for duplicate code before adding functions
6. If `foo.go` is deleted, check if `foo_test.go` should be too
7. Don't duplicate functions across fusion/fusionplus — put shared logic in `fusionorder`
8. Ethereum addresses are case-insensitive; use `strings.ToLower()` for comparisons
9. Codegen no longer mutates specs — safe to re-run anytime with `make codegen-types`
10. ABIs are embedded via `//go:embed` in `constants/abis.go`
11. Transaction builder auto-fetches nonce/gas if not set
12. API errors are JSON: `statusCode`, `error`, `description`
