.PHONY: build run test lint fmt clean docker-up docker-down docker-test

build:
	go build -o server ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./tests/

lint:
	golangci-lint run

fmt:
	gofmt -w .
	goimports -w .

clean:
	rm -f server

docker-up:
	docker-compose up -d api

docker-down:
	docker-compose down

docker-test:
	docker-compose run --rm test
