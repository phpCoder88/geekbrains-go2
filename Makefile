HAS_LINT := $(shell command -v golangci-lint;)

build:
	cd ./duplicate-file-finder && go build -o finder

test:
	go test ./...

check: bootstrap
	golangci-lint run


bootstrap:
ifndef HAS_LINT
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0
endif