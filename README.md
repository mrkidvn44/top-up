# Top-Up API

## Overview

Top-Up API is a backend service built with **Go (Golang)**. It provides APIs for managing users, orders, card details, providers, and purchase histories, typically for a mobile top-up or prepaid card system. The project uses RESTful endpoints and supports authentication, order processing, and integration with external services like Redis, Kafka, and PostgreSQL.

## Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin (HTTP server)
- **Database:** PostgreSQL (via GORM)
- **Cache:** Redis
- **Messaging:** Kafka
- **gRPC:** For internal service communication
- **Swagger:** API documentation

## Project Structure

- `cmd/api/` - Main entry point for the API server
- `internal/` - Application logic (controllers, services, models, etc.)
- `pkg/` - Shared packages (logger, auth, utils, etc.)
- `config/` - Configuration files and loader
- `docs/` - Swagger/OpenAPI documentation
- `proto/` - gRPC protobuf definitions
- `sql/` - Database initialization scripts
- `tests/` - Unit and integration tests

## Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose (for running dependencies)
- PostgreSQL
- Redis
- Kafka

### Setup Steps

1. **Clone the repository**
    ```sh
    git clone <your-repo-url>
    cd top-up-api
    ```

2. **Copy and edit configuration**
    ```sh
    cp config/example.yaml config/config.yaml
    # Edit config/config.yaml as needed
    ```

3. **Start dependencies with Docker Compose**
    ```sh
    docker compose up -d
    ```

4. **Install Go dependencies**
    ```sh
    go mod download
    ```

5. **Run the API server**
    ```sh
    go run cmd/api/main.go
    ```

6. **Access API documentation**
    - Visit: `http://localhost:<port>/swagger/index.html`

## Running Tests

```sh
go test ./...
```

## API Documentation

- Swagger/OpenAPI docs are available in the `docs/` folder and served by the API server.

---