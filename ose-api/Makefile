.PHONY: help proto deps test build run docker clean

# Version information (for build metadata)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go build flags
LDFLAGS := -X main.Version=$(VERSION) \
           -X main.GitCommit=$(COMMIT) \
           -X main.BuildTime=$(BUILD_TIME)

help: ## Show this help message
	@echo "OSE API Gateway - Make targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## Download Go dependencies
	go mod download
	go mod tidy

proto: ## Generate Go code from protobuf definitions
	@echo "Checking protobuf dependencies..."
	@which protoc > /dev/null || (echo "ERROR: protoc not found. See PROTO_GENERATION.md for installation instructions." && exit 1)
	@echo "Installing Go protobuf plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	@echo "Generating protobuf code..."
	@mkdir -p api/proto/advisory/v1
	@protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       --proto_path=. \
	       proto/advisory/v1/advisory.proto
	@echo "✓ Protobuf code generated successfully"
	@echo "  Output: api/proto/advisory/v1/advisory.pb.go"
	@echo "  Output: api/proto/advisory/v1/advisory_grpc.pb.go"

test: ## Run tests with coverage
	@echo "Running tests..."
	go test ./... -v -cover -coverprofile=coverage.out
	@echo "✓ Tests complete"

test-verbose: ## Run tests with verbose output
	go test ./... -v -count=1

build: deps ## Build the advisory server binary
	@echo "Building advisory-server..."
	go build -ldflags "$(LDFLAGS)" -o bin/advisory-server cmd/advisory-server/main.go
	@echo "✓ Built: bin/advisory-server"
	@echo "  Version: $(VERSION)"
	@echo "  Commit:  $(COMMIT)"

run: ## Run the server locally
	@echo "Starting OSE Advisory Server..."
	go run cmd/advisory-server/main.go

docker: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t ose-advisory:$(VERSION) -f deployments/docker/Dockerfile .
	@echo "✓ Docker image built: ose-advisory:$(VERSION)"

docker-run: docker ## Run server in Docker
	docker run --rm -p 50051:50051 -p 8080:8080 \
		-e OSE_LOG_LEVEL=debug \
		ose-advisory:$(VERSION)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out
	go clean
	@echo "✓ Clean complete"

# Development helpers
dev: ## Run server with hot reload (requires entr)
	find . -name '*.go' | entr -r make run

lint: ## Run linters
	golangci-lint run ./...

fmt: ## Format code
	go fmt ./...
	goimports -w .

# Quick checks before commit
check: fmt lint test ## Run fmt, lint, and test

.DEFAULT_GOAL := help
