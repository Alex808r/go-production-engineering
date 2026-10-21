.PHONY: all build run test test-concurrency test-race docker-up docker-down clean

all: test build

build:
	go build -o bin/wallet-app ./cmd/app

run:
	go run ./cmd/app

test:
	go test -v ./...

test-race:
	go test -v -race ./...

test-concurrency:
	go test -v -run TestConcurrency ./tests

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v

clean:
	rm -rf bin/
