# Makefile for go-code

.PHONY: build test vet docs build-all clean help

# Default target
build: go-build

# Build the Go application
go-build:
	go build -o bin/go-code ./cmd/go-code

# Run tests (Go tests + Python tests if harness exists)
test:
	go test -v ./...
	@if [ -d harness ] && [ -f harness/pytest.ini -o -d harness/tests ]; then \
		cd harness && python -m pytest -v; \
	fi

# Run Go vet
vet:
	go vet ./...

# Serve documentation locally
docs:
	mkdocs serve

# Build for all platforms
build-all:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/go-code-linux-amd64 ./cmd/go-code
	GOOS=darwin GOARCH=amd64 go build -o bin/go-code-darwin-amd64 ./cmd/go-code
	GOOS=darwin GOARCH=arm64 go build -o bin/go-code-darwin-arm64 ./cmd/go-code

# Clean build artifacts
clean:
	rm -rf bin/

# Display help
help:
	@echo "Available targets:"
	@echo "  build       - Build the Go application (default)"
	@echo "  test        - Run Go and Python tests"
	@echo "  vet         - Run go vet for static analysis"
	@echo "  docs        - Serve documentation locally"
	@echo "  build-all   - Build for linux/amd64, darwin/amd64, darwin/arm64"
	@echo "  clean       - Remove bin/ directory"
	@echo "  help        - Display this help message"