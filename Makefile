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

go-lint-install:
	@echo "  >  Installing golint"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.54.1

go-lint:
	@echo "  >  Checking if golint is installed..."
	@if [ ! -x "./bin/golangci-lint" ]; then \
		echo "golangci-lint not found, installing..."; \
		$(MAKE) go-lint-install; \
	fi
	@echo "  >  Running golint"
	@./bin/golangci-lint version
	@./bin/golangci-lint run --timeout=2m

codegen-types:
	@echo "Running generate_types.sh from the codegen directory..."
	@cd codegen && ./generate_types.sh
	@echo "Script execution completed."