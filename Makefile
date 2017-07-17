BIN := ide
PKG := github.com/nrocco/ide
VERSION := $(shell git describe --tags --always --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v ${PKG}/vendor/)
GO_FILES := $(shell find * -type d -name vendor -prune -or -name '*.go' -type f | grep -v vendor)

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
LDFLAGS = "-d -s -w -X ${PKG}/cmd.Version=${VERSION}"
BUILD_ARGS = -a -tags netgo -installsuffix netgo -ldflags $(LDFLAGS)

PREFIX = /usr/local

.DEFAULT_GOAL: build

build/$(BIN)-$(GOOS)-$(GOARCH): $(GO_FILES)
	mkdir -p build
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build ${BUILD_ARGS} -o $@ ${PKG}

.PHONY: deps
deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: lint
lint:
	@for file in ${GO_FILES}; do golint $${file}; done

.PHONY: vet
vet:
	@go vet ${PKG_LIST}

.PHONY: test
test:
	@go test ${PKG_LIST}

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: clean
clean:
	rm -rf build

.PHONY: build
build: build/$(BIN)-$(GOOS)-$(GOARCH)

.PHONY: build-all
build-all:
	$(MAKE) build GOOS=linux GOARCH=amd64
	$(MAKE) build GOOS=darwin GOARCH=amd64

.PHONY: install
install: build/$(BIN)-$(GOOS)-$(GOARCH)
	mkdir -p "$(DESTDIR)$(PREFIX)/bin"
	cp "$<" "$(DESTDIR)$(PREFIX)/bin/$(BIN)"
	cp bin/rgit "$(DESTDIR)$(PREFIX)/bin/rgit"
	cp bin/ctags "$(DESTDIR)$(PREFIX)/bin/ctags"
	cp completion.zsh "$(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BIN)"

.PHONY: uninstall
uninstall:
	rm -f "$(DESTDIR)$(PREFIX)/bin/$(BIN)"
	rm -f "$(DESTDIR)$(PREFIX)/bin/rgit"
	rm -f "$(DESTDIR)$(PREFIX)/bin/ctags"
	rm -f "$(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BIN)"
