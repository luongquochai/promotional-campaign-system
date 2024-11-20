# Project settings
APP_NAME := promotional-campaign-system
MAIN_FILE := cmd/server/main.go
MIGRATE_FILE := migrations/migrate.go
SWAGGER_CMD := swag
BUILD_DIR := bin
OUTPUT_BINARY := $(BUILD_DIR)/$(APP_NAME)

# Flags for the Go command
GOFLAGS := 
GOTESTFLAGS := -v

# Default goal
.PHONY: all
all: build

# Run the migration
.PHONY: migrate
migrate:
	go run $(MIGRATE_FILE)

# Run the application
.PHONY: run
run:
	go run $(MAIN_FILE)

# Build the application
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(OUTPUT_BINARY) $(MAIN_FILE)
	@echo "Built $(OUTPUT_BINARY)"

# Generate Swagger documentation
.PHONY: swagger
swagger:
	$(SWAGGER_CMD) init -g $(MAIN_FILE)

# Run tests
.PHONY: test
test:
	go test $(GOFLAGS) $(GOTESTFLAGS) ./...

# Run tests with coverage
.PHONY: coverage
coverage:
	go test -cover $(GOFLAGS) $(GOTESTFLAGS) ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	@echo "Cleaned build artifacts"

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download

# Lint the code
.PHONY: lint
lint:
	golangci-lint run ./...

# Format the code
.PHONY: fmt
fmt:
	go fmt ./...

# Validate code quality (lint + format)
.PHONY: validate
validate: lint fmt
	@echo "Code validated successfully"

# Help menu
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run        - Run the application"
	@echo "  make build      - Build the application binary"
	@echo "  make swagger    - Generate Swagger documentation"
	@echo "  make test       - Run tests"
	@echo "  make coverage   - Run tests with coverage"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Install dependencies"
	@echo "  make lint       - Lint the code"
	@echo "  make fmt        - Format the code"
	@echo "  make validate   - Lint and format the code"
