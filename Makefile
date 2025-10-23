# Project info
APP_NAME := convert
CMD_DIR  := ./cmd/convert

# Go tools
GO      := go
LDFLAGS := -s -w

# Default target
.PHONY: all
all: build

## Build the binary
.PHONY: build
build:
	$(GO) build -o $(APP_NAME) $(CMD_DIR)

## Build a smaller binary (release)
.PHONY: build-release
build-release:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(APP_NAME) $(CMD_DIR)

## Run program (usage: make run DIR=images)
.PHONY: run
run: build
	./$(APP_NAME) $(DIR)

## Run with custom formats (usage: make run-io DIR=images IN=png OUT=jpg)
.PHONY: run-io
run-io: build
	./$(APP_NAME) -i=$(IN) -o=$(OUT) $(DIR)

## Run tests
.PHONY: test
test:
	$(GO) test -v ./...

## Coverage report
.PHONY: cover
cover:
	$(GO) test -cover ./...

## Format and vet
.PHONY: fmt vet
fmt:
	$(GO) fmt ./...
vet:
	$(GO) vet ./...

## Clean build artifacts
.PHONY: clean
clean:
	rm -f $(APP_NAME)

## Module tidy
.PHONY: tidy
tidy:
	$(GO) mod tidy

## Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build           - Build binary"
	@echo "  make build-release   - Build optimized binary"
	@echo "  make run DIR=images  - Run JPG->PNG conversion"
	@echo "  make run-io DIR=images IN=png OUT=jpg - Custom formats"
	@echo "  make test / cover    - Run tests / show coverage"
	@echo "  make fmt / vet       - Format / analyze code"
	@echo "  make tidy / clean    - Manage deps / cleanup"
