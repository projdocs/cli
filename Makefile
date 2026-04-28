BINARY  := projdocs
MODULE  := $(shell go list -m)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X '$(MODULE)/pkg.Version=$(VERSION)'"
GOOS := darwin
GOARCH := arm64
.PHONY: build clean

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BINARY) .

dist: build
	mkdir -p ./dist
	cp $(BINARY) ./dist/$(BINARY)-$(VERSION)-$(GOOS)-$(GOARCH)

clean:
	rm -f $(BINARY)
