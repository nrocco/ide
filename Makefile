BUILD_VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_COMMIT ?= $(shell git describe --always --dirty)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: help
help:
	@echo 'make build-all dist/ide-amd64-freebsd dist/ide-amd64-darwin dist/ide-amd64-linux clear'

.PHONY: dist/ide-amd64-freebsd
dist/ide-amd64-freebsd:
	mkdir -p dist/ide-amd64-freebsd
	docker build --platform freebsd/amd64 --output dist/ide-amd64-freebsd .
	cp bin/* dist/ide-amd64-freebsd
	cp LICENSE README.md dist/ide-amd64-freebsd

.PHONY: dist/ide-amd64-darwin
dist/ide-amd64-darwin:
	mkdir -p dist/ide-amd64-darwin
	docker build --platform darwin/amd64 --output dist/ide-amd64-darwin .
	cp bin/* dist/ide-amd64-darwin
	cp LICENSE README.md dist/ide-amd64-darwin

.PHONY: dist/ide-amd64-linux
dist/ide-amd64-linux:
	mkdir -p dist/ide-amd64-linux
	docker build --platform linux/amd64 --output dist/ide-amd64-linux .
	cp bin/* dist/ide-amd64-linux
	cp LICENSE README.md dist/ide-amd64-linux

.PHONY: clear
clear:
	rm -rf dist

.PHONY: build-all
build-all: dist/ide-amd64-freebsd dist/ide-amd64-darwin dist/ide-amd64-linux

# .PHONY: release
# release: archive-all
# 	sha256sum dist/*.tar.gz > dist/checksums.txt
# 	tools/release-to-github.py $(REPO) $(BUILD_VERSION) dist/checksums.txt dist/*.tar.gz
