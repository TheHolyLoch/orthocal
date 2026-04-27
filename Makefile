# Orthocal - Developed by dgm (dgm@tuta.com)
# orthocal/Makefile

PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

BIN = bin/orthocal

.PHONY: build clean fmt install test uninstall

build:
	go build -o $(BIN) ./cmd/orthocal

clean:
	rm -rf bin

fmt:
	find . -name '*.go' -exec gofmt -w {} +

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 0755 $(BIN) $(DESTDIR)$(BINDIR)/orthocal

test:
	go test ./...

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/orthocal
