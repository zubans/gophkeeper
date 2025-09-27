# Makefile for GophKeeper

.PHONY: help build test clean run-server run-client deps docker-build docker-up docker-down docker-dev

# Default target
help:
	@echo "GophKeeper - Secure Password Manager"
	@echo ""
	@echo "Available targets:"
	@echo "  build        Build client and server"
	@echo "  build-all    Build for all platforms"
	@echo "  test         Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  clean        Clean build artifacts"
	@echo "  run-server   Run the server"
	@echo "  run-client   Run the client (interactive)"
	@echo "  deps         Download dependencies"
	@echo "  lint         Run linter"
	@echo "  fmt          Format code"
	@echo "  migrate-srv  Run server migrations"
	@echo "  migrate-cli  Run client migrations"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build Build Docker images"
	@echo "  docker-up    Start all services with Docker Compose"
	@echo "  docker-down  Stop all Docker services"
	@echo "  docker-dev   Start development environment"
	@echo "  docker-logs  Show Docker logs"
	@echo "  docker-clean Clean Docker resources"

# Build client and server for current platform
build:
	@echo "Building GophKeeper..."
	go build -o bin/gophkeeper-server ./cmd/server
	go build -o bin/gophkeeper-client ./cmd/client
	@echo "Build completed!"

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@./scripts/build.sh

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf build/
	rm -f coverage.out coverage.html
	@echo "Clean completed!"

# Run server
run-server:
	@echo "Starting GophKeeper server..."
	go run ./cmd/server

# Run client (interactive)
run-client:
	@echo "Starting GophKeeper client..."
	go run ./cmd/client

# Run server migrations (requires DB flags env or defaults)
migrate-srv:
	@echo "Running server migrations..."
	go run ./cmd/server -port=0 || true

# Create/upgrade local client DB
migrate-cli:
	@echo "Initializing local client database..."
	GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go run ./cmd/client -version || true

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Install dependencies for development
install-deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Database setup (requires PostgreSQL)
setup-db:
	@echo "Setting up database..."
	@echo "Please ensure PostgreSQL is running and create a database named 'gophkeeper'"
	@echo "You can use the following command:"
	@echo "createdb gophkeeper"

# Run server with database
run-server-db: setup-db
	@echo "Starting server with database..."
	go run ./cmd/server -db-host=localhost -db-port=5432 -db-user=gophkeeper -db-password=password -db-name=gophkeeper

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting GophKeeper with Docker Compose..."
	docker-compose up -d
	@echo "Services started! Server: http://localhost:8080, pgAdmin: http://localhost:5050"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker-dev:
	@echo "Starting development environment..."
	docker-compose -f docker-compose.dev.yml up -d postgres pgadmin
	@echo "Database and pgAdmin started! pgAdmin: http://localhost:5050"
	@echo "Run 'make run-server' to start the server locally"

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f
