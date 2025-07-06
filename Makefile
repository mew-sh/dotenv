# Makefile for dotenv library

.PHONY: all build test test-verbose coverage benchmark fmt vet lint clean install example help

# Default target
all: fmt vet test

# Build the library
build:
	@echo "Building..."
	@go build ./...

# Build the command-line tool
build-cmd:
	@echo "Building command-line tool..."
	@go build -o bin/dotenv ./cmd/dotenv

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@go test -race -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Run golangci-lint (requires golangci-lint to be installed)
lint:
	@echo "Running linter..."
	@golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f bin/dotenv
	@rm -f coverage.out coverage.html
	@rm -rf dist/

# Install the command-line tool
install:
	@echo "Installing dotenv command..."
	@go install ./cmd/dotenv

# Run the example
example:
	@echo "Running example..."
	@cd examples/basic && go run main.go

# Install development dependencies
deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all quality checks
check: fmt vet lint test

# Build release binaries (requires goreleaser)
release-snapshot:
	@echo "Building release snapshot..."
	@goreleaser release --snapshot --rm-dist

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Run fmt, vet, and test (default)"
	@echo "  build         - Build the library"
	@echo "  build-cmd     - Build command-line tool"
	@echo "  test          - Run tests"
	@echo "  test-race     - Run tests with race detection"
	@echo "  coverage      - Run tests with coverage report"
	@echo "  benchmark     - Run benchmarks"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Run golangci-lint"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install dotenv command"
	@echo "  example       - Run example"
	@echo "  deps          - Install development dependencies"
	@echo "  check         - Run all quality checks"
	@echo "  help          - Show this help"
