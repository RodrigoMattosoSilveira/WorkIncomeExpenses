APP_NAME := people-svc
PKG      := ./...
BIN_DIR  := ../../bin
BIN      := $(BIN_DIR)/$(APP_NAME)

# Docker
IMAGE    := yourorg/$(APP_NAME)
TAG      ?= dev

# Runtime config
PORT     ?= 8081
DB_PATH  ?= /data/people.db

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
	@echo "  docker-run  - run docker container (with volume for sqlite)"
	@echo "  clean       - remove bin"

tidy:
	@cd ../.. && go mod tidy

fmt:
	@cd ../.. && gofmt -w ./cmd/people-svc ./internal

test:
	@cd ../.. && go test $(PKG)

vet:
	@cd ../.. && go vet $(PKG)

build:
	@mkdir -p $(BIN_DIR)
	@cd ../.. && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BIN) ./cmd/people-svc

run:
	@cd ../.. && PORT="$(PORT)" DB_PATH="$(DB_PATH)" go run ./cmd/people-svc

clean:
	@rm -rf $(BIN_DIR)

docker-build:
	@cd ../.. && docker build -f cmd/people-svc/Dockerfile -t $(IMAGE):$(TAG) .

docker-run:
	@docker run --rm -p $(PORT):8081 \
		-e DB_PATH="$(DB_PATH)" \
		-v people_data:/data \
		$(IMAGE):$(TAG)

docker-push:
	@docker push $(IMAGE):$(TAG)
