package := $(shell basename `pwd`)

LDFLAGS=-ldflags "-s -w"

.PHONY: default get codetest test fmt lint vet

default: fmt codetest

get:
	go get -v ./...
	go mod tidy

codetest: lint vet

fmt:
	go fmt ./...

lint:
	$(shell go env GOPATH)/bin/golangci-lint run --fix

vet:
	go vet -all .

test:
	go test -v -cover