# Top-Up API

## Overview

Top-Up API is a backend service built with **Go (Golang)** for mobile top-up and prepaid card systems. It provides RESTful APIs for managing orders, SKUs, suppliers, providers, and purchase histories. The project supports authentication, order processing, and integration with PostgreSQL, Redis, Kafka, and gRPC.

## Project Structure

The project follows a clean architecture pattern with the following main directories:

- `cmd/api/` - Application entry point
- `internal/` - Private application code
  - `app/` - Application setup and initialization
  - `controller/http/` - HTTP handlers and routing
  - `grpc/` - gRPC client and server implementations
  - `kafka/` - Kafka consumers
  - `service/` - Business logic layer
  - `repository/` - Data access layer
  - `model/` - Domain models
  - `schema/` - Request/response schemas
- `pkg/` - Reusable packages
- `proto/` - Protocol buffer definitions
- `tests/` - Test files and mocks
- `config/` - Configuration files
- `docs/` - Swagger documentation
- `sql/` - Database initialization scripts

## Tech Stack

- **Language:** Go 1.24+
- **Framework:** Gin (HTTP server)
- **Database:** PostgreSQL (via GORM)
- **Cache:** Redis
- **Messaging:** Apache Kafka
- **gRPC:** Internal service communication
- **Documentation:** Swagger/OpenAPI
- **Testing:** Testify
- **Logging:** Zap

## Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL
- Redis
- Apache Kafka

### Setup Steps

1. **Clone the repository**

   ```sh
   git clone https://github.com/mrkidvn44/top-up-api
   cd top-up-api
   ```

2. **Copy and edit configuration**

   ```sh
   cp config/example.yaml config/config.yaml
   # Edit the configuration file as needed for your environment
   ```

3. **Start dependencies with Docker Compose**

   ```sh
   docker compose up -d
   ```

4. **Install Go dependencies**

   ```sh
   cd main-service/top-up
   go mod download
   ```

5. **Run database migrations** (if needed)

   ```sh
   # The database will be initialized with the SQL scripts in sql/ directory
   ```

6. **Run the API server**

   ```sh
   go run cmd/api/main.go
   ```

7. **Access API documentation**
   - Visit: `http://localhost:8080/swagger/index.html` (default port, check your config)

## Running Tests

```sh
go test ./...
```

To run tests with coverage:

```sh
go test -v -cover ./...
```

## Configuration

The application uses YAML configuration files located in the `config/` directory:

- `config.yaml` - Main configuration file (create from example.yaml)
- `example.yaml` - Example configuration with default values

Key configuration sections:
- **HTTP Server:** Port and server settings
- **Database:** PostgreSQL connection details
- **Redis:** Cache configuration
- **Kafka:** Message broker settings
- **JWT:** Authentication settings
- **Logging:** Log level and format

## API Endpoints

The API provides the following main endpoints:

- **Orders:** `/order/*` - Order management and processing
- **SKUs:** `/sku/*` - Stock Keeping Unit operations
- **Suppliers:** `/supplier/*` - Supplier management
- **Purchase History:** `/purchase-history/*` - Transaction history
- **Health Check:** Health and status endpoints

## API Documentation

- Swagger/OpenAPI documentation is automatically generated and served at `/swagger/index.html`
- API documentation source files are in the `docs/` directory

## Development

### Project Layout

This project follows the standard Go project layout:

- `/cmd` - Main applications
- `/internal` - Private application and library code
- `/pkg` - Library code that can be used by external applications
- `/proto` - Protocol buffer files
- `/tests` - Additional test utilities and mock data
- `/docs` - Design and user documents
- `/sql` - Database schemas and migrations

### Code Generation

The project uses several code generation tools:

```sh
# Generate Swagger documentation
swag init -g cmd/api/main.go

# Generate Protocol Buffer files (if proto files are modified)
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

### Docker Support

The project includes Docker Compose configuration for running dependencies:

- **PostgreSQL** - Main database (port 5432)
- **Redis** - Cache store (port 6379)  
- **Apache Kafka** - Message broker (port 9092)
