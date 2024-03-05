NAME = ide
DOCKER_IMAGE = nrocco/$(NAME)
DOCKER_IMAGE_VERSION = latest
BUILD_VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_COMMIT ?= $(shell git describe --always --dirty)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: help
help:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

.PHONY: build-all
build-all: \
	build-arm64-darwin \
	build-amd64-linux

.PHONY: build-arm64-darwin
build-arm64-darwin: dist/$(NAME)-arm64-darwin

.PHONY: build-amd64-linux
build-amd64-linux: dist/$(NAME)-amd64-linux

.PHONY: dist/$(NAME)-arm64-darwin
dist/$(NAME)-arm64-darwin:
	mkdir -p dist/$(NAME)-arm64-darwin
	docker image build --pull \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--target bin \
		--platform darwin/arm64 \
		--output dist/$(NAME)-arm64-darwin \
		.

.PHONY: dist/$(NAME)-amd64-linux
dist/$(NAME)-amd64-linux:
	mkdir -p dist/$(NAME)-amd64-linux
	docker image build --pull \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--target bin \
		--platform linux/amd64 \
		--output dist/$(NAME)-amd64-linux \
		.

.PHONY: coverage
coverage:
	mkdir -p coverage
	go test -v -coverpkg=./... -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

.PHONY: clear
clear:
	rm -rf dist coverage

.PHONY: release
release:
	tar czf "dist/$(NAME)-arm64-darwin.tar.gz" -C dist/ "$(NAME)-arm64-darwin"
	tar czf "dist/$(NAME)-amd64-linux.tar.gz" -C dist/ "$(NAME)-amd64-linux"
	sha256sum dist/*.tar.gz > dist/checksums.txt
	tools/release-to-github.py nrocco/$(NAME) $(BUILD_VERSION) dist/checksums.txt dist/*.tar.gz

.PHONY: container
container:
	docker image build --pull \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--tag "$(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION)" \
		.

.PHONY: push
push: container
	docker image push "$(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION)"
