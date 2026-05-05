.PHONY: run build test clean migrate-up migrate-down

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v -cover ./...

clean:
	rm -rf bin/