# Define the Go source files
SRC_DIR := ./cmd/api
MAIN := $(SRC_DIR)/main.go
TEST_DIR := ./test

# Define the binary output directory
BIN_DIR := ./bin
BIN_NAME := api

# Configuration file
CONFIG := config.yaml

# Docker Compose file
DOCKER_COMPOSE := docker-compose.yaml

# Default target
.PHONY: all
all: build

# Run the application
.PHONY: run
run:
	@echo "Running the application..."
	go run $(MAIN) --config $(CONFIG)

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) $(MAIN)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v $(TEST_DIR)/...

# Run lint
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Install dependencies
.PHONY: install-deps
install-deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod vendor

# Update dependencies
.PHONY: update-deps
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	go mod vendor

# Run Docker Compose
.PHONY: dockerize
dockerize:
	@echo "Running Docker Compose..."
	rm -rf $(BIN_DIR)
	docker compose down --remove-orphans
	docker compose -f $(DOCKER_COMPOSE) up --build -d


# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)