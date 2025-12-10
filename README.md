# Simple Bank Backend

This is a junior-level Go backend project implementing a simple banking system with concurrent processing.

## Features
- **User Management**: Registration and Authentication (Bcrypt).
- **Transaction Processing**: 
  - Deposit, Withdraw, and Transfer operations.
  - Asynchronous processing using a Worker Pool.
  - Transaction Rollback simulation for transfers.
- **Balance Management**: Thread-safe operations using `sync.RWMutex`.
- **Database**: SQLite with graceful schema migration on startup.
- **Architecture**: 
  - Clean Architecture (Handler -> Service -> Repository).
  - Configurable via environment variables.

## Getting Started

### Prerequisites
- Go 1.21+

### Running the Application
```sh
# Install dependencies
go mod tidy

# Run the server
go run cmd/api/main.go
```

The server will start on port 8080 (default).

## API Endpoints

- `POST /api/register` - Create a new user
- `POST /api/login` - Authenticate
- `POST /api/transactions` - Create a transaction
- `GET /api/balance?user_id=1` - Get user balance

## Project Structure
- `cmd/api`: Entry point.
- `internal/config`: Configuration loader.
- `internal/db`: Database connection and migrations.
- `internal/models`: Domain entities.
- `internal/repository`: Data access layer.
- `internal/service`: Business logic.
- `internal/worker`: Concurrent task processor.
- `internal/handler`: HTTP request handlers.
