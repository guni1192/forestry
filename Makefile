GO=go
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS)
GO_IMPORTS=goimports
GO_LDFLAGS=-ldflags="-s -w"
TARGET_DIR=bin/

.PHONY: build test fmt vet clean

build:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=$(GOOS) GO_ARCH=$(GOARCH) $(GO) build $(GO_LDFLAGS) -o $(TARGET_DIR) ./...

test:
	$(GO) test -v ./...

fmt:
	$(GO_IMPORTS) -w .

vet:
	$(GO) vet -v ./...

clean:
	rm -rf bin
