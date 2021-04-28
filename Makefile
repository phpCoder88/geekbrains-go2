build:
	cd ./duplicate-file-finder && go build -o finder

test:
	go test ./...

check:
	golangci-lint run