# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install required packages
RUN go get github.com/gin-gonic/gin
RUN go get github.com/sirupsen/logrus
RUN go get github.com/golang-jwt/jwt/v5
RUN go get github.com/joho/godotenv
RUN go get github.com/swaggo/gin-swagger
RUN go get github.com/swaggo/files

# Copy source code
COPY . .

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