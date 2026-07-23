# Mainnet-fork integration tests

These tests validate SDK behavior end to end against real mainnet contracts on a local
[anvil](https://getfoundry.sh) fork. They are excluded from normal `make test` and CI by the
`integration` build tag.

## Prerequisites

- Foundry installed (`anvil` on PATH): `curl -L https://foundry.paradigm.xyz | bash && foundryup`
- A mainnet RPC endpoint. Set `FORK_URL`, or the harness probes a list of public RPCs.
  An archive-capable endpoint is more reliable for forking.

## Running

```bash
make test-integration
# or
FORK_URL=https://your-mainnet-rpc go test -tags integration -v -timeout 15m ./tests/integration/...
```

## What is covered

`TestFusionOrderPermit2Fork` validates permit2 fusion orders end to end:

1. Boots an anvil mainnet fork and deploys `SimpleSettlement`.
2. The maker approves WETH to the canonical Permit2 contract and signs a Permit2
   `PermitSingle` (AllowanceTransfer) message with the Limit Order Protocol v4 as spender.
3. A fusion order is built through `fusion.CreateFusionOrderData` with `IsPermit2: true`
   and the permit calldata attached. The first 20 bytes of the maker permit carry the
   maker asset: the protocol's `tryPermit` receives them as its token parameter.
4. A whitelisted taker fills the order via `fillOrderArgs` on the live LOP v4.
5. Assertions: WETH and USDC move between maker and taker, and the finite Permit2
   allowance granted by the embedded permit is fully consumed, proving maker funds were
   transferred through Permit2. Partial fills are covered by a dedicated subtest.
6. A negative control reproduces the pre-fix Go SDK behavior (the permit silently dropped
   from the order) and asserts the fill reverts for a maker who only approved Permit2.
7. A compact (96-byte) permit subtest pins the current router behavior: the router's
   expansion of the compact amount leaves dirty upper bits that Permit2's uint160
   calldata validation rejects, so such fills revert today. The subtest fails loudly
   if a future router deployment starts accepting compact permits.
8. A regular EIP-2612 permit subtest sells USDC through a permit built with
   `Wallet.TokenPermit`, proving the non-Permit2 permit path fills with no prior
   ERC20 approval and consumes the granted allowance exactly.

`TestPermit2TargetSemanticsFork` pins the token field semantics of the full 352-byte
maker permit: the protocol dispatches full permits to the canonical Permit2 contract by
calldata length and ignores the leading token field, so an order encoded with the
Permit2 address in place of the maker asset still fills. The SDK always encodes the
maker asset there for consistency across the SDK, API, and decoders.

## Regenerating testdata

`testdata/SimpleSettlement.json` (abi + creation bytecode) comes from the fusion-sdk repo:

```bash
cd <fusion-sdk> && forge build
jq '{abi: .abi, bytecode: .bytecode.object}' dist/contracts/SimpleSettlement.sol/SimpleSettlement.json \
  > <1inch-sdk-go>/tests/integration/testdata/SimpleSettlement.json
```

## Production canaries

The canaries place real, dust-sized trades through the production API. They are the
live counterpart to the fork suite: the fork proves the on-chain mechanics, the
canaries prove the production API accepts Go-built requests and real counterparties
execute them.

Coverage matrix:

| Product | Chain(s) | Allowance mechanisms |
|---|---|---|
| `TestProductionCanaryFusion` | Base | direct approval, EIP-2612 permit, Permit2 permit |
| `TestProductionCanaryFusionPlus` | Base and Arbitrum | direct approval |
| `TestProductionCanaryFusionPlusEip2612` | Base and Arbitrum | EIP-2612 permit |
| `TestProductionCanaryFusionPlusPermit2` | Base and Arbitrum | Permit2 permit (opt-in) |
| `TestProductionCanaryAggregation` | Arbitrum | direct approval, EIP-2612 permit, Permit2 allowance |

The Fusion+ Permit2 canary is opt-in (`CANARY_FUSION_PLUS_PERMIT2=1`): the
cross-chain order validator currently recovers Permit2 signers with the maker
asset's ERC-2612 nonce instead of the Permit2 allowance nonce, so the order is
rejected as "invalid permit signer" whenever the two nonces differ. The test
aligns the nonces on-chain first (proven end to end in production on 2026-07-23,
order 0xc78c876bea0e2f8f4129fec856af75a8e87904bc020b49dca9a481f617858e7a), and can
join the routine runs once the validator reads the Permit2 nonce.

The tests skip unless all env vars are set:

```bash
DEV_PORTAL_TOKEN=<api key> CANARY_WALLET_KEY=<private key> \
CANARY_BASE_RPC_URL=<base rpc> CANARY_ARBITRUM_RPC_URL=<arbitrum rpc> make test-canary
```

Direction alternates automatically: single-chain canaries sell whichever of
WETH/USDC the wallet holds more trades of (0.0002 WETH or 0.5 USDC per trade), and
the Fusion+ canary bridges 1.5 USDC from whichever chain holds more, so the same
funds recycle indefinitely. Fusion fills are gasless for the maker; gas is needed
for one-time approvals and for aggregation swaps.

Wallet setup and security posture:

- Fund a dedicated wallet on BOTH chains: about 2.5 USDC, 0.0005 WETH, and 0.0005
  ETH on each of Base and Arbitrum. Anyone with write access to repository
  workflows can exfiltrate Actions secrets, so the key must protect nothing beyond
  that dust.
- Signed permits are scoped to the exact trade amount and expire after 30 minutes.
  The one-time approvals (ERC20 to Permit2/router, and the aggregation Permit2
  standing allowance) are unlimited by the canonical pattern.
- The `canary.yml` workflow runs weekly and on manual dispatch only; secrets are
  never exposed to pull_request events. Without the secrets configured the jobs
  skip and stay green.
