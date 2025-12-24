# Prerequisite

- Go 1.25.5

# Quickstart
```bash
make install # or go mod download
make # or go run cmd/main.go
```

# API Endpoints

- **GraphQL API**: `/query`
- **GraphQL Playground**: `/`
- **Health Check**: `/health`

# Environment Variables

The necessary environment variables can be seen in the .env.example file.

# Development

## Generate GraphQL Code
```bash
make generate # or go generate ./...
```

## Format Code
```bash
make format
```

## Run Linter
```bash
make lint # or golangci-lint run ./...
```

## Run Test
```bash
make test # or go test ./...
```