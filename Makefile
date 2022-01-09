GO := go
GO_BUILD := $(GO) build
export GOBIN ?= $(shell pwd)/bin
export GO111MODULE := on

rec_wildcard = $(foreach d,$(wildcard $1*),$(call rec_wildcard,$d/,$2) $(filter $(subst *,%,$2),$d))
SRC := $(call rec_wildcard,,*.go) go.mod go.sum

MAKE2HELP := $(GOBIN)/make2help
STATIC_CHECK = $(GOBIN)/staticcheck

$(GOBIN)/%:
	@scripts/tools.sh $(notdir $@)

.DEFAULT_GOAL := help
.PHONY: fmt lint test tools clean help

## Format files
fmt:
	$(GO) fmt ./...

## Run linters
lint: $(STATIC_CHECK)
	$(GO) vet ./...
	$(STATIC_CHECK) ./...

## Run tests concurrently
test:
	$(GO) test -race ./...

## Install tools
tools:
	@scripts/tools.sh

## Clean up artifacts
clean:
	$(GO) clean

## Show help via make2help
help: $(MAKE2HELP)
	@$(MAKE2HELP)
