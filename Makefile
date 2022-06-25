GO=go
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS)
GO_IMPORTS=goimports
GO_LDFLAGS=-ldflags="-s -w"
TARGET_DIR=bin/
CONTAINER_IMAGE=guni1192/forestry:dev

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

lint:
	golangci-lint run

clean:
	rm -rf bin

docker-build:
	DOCKER_BUILDKIT=1 docker build -t $(CONTAINER_IMAGE) .

docker-run:
	docker run -d --cidfile=/tmp/forestry-dev-server -p 1192:1192 $(CONTAINER_IMAGE)

docker-stop:
	docker stop $(shell cat /tmp/forestry-dev-server)
	rm -rf /tmp/forestry-dev-server

deploy:
	kubectl apply -f k8s/
