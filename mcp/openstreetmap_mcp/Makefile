.PHONY: help build run test lint clean

# Default target: help
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the osm-mcp binary
	GOWORK=off go build -o osm-mcp cmd/server/main.go

run: build ## Build and run the server (stdio by default)
	./osm-mcp

test: ## Run tests
	GOWORK=off go test -v ./...

lint: ## Run linter (requires golangci-lint)
	GOWORK=off golangci-lint run

clean: ## Remove build artifacts
	rm -f osm-mcp
