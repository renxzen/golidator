# Golidator Makefile
# A comprehensive set of commands for Go development workflow
.PHONY: help install-tools bench test clean fmt vet lint ci

# Default target
.DEFAULT_GOAL := help

# Variables
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

## help: Show this help message
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "Tools installed successfully"

## bench: Run benchmarks with memory profiling
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem -v ./...

## test: Run tests with race detection and coverage
test:
	@echo "Running tests..."
	go test -race -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## lint: Run golangci-lint (requires installation)
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run 'make install-tools' to install it."; \
	fi

## security: Run gosec (requires installation)
security:
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Run 'make install-tools' to install it."; \
	fi

## ci: Run CI pipeline (what would run in continuous integration)
ci: mod-verify fmt vet test lint security
	@echo "CI pipeline completed successfully!"

## clean: Clean build artifacts and coverage files
clean:
	@echo "Cleaning..."
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -f cpu.prof mem.prof
	go clean ./...

