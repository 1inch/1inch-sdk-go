#! /usr/bin/make -f

MAKEFLAGS += --silent

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOPKG := $(.)
# A valid GOPATH is required to use the `go get` command.
# If $GOPATH is not specified, $HOME/go will be used by default
GOPATH := $(if $(GOPATH),$(GOPATH),~/go)

get:
	@echo "  >  Checking if there are any missing dependencies..."
	GOBIN=$(GOBIN) go get ./... $(get)

test:
	@echo "  >  Running unit tests"
	GOBIN=$(GOBIN) go test -race ./...

# Runs mainnet-fork integration tests. Requires anvil (foundry) on PATH and,
# optionally, FORK_URL pointing to a mainnet RPC (falls back to public RPCs).
test-integration:
	@echo "  >  Running mainnet-fork integration tests"
	GOBIN=$(GOBIN) go test -tags integration -v -timeout 15m ./tests/integration/...

# Places real dust-sized trades on Base and Arbitrum through the production API,
# covering fusion, fusion plus, and aggregation across every allowance mechanism.
# Requires DEV_PORTAL_TOKEN, CANARY_WALLET_KEY, CANARY_BASE_RPC_URL, and
# CANARY_ARBITRUM_RPC_URL.
test-canary:
	@echo "  >  Running production canary trades"
	GOBIN=$(GOBIN) go test -tags integration -v -timeout 60m -run TestProductionCanary ./tests/integration/

fmt:
	@echo "  >  Running go fmt"
	GOBIN=$(GOBIN) go fmt ./...

lint: go-lint

# Pinned to the version CI uses (.github/workflows/pr.yml) so local lint
# results match pull request checks.
GOLANGCI_LINT_VERSION := v2.10.1

go-lint-install:
	@echo "  >  Installing golangci-lint $(GOLANGCI_LINT_VERSION)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin $(GOLANGCI_LINT_VERSION)

go-lint:
	@echo "  >  Checking if golangci-lint is installed..."
	@if [ ! -x "./bin/golangci-lint" ] || ! ./bin/golangci-lint version 2>/dev/null | grep -q "$(GOLANGCI_LINT_VERSION:v%=%)"; then \
		echo "golangci-lint $(GOLANGCI_LINT_VERSION) not found, installing..."; \
		$(MAKE) go-lint-install; \
	fi
	@echo "  >  Running golangci-lint"
	@./bin/golangci-lint version
	@./bin/golangci-lint run --timeout=3m

codegen-types:
	@echo "  >  Generating types from OpenAPI specs..."
	go run ./codegen/cmd/generate-types
	@echo "Script execution completed."