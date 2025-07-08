# Top-Up API

## Overview

Top-Up API is a backend service built with **Go (Golang)** for mobile top-up and prepaid card systems. It provides RESTful APIs for managing users, orders, card details, providers, and purchase histories. The project supports authentication, order processing, and integration with PostgreSQL, Redis, Kafka, and gRPC.

## Project Structure

- `main-service/`  
  The core API service for production use.
- `mock-service/`  
  Optional mock services for development, testing, and integration (includes `provider/` and `payment/`).

## Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin (HTTP server)
- **Database:** PostgreSQL (via GORM)
- **Cache:** Redis
- **Messaging:** Kafka
- **gRPC:** Internal service communication
- **Swagger:** API documentation

## Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL
- Redis
- Kafka

### Setup Steps

1. **Clone the repository**

   ```sh
   git clone https://github.com/mrkidvn44/top-up-api
   cd top-up-api
   ```

2. **Copy and edit configuration**

   ```sh
   cp main-service/top-up/config/example.yaml main-service/top-up/config/config.yaml
   # Edit as needed
   ```

3. **Start dependencies**

   ```sh
   docker compose up -d
   ```

4. **Install Go dependencies**

   ```sh
   cd main-service/top-up
   go mod download
   ```

5. **Run the main API server**

   ```sh
   go run cmd/api/main.go
   ```

6. **(Optional) Run mock services**

   ```sh
   cd ../../mock-service/provider
   go run cmd/api/main.go
   # Similarly for mock-service/payment if needed
   ```

7. **Access API documentation**
   - Visit: `http://localhost:<port>/swagger/index.html`

## Running Tests

```sh
cd main-service/top-up
go test ./...
```

## API Documentation

- Swagger/OpenAPI docs are available in the `docs/` folder of each service and are served by the API server.

---
