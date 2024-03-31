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

fmt:
	@echo "  >  Running go fmt"
	GOBIN=$(GOBIN) go fmt ./...

lint: go-lint

go-lint-install:
	@echo "  >  Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.1

go-lint:
	@echo "  >  Running golint "
	bin/golangci-lint version
	bin/golangci-lint run --timeout=2m