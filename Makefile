SHELL := /bin/bash

.PHONY: lint test fmt

fmt:
	@gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

lint:
	@unformatted=$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*')); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt check failed for:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	@go vet ./...

test:
	@go test ./...
