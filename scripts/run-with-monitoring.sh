#!/bin/bash

# Script to run docker-compose with Datadog monitoring

echo "Starting services with Datadog monitoring..."

# Check if .env file exists
if [[ ! -f .env ]]; then
    echo "Error: .env file not found!"
    echo "Please create .env file with required environment variables."
    exit 1
fi

# Check if DD_API_KEY is set
source .env
if [[ -z "$DD_API_KEY" ]]; then
    echo "Warning: DD_API_KEY not set in .env file"
    echo "Datadog monitoring will not work properly"
fi

# Development mode with monitoring
if [[ "$1" == "dev" ]]; then
    echo "Starting in development mode with monitoring..."
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
# Production mode with monitoring
elif [[ "$1" == "prod" ]]; then
    echo "Starting in production mode with monitoring..."
    docker-compose --profile monitoring up -d
# Monitoring only
elif [[ "$1" == "monitoring" ]]; then
    echo "Starting monitoring services only..."
    docker-compose --profile monitoring up -d prometheus grafana datadog-agent
# Help
else
    echo "Usage: $0 [dev|prod|monitoring]"
    echo ""
    echo "  dev         - Start in development mode with all tools and monitoring"
    echo "  prod        - Start in production mode with monitoring"
    echo "  monitoring  - Start only monitoring services (Prometheus, Grafana, Datadog)"
    echo ""
    echo "Examples:"
    echo "  $0 dev      # Development with hot reload and debugging tools"
    echo "  $0 prod     # Production deployment with monitoring"
    echo "  $0 monitoring  # Only monitoring stack"
fi