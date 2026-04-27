# Orthocal - Developed by dgm (dgm@tuta.com)
# orthocal/Makefile

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/man
BASH_COMPLETION_DIR ?= $(PREFIX)/share/bash-completion/completions
ZSH_COMPLETION_DIR ?= $(PREFIX)/share/zsh/site-functions
FISH_COMPLETION_DIR ?= $(PREFIX)/share/fish/vendor_completions.d
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || printf unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || printf unknown)

BIN = bin/orthocal
LDFLAGS = -X orthocal/internal/version.Version=$(VERSION) -X orthocal/internal/version.Commit=$(COMMIT) -X orthocal/internal/version.BuildDate=$(BUILD_DATE)

.PHONY: build clean fmt install install-all install-completions install-man test uninstall

build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/orthocal

clean:
	rm -rf bin

fmt:
	find . -name '*.go' -exec gofmt -w {} +

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 0755 $(BIN) $(DESTDIR)$(BINDIR)/orthocal

install-all: install install-completions install-man

install-completions:
	install -d $(DESTDIR)$(BASH_COMPLETION_DIR)
	install -d $(DESTDIR)$(ZSH_COMPLETION_DIR)
	install -d $(DESTDIR)$(FISH_COMPLETION_DIR)
	install -m 0644 completions/orthocal.bash $(DESTDIR)$(BASH_COMPLETION_DIR)/orthocal
	install -m 0644 completions/orthocal.zsh $(DESTDIR)$(ZSH_COMPLETION_DIR)/_orthocal
	install -m 0644 completions/orthocal.fish $(DESTDIR)$(FISH_COMPLETION_DIR)/orthocal.fish

install-man:
	install -d $(DESTDIR)$(MANDIR)/man1
	install -m 0644 docs/orthocal.1 $(DESTDIR)$(MANDIR)/man1/orthocal.1

test:
	go test ./...

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/orthocal
	rm -f $(DESTDIR)$(BASH_COMPLETION_DIR)/orthocal
	rm -f $(DESTDIR)$(ZSH_COMPLETION_DIR)/_orthocal
	rm -f $(DESTDIR)$(FISH_COMPLETION_DIR)/orthocal.fish
	rm -f $(DESTDIR)$(MANDIR)/man1/orthocal.1
