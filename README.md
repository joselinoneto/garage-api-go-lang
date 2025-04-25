# Garage API

A RESTful API for managing products with JWT authentication built with Go.

## Features

- Product CRUD operations
- JWT Authentication
- PostgreSQL database
- Swagger documentation
- Database migrations

## Prerequisites

- Go 1.21 or later
- PostgreSQL
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations

## Configuration

The API is configured to connect to a PostgreSQL database with the following credentials:

- Host: 192.168.1.2
- Port: 5432
- User: casaos
- Password: casaos
- Database: garage-web

## Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run database migrations:
   ```bash
   migrate -path migrations -database "postgresql://casaos:casaos@192.168.1.2:5432/garage-web?sslmode=disable" up
   ```

3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

The server will start on `http://localhost:8080`

## API Documentation

Swagger documentation is available at `http://localhost:8080/swagger/`

## API Endpoints

### Authentication

- POST `/api/v1/auth/register` - Register a new user
- POST `/api/v1/auth/login` - Login and get JWT token

### Products (Protected Routes)

- GET `/api/v1/products` - Get all products
- GET `/api/v1/products/{id}` - Get a specific product
- POST `/api/v1/products` - Create a new product
- PUT `/api/v1/products/{id}` - Update a product
- DELETE `/api/v1/products/{id}` - Delete a product

## Authentication

All product endpoints require JWT authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-token>
``` 