# Docker Compose Setup

## Overview

This Docker Compose setup provides a complete development and production environment for the Auth Service Go application with Garmin integration.

## Services

### Core Services
- **auth-service**: Main Go application (port 5051)
- **postgres**: PostgreSQL database (port 5432)
- **redis**: Redis cache/session store (port 6379)

### Development Tools (profile: development)
- **pgadmin**: PostgreSQL admin interface (port 8080)
- **redis-commander**: Redis management interface (port 8081)

### Monitoring Tools (profile: monitoring)
- **prometheus**: Metrics collection (port 9090)
- **grafana**: Monitoring dashboards (port 3000)

## Quick Start

### Development Environment

1. **Setup environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Start all services:**
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
   ```

3. **Access services:**
   - Application: http://localhost:5051
   - Swagger UI: http://localhost:5051/swagger/
   - PgAdmin: http://localhost:8080 (admin@example.com / admin123)
   - Redis Commander: http://localhost:8081

### Production Environment

1. **Setup production environment:**
   ```bash
   cp .env.production .env
   # Set production values for all variables
   ```

2. **Deploy:**
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
   ```

## Environment Profiles

### Default (Minimal)
```bash
docker-compose up -d
```
Starts: auth-service, postgres, redis

### Development
```bash
docker-compose --profile development up -d
# or
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```
Adds: pgadmin, redis-commander, development configuration

### Monitoring
```bash
docker-compose --profile monitoring up -d
```
Adds: prometheus, grafana

### All Services
```bash
docker-compose --profile development --profile monitoring up -d
```

## Configuration

### Environment Variables

Key environment variables in `.env`:

```bash
# Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your-secure-password
POSTGRES_DB=garmin_db

# Redis
REDIS_PASSWORD=your-redis-password

# Application
VERSION=v1.0.0
APP_SECRET_KEY=your-secret-key
JWT_SECRET=your-jwt-secret
GARMIN_CLIENT_ID=your-garmin-client-id
GARMIN_CLIENT_SECRET=your-garmin-client-secret
```

### Volume Mounts

Persistent data:
- `postgres_data`: PostgreSQL database files
- `redis_data`: Redis persistence files
- `grafana_data`: Grafana dashboards and config
- `prometheus_data`: Prometheus metrics storage

Development mounts:
- `./logs`: Application logs
- `./garmin.sql`: Database initialization

## Health Checks

All services include health checks:
- **auth-service**: Custom health endpoint
- **postgres**: `pg_isready` command
- **redis**: `redis-cli ping`

## Networking

All services communicate via the `auth-network` bridge network with subnet `172.20.0.0/16`.

## Resource Limits

Production deployment includes resource limits:
- **auth-service**: 512MB RAM, 0.5 CPU
- **postgres**: 1GB RAM, 1.0 CPU  
- **redis**: 512MB RAM, 0.5 CPU

## Logging

Production logging configuration:
- JSON format
- 10MB max file size
- 3 file rotation

## Common Commands

### Build and Start
```bash
# Development
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# Production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build -d
```

### Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f auth-service

# Last 100 lines
docker-compose logs --tail=100 -f
```

### Database Operations
```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U postgres -d garmin_db

# Backup database
docker-compose exec postgres pg_dump -U postgres garmin_db > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres -d garmin_db < backup.sql
```

### Redis Operations
```bash
# Connect to Redis
docker-compose exec redis redis-cli

# Monitor Redis
docker-compose exec redis redis-cli monitor
```

### Scaling Services
```bash
# Scale auth service to 3 replicas
docker-compose up -d --scale auth-service=3
```

### Cleanup
```bash
# Stop and remove containers
docker-compose down

# Remove volumes (⚠️ deletes data)
docker-compose down -v

# Remove images
docker-compose down --rmi all
```

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 5051, 5432, 6379 are not in use
2. **Permission errors**: Check file permissions on mounted volumes
3. **Database connection**: Verify `depends_on` and health checks
4. **Memory issues**: Increase Docker memory limits if needed

### Debug Mode

Enable debug logging:
```bash
# Add to .env
GO_ENV=development
LOG_LEVEL=debug
```

### Service Status
```bash
# Check all services
docker-compose ps

# Check specific service health
docker-compose exec auth-service /app --health
```

## Security Considerations

### Production Security

1. **Change default passwords** in `.env.production`
2. **Use secrets management** for sensitive values
3. **Enable TLS** for external connections
4. **Limit network exposure** (remove port mappings for internal services)
5. **Regular updates** of base images
6. **Backup strategy** for persistent data

### Network Security

For production, consider:
- Using external networks
- Implementing reverse proxy (Traefik/Nginx)
- SSL/TLS termination
- Rate limiting

## Monitoring

### Metrics Endpoints

- Application metrics: `http://localhost:5051/metrics`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000`

### Alerting

Configure alerts in Prometheus for:
- Service availability
- Response times
- Database connections
- Memory/CPU usage

## Backup Strategy

### Automated Backups

Create backup scripts:
```bash
#!/bin/bash
# backup.sh
docker-compose exec postgres pg_dump -U postgres garmin_db | gzip > "backup_$(date +%Y%m%d_%H%M%S).sql.gz"
```

### Recovery Testing

Regularly test backup restoration:
```bash
# Test restore
gunzip -c backup_20251009_150000.sql.gz | docker-compose exec -T postgres psql -U postgres -d garmin_db_test
```