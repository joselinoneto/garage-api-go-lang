# Build stage
FROM golang:1.21-alpine AS builder

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

# Generate go.sum
RUN go mod tidy

# Build the application for ARM
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o main cmd/api/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 