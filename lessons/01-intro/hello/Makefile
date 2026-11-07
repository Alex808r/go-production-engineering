.DEFAULT_GOAL := fmt
.PHONY: fmt vet build clean

fmt:
	go fmt ./...
vet: fmt
	go vet ./...
build:
	go build
clean:
	go clean