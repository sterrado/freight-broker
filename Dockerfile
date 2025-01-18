# Build stage
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]