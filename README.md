# Simple Banking API Backend

This is a comprehensive Banking API backend built with Go, implementing Clean Architecture, concurrent processing, and a robust telemetry stack.

## Key Features

- **Architecture**:
  - **Clean Architecture**: Separation of concerns (Handler -> Service -> Repository).
  - **Dependency Injection**: Modular and testable code structure.
  - **Worker Pool**: Asynchronous transaction processing for high throughput.
  - **Redis Caching**: Improved performance for balance inquiries using Cache-Aside pattern.

- **Security**:
  - **JWT Authentication**: Secure API access.
  - **Password Hashing**: Bcrypt for password security.
  - **Role-Based Access Control (RBAC)**: Admin-specific endpoints.

- **Database**:
  - **PostgreSQL**: Production-grade relational database.
  - **Migrations**: Automated schema management on startup.
  - **Transactions**: ACID compliance with proper rollback mechanisms.

- **Observability & Telemetry**:
  - **Structured Logging**: JSON logging (slog) for production.
  - **Prometheus Metrics**: Custom metrics for API request duration and transaction counts.
  - **Grafana Dashboards**: Visualizing key performance indicators.
  - **Distributed Tracing**: OpenTelemetry integration with Jaeger for request tracing.

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Telemetry**: OpenTelemetry, Prometheus, Grafana, Jaeger
- **Containerization**: Docker & Docker Compose

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development without Docker)

### Running the Application (Recommended)

The easiest way to run the entire stack is using Docker Compose:

```bash
docker-compose up --build
```

This will start the following services:
- **API**: `http://localhost:8080`
- **PostgreSQL**: `localhost:5432`
- **Redis**: `localhost:6379`
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (User/Pass: admin/admin)
- **Jaeger UI**: `http://localhost:16686`

### Running Locally (Without Docker)

If you prefer to run the Go application locally, ensure you have PostgreSQL and Redis running, then configure the environment variables in a `.env` file or export them directly.

```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/api/main.go
```

## API Documentation

The API includes endpoints for User Management, Authentication, Transactions, and Reporting.

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and receive JWT
- `POST /api/v1/auth/refresh` - Refresh access token

### Transactions (Authenticated)
- `POST /api/v1/transactions` - Create a new transaction (Deposit, Withdraw, Transfer)
- `GET /api/v1/transactions/history` - Get transaction history

### Balances (Authenticated)
- `GET /api/v1/balances/current` - Get current balance (Cached via Redis)
- `GET /api/v1/balances/historical` - Get historical balance data

### User Management (Admin Only)
- `GET /api/v1/users` - List all users
- `DELETE /api/v1/users/delete?id={id}` - Delete a user

## Monitoring

- **Metrics**: Access `http://localhost:9090` to query Prometheus metrics (e.g., `http_requests_total`).
- **Dashboards**: Access `http://localhost:3000` for Grafana dashboards.
- **Tracing**: Access `http://localhost:16686` to view traces in Jaeger.

## Configuration

The application is configured via environment variables. Key variables include:

- `SERVER_PORT`: Port to run the server on (default: 8080)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`: Database connection details.
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`: Redis connection details.
- `OTEL_EXPORTER_OTLP_ENDPOINT`: OpenTelemetry collector endpoint.
