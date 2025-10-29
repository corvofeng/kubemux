SHELL=/bin/bash -o pipefail

.PHONY: docs build build-mcp client install test vet chart

docs:
	rm -rf docs && mkdir docs
	rm -rf etc && mkdir -p etc/man/man1 && mkdir -p etc/completion
	go run cmd/gendoc/main.go

test:
	go test kubemux/lib -v

build:
	go build -o kubemux ./cmd/kubemux

build-mcp:
	go build -o kubemux-mcp-server ./mcp-server

build-all: build build-mcp

install:
	go mod tidy
