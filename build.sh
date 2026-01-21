#!/bin/bash
# Build script for the auth service

set -euo pipefail

# Configuration
SERVICE_NAME=${SERVICE_NAME:-auth-service}
GO_VERSION=${GO_VERSION:-1.23}
PORT=${PORT:-5051}
TAG=${TAG:-latest}
PLATFORMS=${PLATFORMS:-linux/amd64,linux/arm64}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# Check if docker buildx is available
check_buildx() {
    if ! docker buildx version >/dev/null 2>&1; then
        error "Docker buildx is not available. Please install Docker with buildx support."
    fi
    log "Docker buildx is available"
}

# Build single platform
build_single() {
    local platform=${1:-linux/amd64}
    log "Building for platform: $platform"
    
    docker buildx build \
        --platform "$platform" \
        --build-arg GO_VERSION="$GO_VERSION" \
        --build-arg SERVICE_NAME="$SERVICE_NAME" \
        --build-arg PORT="$PORT" \
        --build-arg VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
        --build-arg GIT_COMMIT="$(git rev-parse HEAD 2>/dev/null || echo 'unknown')" \
        --build-arg BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
        -t "${SERVICE_NAME}:${TAG}" \
        -t "${SERVICE_NAME}:${TAG}-$(echo $platform | tr '/' '-')" \
        -f cmd/server/Dockerfile \
        --load \
        .
}

# Build multi-platform
build_multi() {
    log "Building for multiple platforms: $PLATFORMS"
    
    docker buildx build \
        --platform "$PLATFORMS" \
        --build-arg GO_VERSION="$GO_VERSION" \
        --build-arg SERVICE_NAME="$SERVICE_NAME" \
        --build-arg PORT="$PORT" \
        --build-arg VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
        --build-arg GIT_COMMIT="$(git rev-parse HEAD 2>/dev/null || echo 'unknown')" \
        --build-arg BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
        -t "${SERVICE_NAME}:${TAG}" \
        -f cmd/server/Dockerfile \
        --push \
        .
}

# Test the built image
test_image() {
    local image="${SERVICE_NAME}:${TAG}"
    log "Testing image: $image"
    
    # Test if the image can start
    local container_id
    container_id=$(docker run -d -p "$PORT:$PORT" "$image")
    
    # Wait a bit for the service to start
    sleep 3
    
    # Check if container is still running
    if docker ps -q --filter "id=$container_id" | grep -q .; then
        log "Container started successfully"
        docker stop "$container_id" >/dev/null
        docker rm "$container_id" >/dev/null
    else
        error "Container failed to start"
    fi
}

# Show image size
show_size() {
    local image="${SERVICE_NAME}:${TAG}"
    local size
    size=$(docker images --format "table {{.Size}}" "$image" | tail -n +2)
    log "Image size: $size"
    
    # Check if size meets requirements (should be < 40MB for distroless)
    local size_mb
    size_mb=$(docker images --format "{{.Size}}" "$image" | sed 's/MB//' | head -1)
    if (( $(echo "$size_mb < 40" | bc -l) )); then
        log "✓ Image size meets requirements (< 40MB)"
    else
        warn "Image size is larger than expected (> 40MB)"
    fi
}

# Main function
main() {
    local command=${1:-single}
    
    log "Starting Docker build for $SERVICE_NAME"
    log "Go version: $GO_VERSION"
    log "Port: $PORT"
    log "Tag: $TAG"
    
    check_buildx
    
    case $command in
        single)
            build_single
            test_image
            show_size
            ;;
        multi)
            build_multi
            log "Multi-platform build completed"
            ;;
        test)
            test_image
            ;;
        size)
            show_size
            ;;
        *)
            echo "Usage: $0 [single|multi|test|size]"
            echo "  single - Build for single platform (default)"
            echo "  multi  - Build for multiple platforms and push"
            echo "  test   - Test the built image"
            echo "  size   - Show image size"
            exit 1
            ;;
    esac
    
    log "Build process completed successfully"
}

# Run main function with all arguments
main "$@"