.PHONY: help build run run-stdio run-sse run-streamable clean test fmt lint build-linux build-windows build-macos build-all

# Binary name.
BINARY_NAME=bin/d2mcp

# Default target: help
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the d2mcp binary for the current host
	go build -o $(BINARY_NAME) ./cmd

build-linux: ## Build for Linux (amd64)
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 ./cmd

build-windows: ## Build for Windows (amd64)
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe ./cmd

build-macos: ## Build for macOS (arm64 & amd64)
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-macos-arm64 ./cmd
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-macos-amd64 ./cmd

build-all: build-linux build-windows build-macos ## Build for all platforms

run: build ## Build and run the server
	./$(BINARY_NAME)

run-stdio: build ## Run with STDIO transport
	./$(BINARY_NAME) -transport=stdio

run-sse: build ## Run with SSE transport
	./$(BINARY_NAME) -transport=sse

run-streamable: build ## Run with Streamable HTTP transport
	./$(BINARY_NAME) -transport=streamable

clean: ## Remove build artifacts
	rm -rf bin/

test: ## Run tests
	go test -v ./...

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

deps: ## Install dependencies
	go mod download
	go mod tidy
