.PHONY: run build clean test test-v test-cover lint check fmt help

.DEFAULT_GOAL := help

# TODO: Set your binary name and main package path
BINARY_NAME := app
BUILD_DIR := ./bin
MAIN_PATH := ./cmd/app

# Run the application
run:
	go run $(MAIN_PATH)

# Build the binary
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run tests
test:
	go test -race ./...

# Run tests with verbose output
test-v:
	go test -v -race ./...

# Run tests with coverage
test-cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run all checks (lint + test)
check: lint test

# Show help
help:
	@echo "Available commands:"
	@echo ""
	@echo "  make run          - Run the application"
	@echo "  make build        - Build binary to $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests with race detector"
	@echo "  make test-v       - Run tests (verbose)"
	@echo "  make test-cover   - Run tests with coverage report"
	@echo "  make lint         - Run golangci-lint"
	@echo "  make fmt          - Format code"
	@echo "  make check        - Run all checks (lint + test)"
	@echo "  make help         - Show this help"
