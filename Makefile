BINARY  := projdocs
MODULE  := $(shell go list -m)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X '$(MODULE)/pkg.Version=$(VERSION)'"
.PHONY: build clean

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY) .

clean:
	rm -f $(BINARY)
