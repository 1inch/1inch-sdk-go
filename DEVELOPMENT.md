# SDK Developer Guide

This guide covers the workflows for developing and contributing to the Go SDK.

## Setup

- Go 1.25 or newer
- `make get` to fetch dependencies
- Optional: [Foundry](https://getfoundry.sh) (`anvil` on PATH) to run the mainnet-fork integration tests
- Optional: `prettier` (`npm install -g prettier`) to format OpenAPI spec files

## Repository layout

Each 1inch API has its own client package under `sdk-clients/<name>/`:

- `client.go`, `configuration.go`: client construction (`NewClient` with wallet + API access, `NewClientOnlyAPI` for read-only use)
- `api.go`: one method per REST endpoint
- `validation.go`: request parameter validation
- `<name>_types.gen.go`: generated types (never edit by hand)
- `<name>_types_extended.go`: manual type additions
- `examples/`: one runnable example per operation

Shared code lives in `common/` (interfaces, plus `common/fusionorder/` for types shared by `fusion` and `fusionplus`), chain IDs and contract addresses in `constants/`, and non-exported implementation in `internal/`.

## Common commands

```bash
make test              # unit tests with -race
make lint              # golangci-lint, same version CI runs
make fmt               # go fmt
make codegen-types     # regenerate types from the OpenAPI specs
```

After making changes, run the full verification sequence:

```bash
go build ./...
go vet ./...
go test ./...
make lint
```

## Type generation

API types are generated from the OpenAPI specs in `codegen/openapi/*-openapi.json` using `oapi-codegen`. To add or update a spec, place the JSON file there and run `make codegen-types` from the repository root; output lands in `sdk-clients/<pkg>/<pkg>_types.gen.go`. Generated files are overwritten on every run, so manual additions belong in the package's `*_types_extended.go` file. Format spec JSON with `prettier --write` before committing.

## Testing

The SDK has three testing tiers:

1. **Unit tests** live alongside the source and run with `make test`. All tests use the table-driven pattern: a `tests` slice of cases, `tc` as the loop variable, a `name` field per case, subtests via `t.Run(tc.name, ...)`, `require` for fatal assertions and `assert` for non-fatal ones.
2. **Mainnet-fork integration tests** (`make test-integration`, build tag `integration`) spawn `anvil` forking Ethereum mainnet and fill orders against the deployed production contracts. No real funds are involved. They run nightly in CI.
3. **Production canaries** (`make test-canary`) place real dust-sized trades through the production API on Base and Arbitrum, covering fusion, fusion plus, and aggregation across direct approvals, EIP-2612 permits, and Permit2. They skip unless `DEV_PORTAL_TOKEN`, `CANARY_WALLET_KEY`, `CANARY_BASE_RPC_URL`, and `CANARY_ARBITRUM_RPC_URL` are set, and run weekly in CI.

See `tests/integration/README.md` for the canary coverage matrix, wallet setup, and security posture.

## Changelog

Every change visible to SDK consumers adds a bullet under `## [Unreleased]` in `CHANGELOG.md`. The audience is external users of the module: cover public API changes and observable behavior, and leave out internal refactors. PR CI warns when shipped code changes without a CHANGELOG entry; changes that merge without one get auto-generated release notes from the PR title. Breaking changes also get a migration entry in `BREAKING_CHANGES.md`.

## Versioning and releases

Versions are git tags following semantic versioning. `internal/version/version.go` mirrors the most recent release for the User-Agent header; the "Release new version" workflow owns that file, and PR CI rejects manual edits to it.

To release, a maintainer dispatches the "Release new version" workflow (Actions tab) with a patch, minor, or major bump. The workflow computes the next version from the latest tag, writes it into the version constant, rolls the CHANGELOG `Unreleased` section into a dated version heading, commits, tags that commit, and publishes a GitHub release using the CHANGELOG section as its notes (auto-generated notes when the section is empty).

## Submitting changes

1. Branch from `main` and open a pull request.
2. CI must pass: `test` (unit tests + integration compile check), `lint`, and `version-and-changelog`.
3. Pull requests are squash-merged.
