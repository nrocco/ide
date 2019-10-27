BIN ?= ide
PKG ?= github.com/nrocco/ide
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git describe --always --dirty)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: build
build:
	mkdir -p build
	go build -v -o build/$(BIN)-$(GOOS)-$(GOARCH) -ldflags "-X ${PKG}/cmd.version=${VERSION} -X ${PKG}/cmd.commit=${COMMIT} -X ${PKG}/cmd.buildDate=${DATE}"

.PHONY: clear
clear:
	rm -rf build

# server/server.pb.go: server/server.proto
# 	protoc -I server/ server/server.proto --go_out=plugins=grpc:server

# .PHONY: lint
# lint: server/server.pb.go
# 	golint -set_exit_status ${PKG_LIST}

# .PHONY: vet
# vet: server/server.pb.go
# 	go vet -v ./...

# .PHONY: test
# test: server/server.pb.go
# 	go test -v -short ./...

# .PHONY: releases
# releases: build-all
# 	mkdir -p "build/${BIN}-${VERSION}"
# 	cp bin/rgit "build/${BIN}-${VERSION}/rgit"
# 	build/${BIN}-${GOOS}-${GOARCH} completion > "build/${BIN}-${VERSION}/completion.zsh"
# 	mv build/${BIN}-linux-amd64 "build/${BIN}-${VERSION}/${BIN}"
# 	tar czf "build/${BIN}-${VERSION}-linux-amd64.tar.gz" -C build/ "${BIN}-${VERSION}"
# 	mv build/${BIN}-darwin-amd64 "build/${BIN}-${VERSION}/${BIN}"
# 	tar czf "build/${BIN}-${VERSION}-darwin-amd64.tar.gz" -C build/ "${BIN}-${VERSION}"
# 	rm -rf "build/${BIN}-${VERSION}"
