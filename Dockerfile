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

# Create non-root app user
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# Copy the binary from builder
COPY --from=builder /app/gapi .

# Set proper permissions for app user
RUN chown -R appuser:appgroup /root/

# Switch to app user
USER appuser

EXPOSE 8080

CMD ["./gapi"]
