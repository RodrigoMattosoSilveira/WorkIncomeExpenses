# Project variables
BINARY_NAME=wee
CMD_DIR=cmd
BUILD_DIR=bin

.PHONY: test   build run test tidy clean help
.PHONY: test_g build run test tidy clean help
.PHONY: test_f build run test tidy clean help
.PHONY: test_c build run test tidy clean help
.PHONY: test_v build run test tidy clean help

# Default target
all: build

## Build the application
build:
	@echo "▶ Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	GO111MODULE=on go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

## Run the application
run: build
	@echo "▶ Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

## Run all tests
test:
	@echo "▶ Running tests..., show only PASS/FAIL"
	go test -v ./... | grep -E 'PASS|FAIL' | tr -d ' '

# Run all tests & build golden files
test_g:
	@echo "▶ Running tests..., update golden files"
	go test -v ./... -update | grep -E 'PASS|FAIL' | tr -d ' '

# Run all tests, show only failed
test_f:
	@echo "▶ Running tests..., show only failed"
	go test -v ./... | grep -E 'FAIL' | tr -d ' '

## Run all tests with coverage
test_c:
	@echo "▶ Running tests with coverage ..."
	@go test ./... -cover

## Run all tests with verbose
test_v:
	@echo "▶ Running tests with verbose ..."
	@go test ./... -v

# Tidy dependencies
tidy:
	@echo "▶ Tidying go.mod..."
	@go mod tidy

## Clean build artifacts
clean:
	@echo "▶ Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

## Show help
help:
	@echo ""
	@echo "Available commands:"
	@echo "  make build   - Build the binary"
	@echo "  make run     - Build and run the binary"
	@echo "  make test    - Run tests"
	@echo "  make test_c  - Run tests with coverage"
	@echo "  make test_g  - Run tests, generate golden files"
	@echo "  make test_v  - Run tests with verbose output"
	@echo "  make test_f  - Run tests, show failures only"
	@echo "  make tidy    - Run go mod tidy"
	@echo "  make clean   - Remove build directory"
	@echo "  make help    - Show this help message"
	@echo ""
