#!/bin/bash
# Docker Compose management script

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}"
}

# Default values
ENVIRONMENT="dev"
COMPOSE_FILES=""
PROFILES=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--env)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -p|--profile)
            PROFILES="$PROFILES --profile $2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            COMMAND="$1"
            shift
            break
            ;;
    esac
done

show_help() {
    cat << EOF
Docker Compose Management Script

Usage: $0 [OPTIONS] COMMAND

Options:
  -e, --env ENV         Environment (dev|prod) [default: dev]
  -p, --profile PROFILE Add profile (development|monitoring)
  -h, --help           Show this help

Commands:
  up                   Start services
  down                 Stop services
  build                Build images
  logs                 Show logs
  status               Show service status
  clean                Clean up containers and volumes
  backup               Backup database
  restore FILE         Restore database from backup
  shell SERVICE        Open shell in service
  restart SERVICE      Restart specific service

Examples:
  $0 up                                    # Start dev environment
  $0 -e prod up                           # Start production environment  
  $0 -p development up                    # Start with development tools
  $0 -p development -p monitoring up      # Start with dev tools and monitoring
  $0 logs auth-service                    # Show auth service logs
  $0 shell postgres                       # Open shell in postgres container
  $0 backup                              # Backup database
  $0 restore backup_20251009.sql.gz      # Restore database

EOF
}

# Set compose files based on environment
set_compose_files() {
    COMPOSE_FILES="-f docker-compose.yml"
    
    case $ENVIRONMENT in
        dev|development)
            COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.dev.yml"
            ;;
        prod|production)
            COMPOSE_FILES="$COMPOSE_FILES -f docker-compose.prod.yml"
            ;;
        *)
            warn "Unknown environment: $ENVIRONMENT, using base configuration"
            ;;
    esac
}

# Check if .env file exists
check_env_file() {
    if [[ ! -f .env ]]; then
        warn ".env file not found"
        if [[ -f .env.example ]]; then
            info "Copying .env.example to .env"
            cp .env.example .env
            warn "Please edit .env file with your configuration"
        else
            error ".env.example not found. Please create .env file manually"
        fi
    fi
}

# Wait for services to be healthy
wait_for_services() {
    local services=("postgres" "redis")
    local max_wait=60
    local wait_time=0
    
    log "Waiting for services to be healthy..."
    
    for service in "${services[@]}"; do
        while [[ $wait_time -lt $max_wait ]]; do
            if docker-compose $COMPOSE_FILES ps $service | grep -q "healthy"; then
                log "$service is healthy"
                break
            else
                info "Waiting for $service to be healthy... (${wait_time}s/${max_wait}s)"
                sleep 5
                wait_time=$((wait_time + 5))
            fi
        done
        
        if [[ $wait_time -ge $max_wait ]]; then
            error "$service failed to become healthy within ${max_wait} seconds"
        fi
    done
}

# Commands
cmd_up() {
    check_env_file
    log "Starting services in $ENVIRONMENT environment"
    docker-compose $COMPOSE_FILES $PROFILES up -d "$@"
    wait_for_services
    log "All services started successfully"
    cmd_status
}

cmd_down() {
    log "Stopping services"
    docker-compose $COMPOSE_FILES down "$@"
    log "Services stopped"
}

cmd_build() {
    log "Building images"
    docker-compose $COMPOSE_FILES build "$@"
    log "Build completed"
}

cmd_logs() {
    if [[ $# -eq 0 ]]; then
        docker-compose $COMPOSE_FILES logs -f
    else
        docker-compose $COMPOSE_FILES logs -f "$@"
    fi
}

cmd_status() {
    info "Service Status:"
    docker-compose $COMPOSE_FILES ps
    
    info "\nService Health:"
    for service in auth-service postgres redis; do
        if docker-compose $COMPOSE_FILES ps $service | grep -q "Up"; then
            echo -e "  ${GREEN}✓${NC} $service"
        else
            echo -e "  ${RED}✗${NC} $service"
        fi
    done
}

cmd_clean() {
    warn "This will remove all containers and volumes. Are you sure? (y/N)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        log "Cleaning up containers and volumes"
        docker-compose $COMPOSE_FILES down -v --rmi local
        docker system prune -f
        log "Cleanup completed"
    else
        info "Cleanup cancelled"
    fi
}

cmd_backup() {
    local backup_file="backup_$(date +%Y%m%d_%H%M%S).sql.gz"
    log "Creating database backup: $backup_file"
    
    docker-compose $COMPOSE_FILES exec postgres pg_dump -U postgres garmin_db | gzip > "$backup_file"
    log "Backup created: $backup_file"
}

cmd_restore() {
    local backup_file="$1"
    if [[ ! -f "$backup_file" ]]; then
        error "Backup file not found: $backup_file"
    fi
    
    warn "This will restore database from $backup_file. Continue? (y/N)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        log "Restoring database from $backup_file"
        gunzip -c "$backup_file" | docker-compose $COMPOSE_FILES exec -T postgres psql -U postgres -d garmin_db
        log "Database restored successfully"
    else
        info "Restore cancelled"
    fi
}

cmd_shell() {
    local service="$1"
    log "Opening shell in $service"
    docker-compose $COMPOSE_FILES exec "$service" /bin/sh
}

cmd_restart() {
    local service="$1"
    log "Restarting $service"
    docker-compose $COMPOSE_FILES restart "$service"
    log "$service restarted"
}

# Main execution
set_compose_files

case "${COMMAND:-}" in
    up)
        cmd_up "$@"
        ;;
    down)
        cmd_down "$@"
        ;;
    build)
        cmd_build "$@"
        ;;
    logs)
        cmd_logs "$@"
        ;;
    status)
        cmd_status
        ;;
    clean)
        cmd_clean
        ;;
    backup)
        cmd_backup
        ;;
    restore)
        if [[ $# -eq 0 ]]; then
            error "Please provide backup file: $0 restore backup_file.sql.gz"
        fi
        cmd_restore "$1"
        ;;
    shell)
        if [[ $# -eq 0 ]]; then
            error "Please specify service: $0 shell postgres"
        fi
        cmd_shell "$1"
        ;;
    restart)
        if [[ $# -eq 0 ]]; then
            error "Please specify service: $0 restart auth-service"
        fi
        cmd_restart "$1"
        ;;
    "")
        error "No command specified. Use -h for help"
        ;;
    *)
        error "Unknown command: $COMMAND. Use -h for help"
        ;;
esac