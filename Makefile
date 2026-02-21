.PHONY: build frontend backend clean dev test

build: frontend backend

frontend:
	@echo "Building frontend..."
	cd frontend && npm ci && npm run build

backend:
	@echo "Building backend..."
	go build -o go-sql-practice

clean:
	@echo "Cleaning build artifacts..."
	rm -rf frontend/dist
	rm -f go-sql-practice

dev:
	@echo "Starting development server..."
	./go-sql-practice serve --http="127.0.0.1:8090"

test:
	go test ./...