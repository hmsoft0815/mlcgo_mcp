.PHONY: help build run test lint clean

# Default target: help
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the osm-mcp binary for the current host
	GOWORK=off go build -o osm-mcp cmd/server/main.go

build-linux: ## Build for Linux (amd64)
	GOWORK=off GOOS=linux GOARCH=amd64 go build -o bin/osm-mcp-linux-amd64 cmd/server/main.go

build-windows: ## Build for Windows (amd64)
	GOWORK=off GOOS=windows GOARCH=amd64 go build -o bin/osm-mcp-windows-amd64.exe cmd/server/main.go

build-macos: ## Build for macOS (arm64 & amd64)
	GOWORK=off GOOS=darwin GOARCH=arm64 go build -o bin/osm-mcp-macos-arm64 cmd/server/main.go
	GOWORK=off GOOS=darwin GOARCH=amd64 go build -o bin/osm-mcp-macos-amd64 cmd/server/main.go

build-all: build-linux build-windows build-macos ## Build for all platforms

run: build ## Build and run the server (stdio by default)
	./osm-mcp

test: ## Run tests
	GOWORK=off go test -v ./...

lint: ## Run linter (requires golangci-lint)
	GOWORK=off golangci-lint run

clean: ## Remove build artifacts
	rm -f osm-mcp
	rm -rf bin/
