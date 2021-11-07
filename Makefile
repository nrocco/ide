NAME = ide
DOCKER_IMAGE = nrocco/ide
DOCKER_IMAGE_VERSION = latest
BUILD_VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_COMMIT ?= $(shell git describe --always --dirty)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: help
help:
	@echo 'make build-all dist/$(NAME)-arm64-darwin dist/$(NAME)-amd64-darwin dist/$(NAME)-amd64-linux clear container push'

.PHONY: build-all
build-all: \
	dist/$(NAME)-amd64-freebsd \
	dist/$(NAME)-arm64-darwin \
	dist/$(NAME)-amd64-darwin \
	dist/$(NAME)-amd64-linux

.PHONY: dist/$(NAME)-amd64-freebsd
dist/$(NAME)-amd64-freebsd:
	mkdir -p dist/$(NAME)-amd64-freebsd
	docker image build --pull \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--target bin \
		--platform freebsd/amd64 \
		--output dist/$(NAME)-amd64-freebsd \
		.
	cp bin/* dist/$(NAME)-amd64-freebsd
	cp LICENSE README.md dist/$(NAME)-amd64-freebsd

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
	cp bin/* dist/$(NAME)-arm64-darwin
	cp LICENSE README.md dist/$(NAME)-arm64-darwin

.PHONY: dist/$(NAME)-amd64-darwin
dist/$(NAME)-amd64-darwin:
	mkdir -p dist/$(NAME)-amd64-darwin
	docker image build --pull \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--target bin \
		--platform darwin/amd64 \
		--output dist/$(NAME)-amd64-darwin \
		.
	cp bin/* dist/$(NAME)-amd64-darwin
	cp LICENSE README.md dist/$(NAME)-amd64-darwin

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
	cp bin/* dist/$(NAME)-amd64-linux
	cp LICENSE README.md dist/$(NAME)-amd64-linux

.PHONY: clear
clear:
	rm -rf dist

.PHONY: release
release:
	tar czf "dist/$(NAME)-amd64-darwin.tar.gz" -C dist/ "$(NAME)-amd64-darwin"
	tar czf "dist/$(NAME)-amd64-freebsd.tar.gz" -C dist/ "$(NAME)-amd64-freebsd"
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
