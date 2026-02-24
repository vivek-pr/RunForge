SHELL := /bin/bash

.PHONY: help tools fmt lint test tidy ci

help: ## Print available targets
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_-]+:.*## / {printf "%-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tools: ## Verify required local tools are installed
	@command -v go >/dev/null 2>&1 || (echo "missing: go" && exit 1)
	@echo "tool check passed"

fmt: ## Format Go code
	@files=$$(find . -type f -name '*.go' -not -path './vendor/*'); \
	if [ -n "$$files" ]; then gofmt -w $$files; fi

lint: ## Run golangci-lint
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 run ./...

test: ## Run unit tests
	@go test ./...

tidy: ## Tidy Go module files
	@go mod tidy

ci: ## Run local CI-equivalent checks
	@$(MAKE) test
	@$(MAKE) lint
