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

## 3. Review behavior changes (compile clean, act differently)

- **Cross-chain orders with a native destination token are now signed with the
  destination chain's wrapped token** as the taker asset. v4 signed the
  source chain's wrapper address — a funds-safety bug. If you worked around it
  by passing the wrapped address explicitly, your code is unaffected.
- **`fusionplus` quote fees are now transmitted.** In v4 a configured `Fee`
  never reached the quoter (encoding bug), so quotes were priced without it.
- **`fusionplus.DecodeEscrowExtension` output is corrected**: source and
  destination safety deposits are no longer swapped, deposits are decimal
  strings, and the hashlock is 0x-prefixed hex. Decoded extensions now
  round-trip losslessly through re-encoding.
- **HTTP requests time out after 60s** (previously they could hang forever)
  and response bodies are capped at 64 MiB. Pass a context with a deadline for
  tighter control.
- Optional JSON fields generated without pointers marshal with `omitempty`
  where they previously always emitted.

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

If anything behaves differently after upgrading that is not listed here,
please open an issue — the SDK's CI machine-verifies its public API surface
against each release, and undocumented differences are treated as bugs.
