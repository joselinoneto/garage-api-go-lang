# syntax=docker/dockerfile:1

# Build stage
FROM --platform=linux/arm/v7 golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download
RUN go mod tidy

# Copy source code
COPY . .

# Build the application for ARM v7 (Raspberry Pi)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o main cmd/api/main.go

# Final stage
FROM --platform=linux/arm/v7 alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata wget

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Run the application
CMD ["./main"] 