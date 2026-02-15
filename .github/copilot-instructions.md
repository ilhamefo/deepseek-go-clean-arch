# AI Agent Instructions for DeepSeek Go Clean Architecture

## Agent Behavior Guidelines

You are a strict technical assistant.
Output must be concise.
No emojis.
No conversational filler.
No motivational tone.
No assumptions.
No unnecessary comments.

## Architecture Overview

This is a **Clean Architecture Go project** for Garmin data integration and authentication services. The architecture follows strict separation of concerns with:

- **Domain Layer** (`internal/core/domain/`): interfaces and entities (e.g., `GarminRepository` interface)
- **Service Layer** (`internal/core/service/`): business logic (e.g., `GarminService`, `AuthService`)
- **Repository Layer** (`internal/repository/gorm/`): data access implementations
- **Handler Layer** (`internal/handler/`): HTTP request handling (Fiber framework)
- **Infrastructure** (`internal/infrastructure/`): databases, validators, metrics

**Multiple services** are defined under [`cmd/`](cmd/):
- [`cmd/server/`](cmd/server/) - Main Garmin service (port 5051)
- [`cmd/server_auth/`](cmd/server_auth/) - Authentication service
- [`cmd/server_exporter/`](cmd/server_exporter/) - Metrics exporter

## Dependency Injection: Uber FX

All services use **Uber FX** for dependency injection. See [`cmd/server/main.go`](cmd/server/main.go):

```go
app := fx.New(
    fx.Provide(
        common.Load,           // Config loads first
        config.NewZapLogger,
        database.NewGarminGormDB,
        service.NewGarminService,
        // ...
    ),
    fx.Invoke(route.RegisterGarminRoutes),
)
```

### Multiple Database Pattern

Use **`fx.Annotate`** with result/param tags to provide multiple DB instances:

```go
fx.Annotate(database.NewGarminGormDB, fx.ResultTags(`name:"GarminDB"`)),
fx.Annotate(gorm.NewGarminRepo, fx.ParamTags(`name:"GarminDB"`)),
```

This allows separate databases for Garmin, Auth, and PLN Mobile data. Always use tags when a new database connection is needed.

## Response Handling Pattern

Use the centralized [`Handler`](internal/common/handler.go) for consistent API responses:

```go
// Success
return h.handler.ResponseSuccess(c, data)

// Validation error
return h.handler.ResponseValidationError(c, constant.VALIDATION_ERROR, errors)

// Error with custom status
return h.handler.ResponseError(c, http.StatusBadRequest, constant.INVALID_REQUEST_BODY, err)
```

Every handler function should follow this pattern. See [`internal/handler/garmin_handler.go`](internal/handler/garmin_handler.go) for examples.

## Testing Strategy

Tests use **testify/suite** with **sqlmock** for database mocking. See [`internal/core/service/auth_service_test.go`](internal/core/service/auth_service_test.go):

```go
type AuthServiceIntegrationSuite struct {
    suite.Suite
    db      *gorm.DB
    mock    sqlmock.Sqlmock
    service *service.AuthService
}

func (s *AuthServiceIntegrationSuite) SetupTest() {
    db, mock, cleanup := setupMockDB(s.T())
    // ...
}
```

**Run tests**: `go test ./... -coverprofile=coverage.out`

## Configuration Management

Environment config uses **Viper** ([`internal/common/config.go`](internal/common/config.go)):
- All config fields use `mapstructure` tags
- Multiple database configs: `GarminDB*`, `AuthDB*`, `PostgresPlnMobile*`
- Redis, RabbitMQ, JWT, OAuth configs included
- `common.Load` should always be the **first provider** in FX

## Docker Development Workflow

### Quick Commands

```bash
make dev           # Start development with tools (PgAdmin, Redis Commander)
make prod          # Start production environment
make logs          # View logs
make db-shell      # PostgreSQL shell
make redis-shell   # Redis CLI
make backup        # Database backup
```

Read [`README_DOCKER.md`](README_DOCKER.md) for detailed Docker configurations.

### Docker Compose Profiles

- **Base**: `docker-compose.yml` (core services only)
- **Development**: `docker-compose.dev.yml` (adds PgAdmin, Redis Commander)
- **Production**: `docker-compose.prod.yml` (production configs)

Services: `garmin-service`, `postgres`, `redis`, optional `pgadmin`, `redis-commander`

## Service Structure

Each service follows this pattern:

1. **Domain** defines repository interface
2. **Service** implements business logic using repository
3. **Repository** implements domain interface with GORM
4. **Handler** receives HTTP requests, calls service
5. **Route** registers handler with Fiber app

**Example** (Garmin flow):
[`domain/garmin.go`](internal/core/domain/garmin.go) → [`service/garmin_service.go`](internal/core/service/garmin_service.go) → [`repository/gorm/garmin_repo.go`](internal/repository/gorm/garmin_repo.go) → [`handler/garmin_handler.go`](internal/handler/garmin_handler.go) → [`route/garmin_route.go`](internal/route/garmin_route.go)

## Key Dependencies

- **Web Framework**: Fiber v2
- **ORM**: GORM with PostgreSQL driver
- **DI**: Uber FX
- **Logging**: Uber Zap
- **Config**: Viper
- **Validation**: go-playground/validator
- **Testing**: testify/suite + sqlmock
- **Caching**: Redis (go-redis/v9)
- **Tracing**: DataDog (`dd-trace-go`)
- **Metrics**: Prometheus client

## Conventions

- **File naming**: `{feature}_handler.go`, `{feature}_service.go`, `{feature}_repo.go`, `{feature}_route.go`
- **Package naming**: match directory name (e.g., `package handler`)
- **Error logging**: use `zap.Error(err)` with descriptive message keys
- **Middleware**: registered in [`internal/middleware/`](internal/middleware/)
- **Swagger docs**: generated with `swag init` (annotations in handlers)
- **Health checks**: implement via service layer, exposed at `/health-check`

## Building & Running

### Local Development

```bash
go run ./cmd/server          # Run main server
go run ./cmd/server_auth     # Run auth server
```

### With Docker

```bash
make dev                     # Development mode
docker-compose logs -f       # Follow logs
```

### Swagger Documentation

Access at `http://localhost:5051/swagger/` when server is running.

## Observability

- **Metrics**: Exposed at `/metrics` endpoint (Prometheus format)
- **Tracing**: DataDog APM integrated (configure via env vars: `DD_SERVICE`, `DD_ENV`, `DD_VERSION`)
- **Logging**: Structured JSON logs with Zap (configurable log level)
- **Health**: `/health-check` endpoint checks Redis + Postgres connectivity

## Common Patterns to Follow

1. **Always inject dependencies** via constructor functions compatible with FX
2. **Log liberally** at appropriate levels (Info for events, Error for failures, Debug for details)
3. **Use contexts** for cancellation and timeouts (pass through service → repo)
4. **Validate input** using validator struct via `h.handler.Validator.Struct(request)`
5. **Handle errors gracefully** with appropriate HTTP status codes
6. **Use GORM upserts** with `clause.OnConflict` for idempotent operations

## Future Additions (Planned)

- [ ] `securego/gosec` for static security analysis
- [ ] OpenTelemetry for distributed tracing
- [ ] Kratos integration (placeholder at [`examples/kratos-integration/`](examples/kratos-integration/))
