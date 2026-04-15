.PHONY: build run test clean docker-up docker-down

APP_NAME=project-api

build:
	@echo "==> Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

run: build
	@echo "==> Running $(APP_NAME)..."
	@./bin/$(APP_NAME)

test:
	@echo "==> Running tests..."
	@go test -v -race ./...

vet:
	@echo "==> Running go vet..."
	@go vet ./...

lint: vet
	@echo "==> Lint passed"

clean:
	@echo "==> Cleaning up..."
	@rm -rf bin/

docker-up:
	@echo "==> Starting infrastructure..."
	@docker-compose --env-file .env -f deployments/docker-compose.yml up -d

docker-down:
	@echo "==> Stopping infrastructure..."
	@docker-compose --env-file .env -f deployments/docker-compose.yml down
