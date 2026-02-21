# Configuration variables
BINARY_NAME := go-sql-practice
FRONTEND_DIR := frontend
DIST_DIR := $(FRONTEND_DIR)/dist
SERVER_PORT := 127.0.0.1:8090

.DEFAULT_GOAL := help
.PHONY: build frontend backend clean dev test help validate-build

build: frontend backend validate-build

frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm ci && npm run build
	@if [ ! -d "$(DIST_DIR)/client" ]; then \
		echo "Error: Frontend build failed - client directory not created"; \
		exit 1; \
	fi
	@echo "Frontend build complete"

backend: frontend
	@echo "Building backend..."
	@if [ ! -d "$(DIST_DIR)" ]; then \
		echo "Error: Frontend not built. Run 'make frontend' first."; \
		exit 1; \
	fi
	flox activate -- go build -o $(BINARY_NAME)
	@echo "Backend build complete"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(DIST_DIR)
	rm -f $(BINARY_NAME)
	rm -rf pb_data
	@echo "Clean complete"

dev: build
	@echo "Starting development server..."
	@echo "Server will be available at http://$(SERVER_PORT)"
	@echo "Press Ctrl+C to stop"
	./$(BINARY_NAME) serve --http="$(SERVER_PORT)"

test:
	flox activate -- go test ./...

validate-build:
	@echo "Validating build artifacts..."
	@if [ ! -f "$(BINARY_NAME)" ]; then \
		echo "Error: Binary not created"; exit 1; \
	fi
	@if [ ! -d "$(DIST_DIR)" ]; then \
		echo "Error: Frontend assets missing"; exit 1; \
	fi
	@echo "Build validation successful"

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