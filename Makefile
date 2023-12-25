SHELL=/bin/bash -o pipefail

.PHONY: docs build client install test vet chart

docs:
	rm -rf docs && mkdir docs
	rm -rf etc && mkdir -p etc/man/man1 && mkdir -p etc/completion
	go run cmd/gendoc/main.go

build:
	go build -o gmux ./

install:
	go install ./cmd/... 
