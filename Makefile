DIST_DIR = dist

PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
NAME = ghibp

.PHONY: all clean build windows mac install uninstall test prep

all: build

prep:
	mkdir -p $(DIST_DIR)/

build: prep
	go build -o $(DIST_DIR)/

linux: prep
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/ghibp-linux-amd64

windows: prep
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/ghibp-windows-amd64.exe

mac: prep
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/ghibp-darwin-arm64

install:
	install -m755 $(DIST_DIR)/$(NAME) $(BINDIR)

uninstall:
	rm $(BINDIR)/$(NAME)

clean:
	-rm $(DIST_DIR)/ghibp*
