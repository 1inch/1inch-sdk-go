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

## Regenerating testdata

`testdata/SimpleSettlement.json` (abi + creation bytecode) comes from the fusion-sdk repo:

```bash
cd <fusion-sdk> && forge build
jq '{abi: .abi, bytecode: .bytecode.object}' dist/contracts/SimpleSettlement.sol/SimpleSettlement.json \
  > <1inch-sdk-go>/tests/integration/testdata/SimpleSettlement.json
```
