BINARY  := sg
PKG     := ./cmd/sg
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X snapgit/internal/cli.Version=$(VERSION)

.PHONY: build test vet lint check install clean

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(PKG)

test:
	go test -race -count=1 ./...

vet:
	go vet ./...

lint: vet
	@echo "all checks passed"

check: vet test build
	@echo "build, vet, and tests passed"

install: build
	cp $(BINARY) $(GOPATH)/bin/ 2>/dev/null || cp $(BINARY) /usr/local/bin/

clean:
	rm -f $(BINARY)
