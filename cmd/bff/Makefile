APP_NAME := bff
PKG      := ./...
BIN_DIR  := ../../bin
BIN      := $(BIN_DIR)/$(APP_NAME)

# Docker
IMAGE    := yourorg/$(APP_NAME)
TAG      ?= dev

# Runtime config
PORT     ?= 8080
PEOPLE_SVC_BASE_URL ?= http://localhost:8081

.PHONY: help tidy fmt test vet build run clean docker-build docker-run docker-push

help:
	@echo "Targets:"
	@echo "  tidy        - go mod tidy (repo root)"
	@echo "  fmt         - gofmt"
	@echo "  test        - go test"
	@echo "  vet         - go vet"
	@echo "  build       - build binary to $(BIN)"
	@echo "  run         - run locally"
	@echo "  docker-build- build docker image"
	@echo "  docker-run  - run docker container"
	@echo "  clean       - remove bin"

tidy:
	@cd ../.. && go mod tidy

fmt:
	@cd ../.. && gofmt -w ./cmd/bff ./internal

test:
	@cd ../.. && go test $(PKG)

vet:
	@cd ../.. && go vet $(PKG)

build:
	@mkdir -p $(BIN_DIR)
	@cd ../.. && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BIN) ./cmd/bff

run:
	@cd ../.. && PEOPLE_SVC_BASE_URL="$(PEOPLE_SVC_BASE_URL)" PORT="$(PORT)" go run ./cmd/bff

clean:
	@rm -rf $(BIN_DIR)

docker-build:
	@cd ../.. && docker build -f cmd/bff/Dockerfile -t $(IMAGE):$(TAG) .

docker-run:
	@docker run --rm -p $(PORT):8080 \
		-e PEOPLE_SVC_BASE_URL="$(PEOPLE_SVC_BASE_URL)" \
		$(IMAGE):$(TAG)

docker-push:
	@docker push $(IMAGE):$(TAG)
