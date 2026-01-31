.PHONY: build install clean test lint fmt help

# Build variables
BINARY_NAME=openrouter
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=cmd/openrouter/main.go
VERSION=0.1.0

# Go build flags
GO_BUILD_FLAGS=-ldflags "-X main.Version=$(VERSION)"
CGO_ENABLED=0

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the openrouter binary
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p bin
	@CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "✓ Binary built: $(BINARY_PATH)"

install: build ## Install openrouter to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install $(GO_BUILD_FLAGS) $(MAIN_PATH)
	@echo "✓ Installed: $(GOPATH)/bin/$(BINARY_NAME)"

clean: ## Remove built binaries
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "✓ Cleaned"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Linting..."
	@golangci-lint run ./...

fmt: ## Format code with gofmt and goimports
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .
	@echo "✓ Code formatted"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies downloaded"

dev: ## Run in development mode with debug output
	@$(BINARY_PATH) --debug

run-help: build ## Show CLI help
	@$(BINARY_PATH) --help

run-chat: build ## Run a test chat request (requires OPENROUTER_API_KEY)
	@$(BINARY_PATH) chat "Hello, OpenRouter!"

run-list: build ## List available models (requires OPENROUTER_API_KEY)
	@$(BINARY_PATH) list

.DEFAULT_GOAL := help
