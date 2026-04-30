.PHONY: build run test up down migrate

build:
	go build ./...

run:
	go run ./cmd/orchestrator-api

test:
	go test ./...

up:
	docker compose -f deployments/docker-compose.yml up -d

down:
	docker compose -f deployments/docker-compose.yml down

migrate:
	sql-migrate up -config=dbconfig.yml -env=development
