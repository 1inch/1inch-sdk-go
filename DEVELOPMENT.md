# Development Guide

This SDK is open for contributions from the community.

## Versioning

This library follows [semantic versioning](https://semver.org/). The current major version is **v3**. Breaking changes are documented in `BREAKING_CHANGES.md` and `CHANGELOG.md`.

## Project Structure

The SDK is organized into per-API client packages under `sdk-clients/`, each following a consistent pattern:

```
sdk-clients/{package}/
├── client.go                    # Client struct and constructors (NewClient, NewClientOnlyAPI)
├── api.go                       # API method implementations
├── configuration.go             # Configuration structs and constructors
├── validation.go                # Parameter validation (Validate() methods)
├── *_types.gen.go               # Auto-generated types from OpenAPI specs (DO NOT EDIT)
├── *_types_extended.go          # Manual type extensions and *Fixed workarounds
└── examples/                    # Usage examples per operation
```

Shared types and utilities live in:
- `common/fusionorder/` — Shared types for fusion and fusionplus packages
- `constants/` — Chain IDs, contract addresses, ABIs
- `internal/` — Internal utilities (not exported)

See `CLAUDE.md` for detailed architecture documentation.

## Key Commands

```bash
# Run all unit tests
make test

# Run linter (golangci-lint)
make lint

# Format code
make fmt

# Generate types from OpenAPI specs
make codegen-types

# Get dependencies
make get
```

## Type Generation

Types are auto-generated from OpenAPI specs using `oapi-codegen`:

1. OpenAPI specs live in `codegen/openapi/*-openapi.json`
2. Run `make codegen-types` from the repo root
3. Generated files: `sdk-clients/{package}/*_types.gen.go`

**Do not manually edit `*_types.gen.go` files** — they are overwritten by codegen.

The codegen script applies several pre-processing transforms to specs (see `codegen/generate_types.sh`). Post-processing replaces `form:` struct tags with `url:` tags for `go-querystring` compatibility.

### OpenAPI File Formatting

For consistency, OpenAPI spec files should be formatted with `prettier`:

```bash
npm install -g prettier
prettier --write codegen/openapi/*.json
```

## Testing

- All tests use the **table-driven test pattern** (see `CLAUDE.md` for template)
- Use `github.com/stretchr/testify` for assertions (`require` for fatal, `assert` for non-fatal)
- Run with `make test` or `go test -race ./...`

## Post-Change Verification

After making changes, run these checks:

```bash
go build ./...           # Catch compile errors
go vet ./...             # Static analysis
golangci-lint run        # Linter
go test ./...            # All tests
```

## CI/CD

- **PR Validation** (`.github/workflows/pr.yml`): Runs tests + golangci-lint on PRs
- **Release** (`.github/workflows/release.yml`): Manual dispatch for versioned releases
