prefix ?= /usr

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
	ln -s "$(DESTDIR)$(prefix)/share/ide/bin/ide" "$(DESTDIR)$(prefix)/bin/ide"
