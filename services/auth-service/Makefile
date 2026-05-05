.PHONY: run build test clean migrate-up migrate-down dev-db-up dev-db-down

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v -cover ./...

clean:
	rm -rf bin/

dev-db-up:
	docker compose -f docker-compose.dev.yml up -d

dev-db-down:
	docker compose -f docker-compose.dev.yml down