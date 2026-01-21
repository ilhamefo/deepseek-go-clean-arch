# Docker Compose Quick Start

## 🚀 Quick Commands

```bash
# Development (recommended for local development)
make dev
# or
./docker.sh up

# Production
make prod
# or  
./docker.sh -e prod up

# View logs
make logs
# or
./docker.sh logs

# Stop everything
make down
# or
./docker.sh down
```

## 📋 What's Included

### Core Services
- **auth-service** (Go app) - Port 5051
  - Swagger UI: http://localhost:5051/swagger/
- **postgres** (Database) - Port 5432
- **redis** (Cache/Sessions) - Port 6379

### Development Tools (dev environment)
- **pgadmin** - http://localhost:8080 (admin@example.com / admin123)
- **redis-commander** - http://localhost:8081

### Monitoring (optional)
- **prometheus** - http://localhost:9090
- **grafana** - http://localhost:3000 (admin / admin123)

## 🎯 Environment Profiles

| Profile | Command | Includes |
|---------|---------|----------|
| **Minimal** | `docker-compose up -d` | Core services only |
| **Development** | `make dev` | Core + dev tools |
| **Production** | `make prod` | Core with prod config |
| **Monitoring** | `make monitoring` | Core + monitoring stack |

## 🔧 Environment Setup

1. **Copy environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit configuration:**
   ```bash
   # Database
   POSTGRES_PASSWORD=your-secure-password
   
   # Redis  
   REDIS_PASSWORD=your-redis-password
   
   # Application secrets
   JWT_SECRET=your-jwt-secret
   GARMIN_CLIENT_ID=your-garmin-client-id
   GARMIN_CLIENT_SECRET=your-garmin-client-secret
   ```

## 🔍 Useful Commands

### Using Makefile (recommended)
```bash
make help           # Show all commands
make dev            # Start development
make logs           # View logs  
make status         # Check health
make db-shell       # PostgreSQL shell
make redis-shell    # Redis shell
make backup         # Database backup
make clean          # Clean up everything
```

### Using docker.sh script
```bash
./docker.sh -h                    # Help
./docker.sh up                    # Start dev
./docker.sh -e prod up            # Start production
./docker.sh -p development up     # With dev tools
./docker.sh logs auth-service     # Service logs
./docker.sh shell postgres        # Container shell
./docker.sh backup                # Database backup
```

### Direct Docker Compose
```bash
# Development
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Production  
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# With profiles
docker-compose --profile development up -d
docker-compose --profile monitoring up -d
```

## 🏥 Health Checks

All services include health checks:
- Services start only when dependencies are healthy
- Health status visible in `docker-compose ps`
- Automatic restart on failure

## 📊 Monitoring & Logs

```bash
# Service status
make status

# All logs
make logs

# Specific service logs
make logs SERVICE=auth-service

# Follow logs
docker-compose logs -f auth-service
```

## 💾 Database Operations

```bash
# Connect to database
make db-shell

# Backup database
make backup

# Restore database
make restore FILE=backup_20251009_150000.sql.gz

# View database with PgAdmin
make pgadmin  # Opens http://localhost:8080
```

## 🔄 Development Workflow

```bash
# 1. Start development environment
make dev

# 2. Code changes (auto-detected if using air/live reload)

# 3. Rebuild specific service
make rebuild

# 4. View logs
make app-logs

# 5. Run tests
make test
```

## 🚨 Troubleshooting

### Port Conflicts
```bash
# Check what's using ports
netstat -tulpn | grep :5051
netstat -tulpn | grep :5432
netstat -tulpn | grep :6379
```

### Service Won't Start
```bash
# Check service logs
make logs SERVICE=postgres

# Check health status
make status

# Restart specific service
make restart SERVICE=redis
```

### Database Issues
```bash
# Reset database (⚠️ deletes data)
docker-compose down postgres
docker volume rm deepseek-go-clean-arch_postgres_data
make dev
```

### Clean Start
```bash
# Remove everything and start fresh
make clean
make dev
```

## 🔒 Production Deployment

1. **Set production environment:**
   ```bash
   cp .env.production .env
   # Edit with production values
   ```

2. **Deploy:**
   ```bash
   make prod
   ```

3. **Monitor:**
   ```bash
   make status
   make logs
   ```

4. **Backup:**
   ```bash
   # Setup automated backups
   crontab -e
   # Add: 0 2 * * * cd /path/to/project && make backup
   ```

## 📁 File Structure

```
├── docker-compose.yml           # Main compose file
├── docker-compose.dev.yml       # Development overrides  
├── docker-compose.prod.yml      # Production overrides
├── .env.example                 # Environment template
├── .env.production              # Production template
├── Makefile                     # Make commands
├── docker.sh                    # Management script
├── DOCKER_COMPOSE.md            # Detailed documentation
└── scripts/
    ├── postgres-init.sh         # DB initialization
    ├── redis.conf               # Redis configuration
    ├── pgadmin-servers.json     # PgAdmin setup
    └── dev-data.sql             # Development test data
```

## 🎨 Customization

### Adding New Services
1. Add service to `docker-compose.yml`
2. Add development overrides to `docker-compose.dev.yml`
3. Add production overrides to `docker-compose.prod.yml`
4. Update documentation

### Custom Networks
- Services use `auth-network` (172.20.0.0/16)
- Modify network configuration in compose files

### Resource Limits
- Production limits configured in `docker-compose.prod.yml`
- Adjust based on your server capacity

## 📞 Support

For issues:
1. Check `make status` and `make logs`
2. Review troubleshooting section
3. Check Docker and Docker Compose versions
4. Refer to `DOCKER_COMPOSE.md` for detailed docs