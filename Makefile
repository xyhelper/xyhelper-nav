ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "xyhelper-nav"
DOCKER_NAME = "xyhelper/nav"

include ./hack/hack-cli.mk
include ./hack/hack.mk

# ---------------------------------------------------------------------------
# Project shortcuts
# ---------------------------------------------------------------------------

# Start the development server (auto reloads on file changes).
.PHONY: run
run: cli.install
	@gf run main.go

# Run the unit tests.
.PHONY: test
test:
	@go test ./... -count=1

# Format the Go sources.
.PHONY: fmt
fmt:
	@gofmt -s -w .

# Tidy the module dependencies.
.PHONY: tidy
tidy:
	@go mod tidy

# Build a local binary into ./bin without the GoFrame CLI.
.PHONY: build-local
build-local:
	@go build -trimpath -o bin/nav .

# Build the container image from the local Dockerfile.
.PHONY: docker
docker:
	@docker build -t $(DOCKER_NAME):latest .

# Start the whole stack with docker compose.
.PHONY: up-docker
up-docker:
	@docker compose up -d --build

# Stop the docker compose stack.
.PHONY: down-docker
down-docker:
	@docker compose down
