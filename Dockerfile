# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gapi ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install runtime dependencies if needed
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/gapi .

# Copy any config files if needed
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./gapi"]

