# G4Api - Go REST API

Lightweight REST API for managing users and tasks with their relationships, written in Go.

## Requirements

- Docker & Docker Compose (recommended)
- Or Go 1.23+ and MariaDB 10.4+ / MySQL 5.7+

## Quick Start

### With Docker Compose

```bash
# Start the API (MariaDB + Go server)
docker-compose up -d api

# Run the tests
docker-compose run --rm test
```

The API will be available at `http://localhost:8080`.

### Manual Setup

1. Initialize the database:

```bash
mariadb -u root -p < share/sql/rest_api.sql
```

2. Set environment variables:

```bash
export DB_USER=root
export DB_PWD=
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=rest_api
export PORT=8080
```

3. Run the server:

```bash
go run ./cmd/server
```

## API Endpoints

| Method | URI | Description |
|--------|-----|-------------|
| `GET` | `/user/{id}` | Get user by ID |
| `GET` | `/user/{id}/task` | Get all tasks for a user |
| `POST` / `PUT` | `/user/{id}/task/{taskId}` | Associate a task with a user |
| `POST` / `PUT` | `/task` | Create a new task |
| `POST` / `PUT` | `/task/{id}` | Update a task |
| `DELETE` | `/task/{id}` | Delete a task |
| `DELETE` | `/user/{id}/task/{taskId}` | Remove task-user association |

## Examples

### Get a user

```bash
curl http://localhost:8080/user/1
```

### Create a task

```bash
curl -X POST http://localhost:8080/task \
  -H "Content-Type: application/json" \
  -d '{"title": "My Task", "description": "Task description", "status": 1}'
```

### Associate a task with a user

```bash
curl -X POST http://localhost:8080/user/1/task/1
```

## Testing

### With Docker (recommended, no Go required)

```bash
# Start MariaDB, initialize database, and run tests
docker-compose run --rm test
```

This starts a dedicated test container that:
1. Waits for MariaDB to be ready
2. Initializes the database with seed data
3. Runs all integration tests

### With Go installed

```bash
# Ensure MariaDB is running and initialized
mariadb -u root -p < share/sql/rest_api.sql

# Set database connection
export DB_USER=root
export DB_PWD=
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=rest_api

# Run tests
go test -v ./tests/
```

## Linting and Formatting

### With Docker (no Go required)

```bash
# Format code
docker run --rm -v ${PWD}:/app -w /app golang:1.23 gofmt -w .

# Run go vet
docker run --rm -v ${PWD}:/app -w /app golang:1.23 go vet ./...

# Run golangci-lint
docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:latest golangci-lint run
```

### With Go installed

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Format code
gofmt -w .
```

## Database

The project uses **MariaDB 11** by default via Docker Compose. TLS is disabled locally (`tls=skip-verify` in DSN, `--skip-ssl` for CLI) to avoid self-signed certificate errors. For production, configure proper TLS certificates and remove these flags.

## Project Structure

```
rest_api_go/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── api/
│   │   └── api.go           # API handler
│   ├── config/
│   │   ├── db.go            # Database configuration
│   │   └── route.go         # Route definitions
│   ├── handler/
│   │   ├── user.go          # User handlers
│   │   └── task.go          # Task handlers
│   ├── model/
│   │   ├── task.go          # Task model
│   │   ├── task_status.go   # Task status enum
│   │   ├── user.go          # User model
│   │   └── user_task.go     # User-Task relationship
│   └── router/
│       └── router.go        # HTTP router
├── share/
│   └── sql/
│       └── rest_api.sql     # Database schema + seed data
├── tests/
│   └── api_test.go          # Integration tests
├── .github/
│   └── workflows/
│       └── ci.yml           # GitHub Actions CI/CD
├── docker-compose.yml       # Docker Compose (MariaDB + API + tests)
├── Dockerfile               # Production image
├── Dockerfile.test          # Test image
├── .golangci.yml            # Linter configuration
├── go.mod                   # Go module definition
└── openapi.yaml             # OpenAPI specification
```

## License

MIT
