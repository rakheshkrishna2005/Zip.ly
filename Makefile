.PHONY: build run clean test migrate

# Default target
all: build

# Build the application
build:
	@echo "Building URL Shortener..."
	go build -o bin/url-shortener ./cmd/server

# Run the application
run:
	@echo "Running URL Shortener..."
	go run ./cmd/server/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Database migration setup (requires golang-migrate)
migrate-up:
	@echo "Running database migrations..."
	migrate -path migrations -database "postgresql://postgres:prk@localhost:5432/url_shortener?sslmode=disable" up

migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "${DATABASE_URL}" down

# Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name;
	migrate create -ext sql -dir migrations -seq $$name

# Initialize the database (create database if it doesn't exist)
init-db:
	@echo "Initializing database..."
	psql -U postgres -c "CREATE DATABASE url_shortener;" || true

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go get -u github.com/gorilla/mux
	go get -u github.com/lib/pq
	go get -u github.com/joho/godotenv
	go get -u github.com/jmoiron/sqlx

# Run the application with hot reload (requires air)
dev:
	@echo "Running in development mode with hot reload..."
	air -c .air.toml