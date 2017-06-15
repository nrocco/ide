prefix ?= /usr/local

.PHONY: install
install:
ifeq ($(shell uname -s),Darwin)
	/usr/bin/install -d  "$(DESTDIR)$(prefix)/share/ide"
	find bin lib exec hooks -type d -exec install -d "$(DESTDIR)$(prefix)/share/ide/{}" \;
	find bin lib exec hooks -type f -exec install -m755 "{}" "$(DESTDIR)$(prefix)/share/ide/{}" \;
	find bin lib exec hooks -type l -exec install -m755 "{}" "$(DESTDIR)$(prefix)/share/ide/{}" \;
else
	find bin lib exec hooks -type f -exec install -Dm755 "{}" "$(DESTDIR)$(prefix)/share/ide/{}" \;
	find bin lib exec hooks -type l -exec install -Dm755 "{}" "$(DESTDIR)$(prefix)/share/ide/{}" \;
endif
	/usr/bin/install -d  "$(DESTDIR)$(prefix)/bin"
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/composer" "$(DESTDIR)$(prefix)/bin/composer"
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/ctags" "$(DESTDIR)$(prefix)/bin/ctags"
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/ide" "$(DESTDIR)$(prefix)/bin/ide"
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/node" "$(DESTDIR)$(prefix)/bin/node"
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/npm" "$(DESTDIR)$(prefix)/bin/npm"
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/php" "$(DESTDIR)$(prefix)/bin/php"
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/phpcoverage" "$(DESTDIR)$(prefix)/bin/phpcoverage"
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/phpunit" "$(DESTDIR)$(prefix)/bin/phpunit"
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/rgit" "$(DESTDIR)$(prefix)/bin/rgit"
