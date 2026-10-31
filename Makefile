NAME = ide
DOCKER_IMAGE = registry.gitlab.com/nrocco/$(NAME)
DOCKER_IMAGE_VERSION = latest
DOCKER_IMAGE_PLATFORM = darwin/arm64,linux/amd64,linux/arm64
DOCKER_IMAGE_BUILD_ARGS =
BUILD_VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: help
help:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

.PHONY: lint
lint:
	golint -set_exit_status ./...
	go vet -v ./...
	staticcheck ./...
	govulncheck ./...
	deadcode ./...
	gosec ./...
	go test -v -cover -short ./...

.PHONY: coverage
coverage:
	mkdir -p coverage
	go test -v -coverpkg=./... -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

.PHONY: build-all
build-all: \
	build-amd64-linux \
	build-arm64-darwin \
	build-arm64-linux

.PHONY: build-arm64-darwin
build-arm64-darwin: dist/$(NAME)-arm64-darwin

.PHONY: build-amd64-linux
build-amd64-linux: dist/$(NAME)-amd64-linux

.PHONY: build-amd64-linux
build-arm64-linux: dist/$(NAME)-arm64-linux

.PHONY: dist/$(NAME)-arm64-darwin
dist/$(NAME)-arm64-darwin:
	mkdir -p dist/$(NAME)-arm64-darwin
	$(MAKE) image DOCKER_IMAGE_PLATFORM=darwin/arm64 DOCKER_IMAGE_BUILD_ARGS="--target bin --output dist/$(NAME)-arm64-darwin"

.PHONY: dist/$(NAME)-arm64-linux
dist/$(NAME)-arm64-linux:
	mkdir -p dist/$(NAME)-arm64-linux
	$(MAKE) image DOCKER_IMAGE_PLATFORM=linux/arm64 DOCKER_IMAGE_BUILD_ARGS="--target bin --output dist/$(NAME)-arm64-linux"

.PHONY: dist/$(NAME)-amd64-linux
dist/$(NAME)-amd64-linux:
	mkdir -p dist/$(NAME)-amd64-linux
	$(MAKE) image DOCKER_IMAGE_PLATFORM=linux/amd64 DOCKER_IMAGE_BUILD_ARGS="--target bin --output dist/$(NAME)-amd64-linux"

.PHONY: clear
clear:
	rm -rf dist coverage

.PHONY: image
image:
	docker buildx build $(DOCKER_IMAGE_BUILD_ARGS) \
		--platform $(DOCKER_IMAGE_PLATFORM) \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--tag "$(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION)" \
		.

.PHONY: clear push
push:
	$(MAKE) image DOCKER_IMAGE_BUILD_ARGS='--push'

.PHONY: ci
ci: clear lint coverage
