# Migrating from v4 to v5

This is the step-by-step upgrade checklist. For the full per-change reference
(what changed and why), see [BREAKING_CHANGES.md](BREAKING_CHANGES.md).

**TL;DR: most v4 code compiles unchanged.** Every renamed type keeps its old
name as a deprecated alias, and the SDK's method signatures are
alias-compatible. For a typical integration the upgrade is step 1 plus a
recompile.

## 1. Update the module path

```bash
go get github.com/1inch/1inch-sdk-go/v5
find . -name '*.go' -exec sed -i 's|github.com/1inch/1inch-sdk-go/v4|github.com/1inch/1inch-sdk-go/v5|g' {} +
go mod tidy
```

## 2. Recompile and fix what the compiler flags

These are the only changes that can produce compile errors, roughly ordered by
likelihood:

| If the compiler says | Do this |
|---|---|
| cannot use `float32(...)` as `int` (chain-id fields) | Drop the `float32(...)` conversion — `constants.*ChainId` values and integer literals work as-is |
| undefined: `spotprices.USD` (or other currency consts) | Use the prefixed name: `spotprices.GetPricesForRequestedTokensParamsCurrencyUSD` |
| undefined: `traces.CALL` (or other trace-type consts) | Use the prefixed name: `traces.CoreCustomRootTxEventCallstackTraceFullDtoTypeCALL` |
| cannot use `big.NewInt(...)` as `int` (fusionplus `Fee`) | Pass the bps value as a plain `int` (e.g. `Fee: 100`). Note: in v4 the fee was silently omitted from quote requests; it is now actually sent |
| invalid operation on `tokens` fields (`DisplayedSymbol`, `Eip2612`, `IsFoT`, `LogoURI`, `Extensions`) | These are values now, not pointers: replace nil-checks with zero-value checks and drop dereferences |
| invalid operation on `history.TransactionDetailsDto.Meta`, `nft.Asset.Collection`, `nft.Asset.RarityData` | Same: values now, not pointers |
| unknown field `LimitOrderV3SubscribedApiController...` in `orderbook` wrapper struct literal | Use the new embedded field name, e.g. `GetAllOrdersParams{LimitOrdersQueryParams: q}` |
| cannot use `*PresetClassFixed` in `fusionplus.CreateAuctionDetailsPlus` | Pass the generated `*Preset` (from `GetPreset`) directly |
| undefined: `web3.ApiKeyAuthScopes` | Delete the reference; it was an unused generated artifact |
| unknown field `U` / undefined `RequestPayload.U` | Use `RequestPayload.Path` (the field was renamed from the cryptic `U`) |
| unknown field `EscExtension` on `fusionplus.Order` | Use `Order.EscrowExtension` |
| too many arguments to `traces.NewConfiguration` | It now takes a struct: `traces.NewConfiguration(traces.ConfigurationParams{ChainId: id, ApiUrl: url, ApiKey: key})` |
| cannot use `NetworkEnum` as `int` in `constants.GetWrappedToken` / `ChainToWrapper[…]` | Pass an `int` chain id (`constants.EthereumChainId`, or `int(x)`); `NetworkEnum` is deprecated |
| undefined: `fusion.Preset` (old hand-written struct), `fusionplus.FusionOrderV4`/`CrossChainOrderParams`, or the removed `orderbook` dead types (`GetCountParams`, `CountResponse`, `GetEventParams`, `GetEventsParams`, `EventResponse`, `OrderResponseExtended`, `GetActiveOrdersWithPermitParams`) | These were dead/unused and are removed; use the live equivalents (e.g. `orderbook.GetOrderCountParams`) |

Everything else in the v5 surface cleanup — the intent-based renames (e.g. `history.Item`→`HistoryEvent`, `balances` `…List`-suffix drops, `orderbook.Decode`→`DecodeExtension`, `fusion.Preset` clean names, `GetOrderByOrderHash`→`GetOrderFillsByHash`, `Whitelisted…`→`GetWhitelisted…`, `NetworkEnum`→`*ChainId`) — keeps the old name as a `// Deprecated:` alias or forwarding method, so it compiles unchanged. See [BREAKING_CHANGES.md](BREAKING_CHANGES.md) → "API Surface Cleanup (v5)" for the full per-package list.

## 3. Review behavior changes (compile clean, act differently)

These compile fine but behave differently — verify each applies to you. See
BREAKING_CHANGES.md → "Behavior Changes" for detail.

- [ ] `fusionplus` quote `Fee` is now actually transmitted — re-check your pricing.
- [ ] `fusionplus.DecodeEscrowExtension` output is corrected (safety deposits, hashlock) and round-trips losslessly.
- [ ] HTTP requests now time out at 60s and cap response bodies at 64 MiB — pass a shorter context deadline if you need one.

## 4. Optionally: move off the deprecated names

Old type names still work indefinitely, but they are marked `// Deprecated:`
(your editor shows strikethroughs; `staticcheck` flags them). The most common
renames:

| v4 name | v5 name |
|---|---|
| `fusionplus.QuoterControllerGetQuoteParamsFixed` | `fusionplus.QuoteParams` |
| `fusionplus.GetQuoteOutputFixed` | `fusionplus.Quote` |
| `fusion.QuoterControllerGetQuoteParamsFixed` | `fusion.QuoteParams` |
| `fusion.GetQuoteOutputFixed` | `fusion.Quote` |
| `tokens.ProviderTokenDtoFixed` | `tokens.ProviderTokenDto` |
| `orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams` | `orderbook.LimitOrdersQueryParams` |
| `txbroadcast.TxProcessorApiControllerBroadcastTransactionJSONRequestBody` | `txbroadcast.BroadcastPublicTransactionJSONRequestBody` |

Each package's `deprecated_names.go` lists every alias with its replacement.

## 5. Verify

```bash
go build ./... && go vet ./...
go run honnef.co/go/tools/cmd/staticcheck@latest ./... # flags remaining deprecated-name usages
```

If something behaves differently after upgrading and isn't listed here, please
open an issue.

## Appendix: two common fix patterns

The full list of affected types lives in BREAKING_CHANGES.md; these are the two
edits most callers actually make.

**Chain-id fields are now `int`** — drop `float32(...)` conversions:

```go
fusionplus.QuoteParams{SrcChain: float32(constants.ArbitrumChainId)} // before
fusionplus.QuoteParams{SrcChain: constants.ArbitrumChainId}          // after
```

**Tokens optional fields are now values** — use zero-value checks, not nil
checks (same for `history.TransactionDetailsDto.Meta`, `nft.Asset.Collection`,
`nft.Asset.RarityData`):

```go
if token.DisplayedSymbol != nil { … } // before
if token.DisplayedSymbol != "" { … }  // after
```
