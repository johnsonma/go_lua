# Go Lua Implementation Makefile

# Build variables
BINARY_NAME=go_lua
BUILD_DIR=bin
SOURCE_DIR=./cmd/go_lua

# Go variables
GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

# Build for development (with debug info)
dev: build

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Test the application
test:
	@echo "Running tests..."
	$(GOTEST) ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GO) vet ./...

# Lint and format
lint: fmt vet

# Build and run
br: build run

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  dev      - Build for development"
	@echo "  run      - Build and run the application"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  deps     - Install dependencies"
	@echo "  fmt      - Format code"
	@echo "  vet      - Vet code"
	@echo "  lint     - Format and vet code"
	@echo "  br       - Build and run"
	@echo "  help     - Show this help message"

.PHONY: all build dev run test clean deps fmt vet lint br help
