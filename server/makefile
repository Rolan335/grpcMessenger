.PHONY: generate test lint build

generate:
	go generate cmd/main.go

test:
	go test -v ./test

lint:
	golangci-lint run

build:
	mkdir -p ./build && \
	GOOS=linux GOARCH=amd64 go build -o ./build/server-linux cmd/main.go && \
	GOOS=darwin GOARCH=amd64 go build -o ./build/server-macos cmd/main.go && \
	GOOS=windows GOARCH=amd64 go build -o ./build/server-win.exe cmd/main.go
