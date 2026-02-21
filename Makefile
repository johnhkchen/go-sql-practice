# Configuration variables
BINARY_NAME := go-sql-practice
FRONTEND_DIR := frontend
DIST_DIR := $(FRONTEND_DIR)/dist
SERVER_PORT := 127.0.0.1:8090

.DEFAULT_GOAL := help
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
	rm -rf pb_data
	@echo "Clean complete"

dev:
	@echo "Starting development server..."
	./$(BINARY_NAME) serve --http="$(SERVER_PORT)"

test:
	flox activate -- go test ./...

help:
	@echo "Available targets:"
	@echo "  build       - Build frontend and backend (full pipeline)"
	@echo "  frontend    - Install deps and build Astro frontend"
	@echo "  backend     - Build Go binary (requires frontend)"
	@echo "  clean       - Remove all build artifacts"
	@echo "  dev         - Start development server"
	@echo "  test        - Run Go tests"
	@echo "  validate    - Validate build artifacts"
	@echo "  help        - Show this help message"