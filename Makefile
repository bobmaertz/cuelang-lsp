# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

COMMIT      := $(shell git rev-parse --short=10 HEAD)
BUILD_DATE  := $(shell date -u +%Y-%m-%d)
LDFLAGS     := -X github.com/bobmaertz/cuelang-lsp/pkg/version.Commit=$(COMMIT) -X github.com/bobmaertz/cuelang-lsp/pkg/version.BuildDate=$(BUILD_DATE)

# Main package path
MAIN_PATH=./cmd/lsp

# Binary name
BINARY_NAME=./bin/lsp

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags "$(LDFLAGS)" -v $(MAIN_PATH)

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...

test-coverage-html: test-coverage
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-func: test-coverage
	$(GOCMD) tool cover -func=coverage.out

test-coverage-view: test-coverage-html
	@which open > /dev/null && open coverage.html || \
	which xdg-open > /dev/null && xdg-open coverage.html || \
	echo "Please open coverage.html manually"

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

run:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags "$(LDFLAGS)" -v $(MAIN_PATH)
	./$(BINARY_NAME)

lint: 
	golangci-lint run 

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)_linux -v $(MAIN_PATH)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)_windows.exe -v $(MAIN_PATH)

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) --ldflags "$(LDFLAGS)" -o $(BINARY_NAME)_mac -v $(MAIN_PATH)

.PHONY: all build test test-coverage test-coverage-html test-coverage-func test-coverage-view clean run lint build-linux build-windows build-mac
