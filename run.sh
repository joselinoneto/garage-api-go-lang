#!/bin/bash

# Set environment variables
export DB_HOST=192.168.1.2
export DB_PORT=5432
export DB_USER=casaos
export DB_PASSWORD=casaos
export DB_NAME=garage-web
export DB_SSL_MODE=disable

# Run the application
go run cmd/api/main.go 