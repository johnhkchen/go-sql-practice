# Configuration variables
BINARY_NAME := go-sql-practice
FRONTEND_DIR := frontend
DIST_DIR := $(FRONTEND_DIR)/dist
SERVER_PORT := 127.0.0.1:8090

.PHONY: build frontend backend clean dev test help validate-build

build: frontend backend

frontend:
	@echo "Building frontend..."
	cd frontend && npm ci && npm run build

backend:
	@echo "Building backend..."
	flox activate -- go build -o $(BINARY_NAME)

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(DIST_DIR)
	rm -f $(BINARY_NAME)

dev:
	@echo "Starting development server..."
	./$(BINARY_NAME) serve --http="$(SERVER_PORT)"

test:
	flox activate -- go test ./...