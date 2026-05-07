.PHONY: build build-all run-basics run-websrv run-fiber run-movies-crud run-lesson clean vet fmt

# Build directory
BIN_DIR := bin

# Default database URL
DATABASE_URL ?= host=localhost user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable

# All lesson directories
LESSONS := $(notdir $(wildcard cmd/lessons/*))

build-all:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/ ./cmd/...

build-basics:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/basics ./cmd/basics/

build-movies-crud:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/movies-crud ./cmd/movies-crud-cli/

build-lesson-%:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$* ./cmd/lessons/$*/

run-basics: build-basics
	@$(BIN_DIR)/basics

run-movies-crud: build-movies-crud
	@DATABASE_URL=$(DATABASE_URL) $(BIN_DIR)/movies-crud $(filter-out $@,$(MAKECMDGOALS))

run-fiber: build-basics
	@DATABASE_URL=$(DATABASE_URL) $(BIN_DIR)/basics

clean:
	rm -rf $(BIN_DIR)

vet:
	go vet ./...

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

# List available lessons
list-lessons:
	@echo "Available lessons:"
	@for lesson in $(LESSONS); do echo "  $$lesson"; done

# Build a specific lesson
run-lesson-%: build-lesson-%
	@DATABASE_URL=$(DATABASE_URL) $(BIN_DIR)/$*
