# Makefile for push-to-cdn
.ONESHELL:
.DEFAULT_GOAL := help

BINARY    := pushcdn
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT    ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE      ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS   := -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"
BUILD_DIR := ./dist

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[1;36m%-20s\033[0m \033[2m%s\033[0m\n", $$1, $$2}'

.PHONY: build install clean tidy test vet snapshot release-check tf-fmt tf-validate

build: ## Build the pushcdn binary into ./dist
	mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) .

install: ## Install pushcdn via `go install`
	go install $(LDFLAGS) .

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

tidy: ## Tidy go.mod and go.sum
	go mod tidy

test: ## Run tests
	go test ./...

vet: ## Run go vet
	go vet ./...

snapshot: ## Local goreleaser dry-run (no publish)
	goreleaser release --snapshot --clean

release-check: ## Validate .goreleaser.yaml
	goreleaser check

tf-fmt: ## terraform fmt -check
	terraform -chdir=terraform fmt -check

tf-validate: ## terraform validate
	terraform -chdir=terraform init -backend=false && terraform -chdir=terraform validate
