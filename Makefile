install:
	go mod download

run:
	go run cmd/main.go

test:
	go test ./...

deps:
	go mod tidy

lint:
	golangci-lint run ./...

generate:
	go generate ./...

format:
	go fmt ./...
	gofmt -s -w .

.DEFAULT_GOAL = run