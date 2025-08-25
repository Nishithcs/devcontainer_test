# Application Settings
APP_NAME = go-backend-api
API_MAIN_PATH = cmd/api/main.go
CONSUMER_MAIN_PATH = cmd/consumer/main.go
API_BINARY_NAME = $(APP_NAME)-server
CONSUMER_BINARY_NAME = $(APP_NAME)-consumer

# Go Settings
GO = go
GOFMT = gofmt
GOTEST = go test
GOVET = go vet
GOLINT = golangci-lint
AIR = air

# Docker Settings
DOCKER = docker
DOCKER_COMPOSE = docker-compose
DOCKER_IMAGE = $(APP_NAME)
DOCKER_TAG = latest

# Database Settings
MIGRATION_DIR = internal/db/migrations

# Color Settings
COLOR_RESET = \033[0m
COLOR_GREEN = \033[32m
COLOR_YELLOW = \033[33m
COLOR_BLUE = \033[34m

.PHONY: all build build-api build-consumer clean test coverage lint fmt vet run docker-build docker-run help migrate seed dev deps

# Default target
all: lint test build

# Show help
help:
	@echo "$(COLOR_BLUE)Available commands:$(COLOR_RESET)"
	@echo "$(COLOR_GREEN)make build$(COLOR_RESET)            - Build the application"
	@echo "$(COLOR_GREEN)make build-api$(COLOR_RESET)        - Build the API application"
	@echo "$(COLOR_GREEN)make build-consumer$(COLOR_RESET)   - Build the Consumer application"
	@echo "$(COLOR_GREEN)make clean$(COLOR_RESET)            - Clean build artifacts"
	@echo "$(COLOR_GREEN)make test$(COLOR_RESET)             - Run tests"
	@echo "$(COLOR_GREEN)make coverage$(COLOR_RESET)         - Run tests with coverage"
	@echo "$(COLOR_GREEN)make lint$(COLOR_RESET)             - Run linter"
	@echo "$(COLOR_GREEN)make fmt$(COLOR_RESET)              - Run gofmt"
	@echo "$(COLOR_GREEN)make vet$(COLOR_RESET)              - Run govet"
	@echo "$(COLOR_GREEN)make run$(COLOR_RESET)              - Run the API application"
	@echo "$(COLOR_GREEN)make run-consumer$(COLOR_RESET)     - Run the Consumer application"
	@echo "$(COLOR_GREEN)make dev$(COLOR_RESET)              - Run the API in development mode"
	@echo "$(COLOR_GREEN)make dev-consumer$(COLOR_RESET)     - Run the Consumer in development mode"
	@echo "$(COLOR_GREEN)make dev-env$(COLOR_RESET)          - Start development environment"
	@echo "$(COLOR_GREEN)make dev-env-build$(COLOR_RESET)    - Start development environment with build"
	@echo "$(COLOR_GREEN)make dev-env-down$(COLOR_RESET)     - Stop development environment"
	@echo "$(COLOR_GREEN)make docker-build$(COLOR_RESET)     - Build Docker image"
	@echo "$(COLOR_GREEN)make docker-run$(COLOR_RESET)       - Run Docker container"
	@echo "$(COLOR_GREEN)make deps$(COLOR_RESET)             - Install dependencies"
	@echo "$(COLOR_GREEN)make migrate$(COLOR_RESET)          - Run database migrations"
	@echo "$(COLOR_GREEN)make seed$(COLOR_RESET)             - Seed database with initial data"

# Build the application
build: build-api build-consumer

# Build API application
build-api:
	@echo "$(COLOR_BLUE)Building API application...$(COLOR_RESET)"
	$(GO) build -o bin/$(API_BINARY_NAME) $(API_MAIN_PATH)

# Build Consumer application
build-consumer:
	@echo "$(COLOR_BLUE)Building Consumer application...$(COLOR_RESET)"
	$(GO) build -o bin/$(CONSUMER_BINARY_NAME) $(CONSUMER_MAIN_PATH)

# Clean build artifacts
clean:
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	rm -rf bin/
	$(GO) clean

# Run tests
test:
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_YELLOW)Coverage report generated: coverage.html$(COLOR_RESET)"

# Run linter
lint:
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	$(GOLINT) run

# Format code
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	$(GOFMT) -w .

# Run go vet
vet:
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	$(GOVET) ./...

# Run the API application
run: build-api
	@echo "$(COLOR_BLUE)Running API application...$(COLOR_RESET)"
	./bin/$(API_BINARY_NAME)

# Run the Consumer application
run-consumer: build-consumer
	@echo "$(COLOR_BLUE)Running Consumer application...$(COLOR_RESET)"
	./bin/$(CONSUMER_BINARY_NAME)

# Run API in development mode and hot reload
dev:
	@echo "$(COLOR_BLUE)Running API in development mode...$(COLOR_RESET)"
	$(AIR) -c .air.toml

# Run Consumer in development mode
dev-consumer:
	@echo "$(COLOR_BLUE)Running Consumer in development mode...$(COLOR_RESET)"
	$(GO) run $(CONSUMER_MAIN_PATH)

# Build Docker image
docker-build:
	@echo "$(COLOR_BLUE)Building Docker image...$(COLOR_RESET)"
	$(DOCKER) build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run Docker container
docker-run: docker-build
	@echo "$(COLOR_BLUE)Running Docker container...$(COLOR_RESET)"
	$(DOCKER) run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Start development environment
dev-env:
	@echo "$(COLOR_BLUE)Starting development environment...$(COLOR_RESET)"
	$(DOCKER_COMPOSE) up -d

# Start development environment and build
dev-env-build:
	@echo "$(COLOR_BLUE)Starting development environment and building...$(COLOR_RESET)"
	$(DOCKER_COMPOSE) up -d --build

# Stop development environment
dev-env-down:
	@echo "$(COLOR_BLUE)Stopping development environment...$(COLOR_RESET)"
	$(DOCKER_COMPOSE) down

# Install development dependencies
deps:
	@echo "$(COLOR_BLUE)Installing dependencies...$(COLOR_RESET)"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/cosmtrek/air@latest
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	$(GO) mod tidy
	$(GO) mod vendor

# Database commands
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $$name

# Seed database
seed:
	@echo "$(COLOR_BLUE)Seeding database...$(COLOR_RESET)"
	$(GO) run internal/db/seeder/main.go
