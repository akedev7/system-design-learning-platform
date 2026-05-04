# Go Service

## Overview
A Go CLI utility service providing file operations. This service follows TDD principles and is designed to be used as a command-line tool.

## Service Details
- **Language**: Go 1.21+
- **Testing Framework**: Go built-in testing
- **Build Tool**: go build
- **Task Runner**: Makefile

## Available Targets
Run targets using `make` command:

- `make setup` - Download dependencies
- `make test` - Run tests
- `make test-verbose` - Run tests with verbose output
- `make lint` - Run linter (golangci-lint)
- `make build` - Build the binary
- `make clean` - Remove build artifacts

## Project Structure
```
go-service/
├── Claude.md              # This file
├── Makefile               # Task definitions
├── go.mod                 # Go module file
├── go.sum                 # Go checksum file
├── src/
│   └── fileops.go         # File operation functions
└── tests/
    └── fileops_test.go    # Tests for file operations
```

## Features

### File Operations (`fileops.go`)
- `ReadFile(filename string) (string, error)` - Reads content of a file
- `WriteFile(filename, content string) error` - Writes content to a file
- `FileExists(filename string) bool` - Checks if file exists
- `DeleteFile(filename string) error` - Deletes a file

## TDD Workflow
This service was built using Test-Driven Development:
1. Write failing test (Red)
2. Implement minimal code to pass (Green)
3. Refactor and lint (Refactor)

## Usage Example
```go
package main

import (
    "fmt"
    "go-service/src"
)

func main() {
    // Write a file
    src.WriteFile("test.txt", "Hello, World!")
    
    // Check if file exists
    exists := src.FileExists("test.txt")
    fmt.Println("File exists:", exists)
    
    // Read the file
    content, _ := src.ReadFile("test.txt")
    fmt.Println("Content:", content)
    
    // Delete the file
    src.DeleteFile("test.txt")
}
```

## Versioning
Check `go.mod` for module version. Follow semantic versioning.
