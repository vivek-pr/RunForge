SHELL := /bin/bash

.PHONY: help tools fmt lint test

help: ## Print available targets
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_-]+:.*## / {printf "%-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tools: ## Verify required local tools are installed
	@command -v awk >/dev/null 2>&1 || (echo "missing: awk" && exit 1)
	@command -v sed >/dev/null 2>&1 || (echo "missing: sed" && exit 1)
	@echo "tool check passed"

fmt: ## Run formatting step
	@echo "format check complete (no source formatters configured yet)"

lint: ## Run lint checks
	@test -f README.md
	@test -f LICENSE
	@test -f docs/local-dev.md
	@test -f developer/progress.md
	@test -f developer/decisions.md
	@test -f developer/commands.log
	@test -f developer/verification.md
	@test -f developer/todos.md
	@echo "lint passed"

test: ## Run test suite
	@echo "no tests yet"
	@true

