.PHONY: build build-all test clean run-server run-client deps lint fmt

DB_HOST=localhost
DB_PORT=5432
DB_USER=gophkeeper
DB_PASSWORD=password
DB_NAME=gophkeeper

build:
	go build -o bin/gophkeeper-server ./cmd/server
	go build -o bin/gophkeeper-client ./cmd/client

build-all:
	go build -o bin/gophkeeper-server ./cmd/server
	GOOS=linux GOARCH=amd64 go build -o bin/gophkeeper-client-linux-amd64 ./cmd/client
	GOOS=linux GOARCH=arm64 go build -o bin/gophkeeper-client-linux-arm64 ./cmd/client
	GOOS=windows GOARCH=amd64 go build -o bin/gophkeeper-client-windows-amd64.exe ./cmd/client
	GOOS=windows GOARCH=arm64 go build -o bin/gophkeeper-client-windows-arm64.exe ./cmd/client
	GOOS=darwin GOARCH=amd64 go build -o bin/gophkeeper-client-darwin-amd64 ./cmd/client
	GOOS=darwin GOARCH=arm64 go build -o bin/gophkeeper-client-darwin-arm64 ./cmd/client

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

run-server:
	go run ./cmd/server -db-host=$(DB_HOST) -db-port=$(DB_PORT) -db-user=$(DB_USER) -db-password=$(DB_PASSWORD) -db-name=$(DB_NAME)

run-client:
	go run ./cmd/client

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

fmt:
	go fmt ./...