include .env
export

# template migrate Variables
OLD_IMPORT=github.com/gauravst/real-time-chat
OLD_CMD_DIR=cmd/real-time-chat

# Variables
BINARY_NAME=go-api-template
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')
MIGRATE_PATH=./migrations
DB_URL=$(DATABASE_URI)
APP_NAME=go-api-template
DOCKER_IMAGE_NAME=go-api-template
DOCKER_TAG=latest
PORT=8080

# Default target
all: build

# Build the application
build:
	@echo "Building the application..."
	go build -o bin/$(BINARY_NAME) cmd/go-api-template/main.go

# Run the application
run:
	@echo "Running the application..."
	go run cmd/go-api-template/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Run all up migrations
migrate-up:
	migrate -path $(MIGRATE_PATH) -database $(DB_URL) up

# Run all down migrations
migrate-down:
	migrate -path $(MIGRATE_PATH) -database $(DB_URL) down

## Build the Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

## Run the Docker container
docker-run:
	docker run -p $(PORT):$(PORT) $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

## Stop the Docker container
docker-stop:
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE_NAME):$(DOCKER_TAG))

## Remove the Docker container
docker-rm:
	docker rm $$(docker ps -a -q --filter ancestor=$(DOCKER_IMAGE_NAME):$(DOCKER_TAG))

## Remove the Docker image
docker-rmi:
	docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

## Clean up Docker resources (stop, remove container, and remove image)
docker-clean: docker-stop docker-rm docker-rmi

# Combined commands

## Build and run the Go application
all: build run

## Build and run the Docker container
docker-all: docker-build docker-run

# Run all checks (format, test, build)
check: fmt test build

# template migrate 
setup:
	@echo "Enter your GitHub username :-"
	@read USERNAME; \
		echo "Enter your project name :-"; \
	read PROJECT; \
	LOWER_USERNAME=$$(echo $$USERNAME | tr '[:upper:]' '[:lower:]'); \
	LOWER_PROJECT=$$(echo $$PROJECT | tr '[:upper:]' '[:lower:]'); \
	NEW_IMPORT=github.com/$$LOWER_USERNAME/$$LOWER_PROJECT; \
	NEW_CMD_DIR=cmd/$$LOWER_PROJECT; \
	echo "Replacing imports..."; \
	find . -type f -name "*.go" -exec sed -i 's|$(OLD_IMPORT)|'$$NEW_IMPORT'|g' {} +; \
	find . -type f -name "go.mod" -exec sed -i 's|$(OLD_IMPORT)|'$$NEW_IMPORT'|g' {} +; \
	find . -type f -name "go.sum" -exec sed -i 's|$(OLD_IMPORT)|'$$NEW_IMPORT'|g' {} +; \
	echo "Renaming main.go..."; \
	mkdir -p $$NEW_CMD_DIR; \
	mv $(OLD_CMD_DIR)/main.go $$NEW_CMD_DIR/main.go; \
	rm -rf $(OLD_CMD_DIR); \
	echo "Updating Makefile..."; \
	sed -i 's|OLD_IMPORT=$(OLD_IMPORT)|OLD_IMPORT='$$NEW_IMPORT'|g' makefile; \
	sed -i 's|OLD_CMD_DIR=$(OLD_CMD_DIR)|OLD_CMD_DIR='$$NEW_CMD_DIR'|g' makefile; \
	echo "Setup completed!"

# Help (list all targets)
help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  clean    - Clean build artifacts"
	@echo "  deps     - Install dependencies"
	@echo "  check    - Run all checks (format, test, build)"
	@echo "  help     - Show this help message"
	@echo "  docker-build  - Build the Docker image"
	@echo "  docker-run    - Run the Docker container"
	@echo "  docker-stop   - Stop the Docker container"
	@echo "  docker-rm     - Remove the Docker container"
	@echo "  docker-rmi    - Remove the Docker image"
	@echo "  docker-clean  - Clean up Docker resources (stop, remove container, and remove image)"
	@echo "  docker-all    - Build and run the Docker container"
	@echo "  setup     - To Setup Project"
