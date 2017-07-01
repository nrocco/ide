BINARY := ide
PKG := github.com/nrocco/ide
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

LDFLAGS = "-X ${PKG}/cmd.Version=${VERSION}"

PREFIX = /usr/local

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(GO_FILES)
	go build -i -v -o ${BINARY} -ldflags ${LDFLAGS} ${PKG}

lint:
	@for file in ${GO_FILES}; do \
		golint $${file}; \
	done

vet:
	@go vet ${PKG_LIST}

test:
	@go test -short ${PKG_LIST}

install: $(BINARY)
	mkdir -p "$(DESTDIR)$(PREFIX)/bin"
	cp "$<" "$(DESTDIR)$(PREFIX)/bin/$(BINARY)"
	cp bin/composer "$(DESTDIR)$(PREFIX)/bin/composer"
	cp bin/ctags "$(DESTDIR)$(PREFIX)/bin/ctags"
	cp bin/node "$(DESTDIR)$(PREFIX)/bin/node"
	cp bin/npm "$(DESTDIR)$(PREFIX)/bin/npm"
	cp bin/php "$(DESTDIR)$(PREFIX)/bin/php"
	cp bin/phpcoverage "$(DESTDIR)$(PREFIX)/bin/phpcoverage"
	cp bin/phpunit "$(DESTDIR)$(PREFIX)/bin/phpunit"
	cp bin/rgit "$(DESTDIR)$(PREFIX)/bin/rgit"
	cp completion.zsh "$(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BINARY)"

uninstall:
	rm -f "$(DESTDIR)$(PREFIX)/bin/$(BINARY)"
	rm -f "$(DESTDIR)$(PREFIX)/bin/composer"
	rm -f "$(DESTDIR)$(PREFIX)/bin/ctags"
	rm -f "$(DESTDIR)$(PREFIX)/bin/node"
	rm -f "$(DESTDIR)$(PREFIX)/bin/npm"
	rm -f "$(DESTDIR)$(PREFIX)/bin/php"
	rm -f "$(DESTDIR)$(PREFIX)/bin/phpcoverage"
	rm -f "$(DESTDIR)$(PREFIX)/bin/phpunit"
	rm -f "$(DESTDIR)$(PREFIX)/bin/rgit"
	rm -f "$(DESTDIR)$(PREFIX)/share/zsh/site-functions/_$(BINARY)"

.PHONY: clean
clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi
