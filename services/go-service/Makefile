.PHONY: setup test test-verbose lint build clean

BINARY_NAME=go-service
SRC_DIR=src
TEST_DIR=tests

setup:
	go mod download

test:
	go test ./$(TEST_DIR)/...

test-verbose:
	go test -v ./$(TEST_DIR)/...

lint:
	golangci-lint run ./...

build:
	go build -o $(BINARY_NAME) ./$(SRC_DIR)/

clean:
	rm -f $(BINARY_NAME)
	go clean
