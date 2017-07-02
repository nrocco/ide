BIN := ide
PKG := github.com/nrocco/ide
VERSION := $(shell git describe --tags --always --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v ${PKG}/vendor/)
GO_FILES := $(shell find * -type d -name vendor -prune -or -name '*.go' -type f | grep -v vendor)

LDFLAGS = "-d -s -w -X ${PKG}/cmd.Version=${VERSION}"

PREFIX = /usr/local

.DEFAULT_GOAL: $(BIN)

$(BIN): $(GO_FILES)
	go build -i -v -o ${BIN} -ldflags ${LDFLAGS} ${PKG}

deps:
	dep ensure

lint:
	@for file in ${GO_FILES}; do golint $${file}; done

vet:
	@go vet ${PKG_LIST}

test:
	@go test -short ${PKG_LIST}

version:
	@echo $(VERSION)

.PHONY: clean
clean:
	if [ -f ${BIN} ]; then rm ${BIN}; fi

install: $(BIN)
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
	mkdir -p build/
	for GOOS in darwin linux; do \
		for GOARCH in amd64; do \
		    echo "==> Building ide for $$GOOS $$GOARCH"; \
			docker run --rm -v "$(PWD)":/go/src/$(PKG) -w /go/src/$(PKG) -e "CGO_ENABLED=0" -e "GOOS=$$GOOS" -e "GOARCH=$$GOARCH" golang:1.9 \
				go build -a -x -v -tags netgo -installsuffix netgo -o build/${BIN}-$$GOOS-$$GOARCH -ldflags ${LDFLAGS} ${PKG}; \
		done; \
	done
