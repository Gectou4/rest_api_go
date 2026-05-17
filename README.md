# G4Api - Go REST API

Lightweight REST API for managing users and tasks with their relationships, written in Go.

## Requirements

- Go 1.23+
- MySQL 5.7+ / MariaDB 10.4+

## Quick Start

### With Docker Compose

```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`.

### Manual Setup

1. Initialize the database:

```bash
mysql -u root -p < share/sql/rest_api.sql
```

2. Set environment variables (optional):

```bash
export DB_USER=root
export DB_PWD=
export DB_DSN="tcp(127.0.0.1:3306)/rest_api?parseTime=true"
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

### With Docker (no Go required)

```bash
# Start MySQL, initialize database, and run tests
docker run --rm --network host -v ${PWD}:/app -w /app golang:1.23 sh -c "
  go mod tidy &&
  until mysqladmin ping -h127.0.0.1 -uroot -proot --silent; do sleep 2; done &&
  mysql -h127.0.0.1 -uroot -proot rest_api < share/sql/rest_api.sql &&
  DB_USER=root DB_PWD=root DB_DSN='tcp(127.0.0.1:3306)/rest_api?parseTime=true' go test -v ./tests/
"
```

### With Go installed

```bash
# Ensure MySQL is running and initialized
mysql -u root -p < share/sql/rest_api.sql

# Set database connection
export DB_USER=root
export DB_PWD=root
export DB_DSN="tcp(127.0.0.1:3306)/rest_api?parseTime=true"

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
├── docker-compose.yml       # Docker Compose for local dev
├── Dockerfile               # Docker image
├── .golangci.yml            # Linter configuration
├── go.mod                   # Go module definition
└── openapi.yaml             # OpenAPI specification
```

## License

MIT
