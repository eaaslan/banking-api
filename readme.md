# Project Structure

```
go-banking-api/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   │
│   ├── database/
│   │   └── database.go            # Database connection & migrations
│   │
│   ├── model/
│   │   ├── user.go                # User model with validations
│   │   ├── balance.go             # Balance model with validations
│   │   ├── transaction.go         # Transaction model with validations
│   │   └── audit_log.go           # AuditLog model with validations
│   │
│   ├── repository/                 # 🔜 Next: Database operations
│   │   ├── user_repository.go
│   │   ├── balance_repository.go
│   │   ├── transaction_repository.go
│   │   └── audit_log_repository.go
│   │
│   ├── service/                    # 🔜 Next: Business logic
│   │   ├── user_service.go
│   │   ├── balance_service.go
│   │   ├── transaction_service.go
│   │   └── audit_service.go
│   │
│   ├── handler/                    # 🔜 Next: HTTP handlers
│   │   ├── user_handler.go
│   │   ├── balance_handler.go
│   │   ├── transaction_handler.go
│   │   └── health_handler.go
│   │
│   ├── middleware/                 # 🔜 Next: HTTP middleware
│   │   ├── auth.go
│   │   ├── logger.go
│   │   ├── cors.go
│   │   └── rate_limit.go
│   │
│   └── util/                       # 🔜 Next: Utility functions
│       ├── hash.go                # Password hashing (bcrypt)
│       ├── jwt.go                 # JWT token generation
│       └── response.go            # Standard API responses
│
├── test/                           # 🔜 Next: Tests
│   ├── integration/
│   │   └── api_test.go
│   └── unit/
│       ├── model_test.go
│       ├── service_test.go
│       └── repository_test.go
│
├── scripts/
│   ├── test_models.go             # ✅ Model validation tests
│   └── test_api.sh                # ✅ API curl tests
│
├── .env                            # Environment variables
├── .env.example                    # Example environment variables
├── .gitignore                      # Git ignore file
├── docker-compose.yml              # 🔜 Docker setup
├── Dockerfile                      # 🔜 Docker build
├── go.mod                          # Go module file
├── go.sum                          # Go dependencies
└── README.md                       # Project documentation
```

---

## Current Status: ✅ Models Complete

### Completed:
1. ✅ **User Model** - Full validation
2. ✅ **Balance Model** - Full validation
3. ✅ **Transaction Model** - Full validation
4. ✅ **AuditLog Model** - Full validation
5. ✅ **Database Connection** - GORM setup with auto-migration
6. ✅ **Configuration** - Environment-based config
7. ✅ **Logging** - Structured logging with slog
8. ✅ **Graceful Shutdown** - Signal handling

---

## File Descriptions

### `/cmd/api/main.go`
- Application entry point
- Server initialization
- Graceful shutdown handling
- Basic HTTP handlers (temporary)

### `/internal/config/config.go`
- Environment variable loading
- Configuration struct
- Default values

### `/internal/database/database.go`
- GORM database connection
- Auto-migration for all models
- Connection pooling

### `/internal/model/*.go`
All models include:
- GORM struct tags
- Validation methods
- Business logic methods
- GORM hooks (BeforeCreate, BeforeUpdate)
- Table name specification

---

## Environment Variables

Required in `.env` file:
```env
# Application
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL=postgres://user:password@localhost:5432/banking_db?sslmode=disable

# JWT (coming soon)
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h

# Rate Limiting (coming soon)
RATE_LIMIT=100
RATE_WINDOW=1m
```

---

## Running the Application

### 1. Setup Database
```bash
# Using Docker
docker run --name postgres-banking \
  -e POSTGRES_USER=bankuser \
  -e POSTGRES_PASSWORD=bankpass \
  -e POSTGRES_DB=banking_db \
  -p 5432:5432 \
  -d postgres:15

# Or use your local PostgreSQL
createdb banking_db
```

### 2. Configure Environment
```bash
# Copy example env file
cp .env.example .env

# Edit .env with your settings
nano .env
```

### 3. Run the Application
```bash
# Download dependencies
go mod tidy

# Run
go run cmd/api/main.go

# Or build and run
go build -o bin/api cmd/api/main.go
./bin/api
```

### 4. Test Models
```bash
# Run model validation tests
go run scripts/test_models.go

# Test API endpoints
chmod +x scripts/test_api.sh
./scripts/test_api.sh
```

---

## Database Schema

Tables created automatically on startup:
- `users` - User accounts
- `balances` - User balances (1-to-1 with users)
- `transactions` - All financial transactions
- `audit_logs` - System audit trail

All tables include:
- `id` (primary key, auto-increment)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `deleted_at` (soft delete, nullable)

---

## Next Steps

### Phase 1: Repository Layer (Next)
Create database operation abstractions:
- CRUD operations
- Query methods
- Transaction handling

### Phase 2: Service Layer
Implement business logic:
- User registration/login
- Balance operations
- Transaction processing
- Audit logging

### Phase 3: HTTP Layer
Build REST API:
- Route setup
- Request validation
- Response formatting
- Error handling

### Phase 4: Security
Add authentication & authorization:
- JWT tokens
- Password hashing (bcrypt)
- Role-based access
- Rate limiting

### Phase 5: Testing
Comprehensive tests:
- Unit tests
- Integration tests
- API tests
- Load tests

### Phase 6: Deployment
Production ready:
- Docker containerization
- CI/CD pipeline
- Monitoring setup
- Documentation

---

## Development Guidelines

### Code Style
- Use Go conventions (gofmt, golint)
- Write descriptive comments
- Keep functions small and focused
- Use meaningful variable names

### Error Handling
- Always handle errors
- Return descriptive error messages
- Use custom error types when needed
- Log errors appropriately

### Database
- Use transactions for multi-step operations
- Add indexes for frequently queried fields
- Keep queries optimized
- Use prepared statements

### Security
- Never log sensitive data
- Validate all inputs
- Use parameterized queries
- Hash passwords with bcrypt
- Implement rate limiting

---

## Useful Commands

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test ./...

# Run with race detector
go run -race cmd/api/main.go

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go

# Check dependencies
go mod verify
go mod tidy

# Update dependencies
go get -u ./...
```

---

## Contributing

1. Create feature branch
2. Write tests
3. Implement feature
4. Run all tests
5. Format code
6. Submit PR

---

## License

MIT License - See LICENSE file for details