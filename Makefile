SHELL := /bin/bash

BIN      	= $(CURDIR)/bin

GOPATH		= $(HOME)/go
GOBIN		= $(GOPATH)/bin
GO			?= GOGC=off $(shell which go)

PATH := $(BIN):$(GOBIN):$(PATH)

# Printing
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

# Tools
$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(@F)…)
	$Q GOBIN=$(BIN) $(GO) install $(shell $(GO) list -e -tags=tools -f '{{ join .Imports "\n" }}' ./tools | grep $(@F))

# golangci-lint is recommended to be installed via the install script instead of go get
$(BIN)/golangci-lint: | $(BIN) ; $(info $(M) installing  golangci-lint…) @
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.2.2

GOLANGCI_LINT = $(BIN)/golangci-lint

# Targets
.PHONY: lint
lint: | $(GOLANGCI_LINT) ; $(info $(M) running golint…) @ ## Run the project linters
	$Q $(GOLANGCI_LINT) run --max-issues-per-linter 10

.PHONY: test
test: ## Run all tests
	$Q $(GO) test ./... -timeout 20s

.PHONY: generate
generate: | ; $(info $(M) running go generate…) @ ## Run go generate
	$Q $(GO) generate ./...
	$Q $(MAKE) fmt

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt $(PKGS)

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@find . -name '*_gen.go' -exec rm -r {} \;

.PHONY: build
build: test lint

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
