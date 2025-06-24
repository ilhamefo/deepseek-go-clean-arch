# DeepSeek Go Clean Architecture

A Go project implementing Clean Architecture principles for an event registration system, featuring authentication (including Google OAuth), JWT, GORM, and comprehensive testing with sqlmock.

## Features
- Clean Architecture structure (domain, service, repository, etc.)
- Authentication with JWT and Google OAuth2
- GORM for database access
- Test suite with sqlmock for integration-like tests
- Configurable via environment variables

## Prerequisites
- Go 1.20 or newer
- PostgreSQL (for production use)
- [Git](https://git-scm.com/)

## Installation

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd deepseek-go-clean-arch
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Configure environment:**
   - Copy `.env.example` to `.env` and edit as needed (or set environment variables directly).
   - Example variables:
     - `DB_HOST`, `DB_USER`, `DB_PASS`, `DB_NAME`, `JWT_SECRET`, etc.

## Running the Application

1. **Run the server:**
   ```bash
   go run ./cmd/server
   ```
   Or, if you have a Makefile:
   ```bash
   make run
   ```

2. **API Documentation:**
   - See `docs/swagger.yaml` or `docs/swagger.json` for OpenAPI docs.

## Running Tests

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Project Structure

- `cmd/` - Application entrypoints (main server, exporters)
- `internal/core/` - Domain logic, services, interfaces
- `internal/repository/` - GORM and other repository implementations
- `internal/common/` - Shared code, config, helpers
- `docs/` - API documentation (Swagger/OpenAPI)
- `files/` - Example data files
- `scripts/` - Utility scripts

## License

## Future add

- `securego/gosec` For static code analysis to find security issues.
- `prometheus/client_golang` For exporting application metrics to Prometheus
- `open-telemetry/opentelemetry-go` For distributed tracing and metrics, especially useful in microservices.

MIT License
