DIST_DIR = dist

PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man/man1

NAME = ghibp

.PHONY: all clean build windows mac install uninstall test prep man docs

all: build man docs

prep:
	mkdir -p $(DIST_DIR)/
	mkdir -p $(DIST_DIR)/man

build: prep
	go build -o $(DIST_DIR)/

man: build
	$(DIST_DIR)/$(NAME) docs --disable-markdown=true

docs: build
	$(DIST_DIR)/$(NAME) docs --disable-manpages=true

linux: prep
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/ghibp-linux-amd64

windows: prep
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/ghibp-windows-amd64.exe

mac: prep
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/ghibp-darwin-arm64

install:
	install -m755 -d $(BINDIR)
	install -m755 $(DIST_DIR)/$(NAME) $(BINDIR)
	install -m755 -d $(MANDIR)
	install -m644 $(DIST_DIR)/man/$(NAME)*.1 $(MANDIR)

uninstall:
	rm $(BINDIR)/$(NAME)
	rm ${MANDIR}/$(NAME)*

clean:
	-rm $(DIST_DIR)/ghibp*
	-rm -rf $(DIST_DIR)/man/*
