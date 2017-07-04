BIN := ide
PKG := github.com/nrocco/ide
VERSION := $(shell git describe --tags --always --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v ${PKG}/vendor/)
GO_FILES := $(shell find * -type d -name vendor -prune -or -name '*.go' -type f | grep -v vendor)

LDFLAGS = "-d -s -w -X ${PKG}/cmd.Version=${VERSION}"
BUILD_ARGS = -a -tags netgo -installsuffix netgo -ldflags $(LDFLAGS)

PREFIX = /usr/local

.DEFAULT_GOAL: build/$(BIN)

build/$(BIN): $(GO_FILES)
	CGO_ENABLED=0 go build ${BUILD_ARGS} -o build/${BIN} ${PKG}

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

install: build/$(BIN)
	mkdir -p "$(DESTDIR)$(PREFIX)/bin"
	cp "$<" "$(DESTDIR)$(PREFIX)/bin/$(BIN)"
	# cp bin/composer "$(DESTDIR)$(PREFIX)/bin/composer"
	# cp bin/ctags "$(DESTDIR)$(PREFIX)/bin/ctags"
	# cp bin/node "$(DESTDIR)$(PREFIX)/bin/node"
	# cp bin/npm "$(DESTDIR)$(PREFIX)/bin/npm"
	# cp bin/php "$(DESTDIR)$(PREFIX)/bin/php"
	# cp bin/phpcoverage "$(DESTDIR)$(PREFIX)/bin/phpcoverage"
	# cp bin/phpunit "$(DESTDIR)$(PREFIX)/bin/phpunit"
	cp bin/rgit "$(DESTDIR)$(PREFIX)/bin/rgit"
	cp completion.zsh "$(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BIN)"

uninstall:
	rm -f "$(DESTDIR)$(PREFIX)/bin/$(BIN)"
	# rm -f "$(DESTDIR)$(PREFIX)/bin/composer"
	# rm -f "$(DESTDIR)$(PREFIX)/bin/ctags"
	# rm -f "$(DESTDIR)$(PREFIX)/bin/node"
	# rm -f "$(DESTDIR)$(PREFIX)/bin/npm"
	# rm -f "$(DESTDIR)$(PREFIX)/bin/php"
	# rm -f "$(DESTDIR)$(PREFIX)/bin/phpcoverage"
	# rm -f "$(DESTDIR)$(PREFIX)/bin/phpunit"
	rm -f "$(DESTDIR)$(PREFIX)/bin/rgit"
	rm -f "$(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BIN)"

.PHONY: build
build:
	mkdir -p build
	for GOOS in darwin linux; do \
		for GOARCH in amd64; do \
		    echo "==> Building ${BIN}-$$GOOS-$$GOARCH"; \
			docker run --rm -v "$(PWD)":/go/src/$(PKG) -w /go/src/$(PKG) -e "CGO_ENABLED=0" -e "GOOS=$$GOOS" -e "GOARCH=$$GOARCH" golang:1.9 \
				go build ${BUILD_ARGS} -o build/${BIN}-$$GOOS-$$GOARCH ${PKG}; \
		done; \
	done
