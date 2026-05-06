.PHONY: setup test lint build run clean

setup:
	go mod download
	go mod tidy

test:
	go test -v ./...

lint:
	golangci-lint run ./...

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

clean:
	rm -rf bin/
	go clean
