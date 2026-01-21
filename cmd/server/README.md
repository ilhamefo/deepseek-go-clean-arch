# Docker Build Guide

## Overview

This directory contains optimized Dockerfile configurations for the Auth Service Go application with the following characteristics:

- **Multi-stage build** for minimal final image size
- **BuildKit cache optimization** for faster builds
- **Security-focused** with non-root user execution
- **Cross-platform support** (linux/amd64, linux/arm64)
- **Reproducible builds** with version information

## Files

- `Dockerfile` - Main Dockerfile (distroless runtime for HTTPS/CA support)
- `Dockerfile.cgo` - Alternative for CGO-dependent applications
- `build.sh` - Build script with various options
- `../../../.dockerignore` - Files to exclude from build context

## Image Sizes

- **Current build**: ~30-35MB (distroless with CA certificates)
- **Alternative (scratch)**: ~15-25MB (if no HTTPS outbound needed)
- **CGO build**: ~35-45MB (distroless/cc runtime)

## Quick Start

### Build Single Platform

```bash
# Basic build
docker buildx build -t auth-service:latest -f cmd/server/Dockerfile .

# With custom arguments
docker buildx build \
  --build-arg GO_VERSION=1.23 \
  --build-arg SERVICE_NAME=auth-service \
  --build-arg PORT=5051 \
  -t auth-service:latest \
  -f cmd/server/Dockerfile .
```

### Build Multi-Platform

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --build-arg GO_VERSION=1.23 \
  --build-arg SERVICE_NAME=auth-service \
  --build-arg PORT=5051 \
  -t auth-service:latest \
  -f cmd/server/Dockerfile \
  --push .
```

### Using Build Script

```bash
# Make executable
chmod +x build.sh

# Build for single platform
./build.sh single

# Build for multiple platforms
./build.sh multi

# Test the built image
./build.sh test

# Check image size
./build.sh size
```

## Build Arguments

| Argument | Default | Description |
|----------|---------|-------------|
| `GO_VERSION` | `1.23` | Go version for build |
| `TARGETOS` | `linux` | Target operating system |
| `TARGETARCH` | `amd64` | Target architecture |
| `SERVICE_NAME` | `auth-service` | Binary name |
| `PORT` | `5051` | Application port |

## Runtime Options

### Standard Build (Recommended)

Uses `gcr.io/distroless/static-debian12:nonroot`:
- ✅ CA certificates included (for HTTPS outbound)
- ✅ Timezone data included
- ✅ Non-root user (uid=65532)
- ✅ No shell or package manager
- 📏 Size: ~30-35MB

### CGO Build

Uses `gcr.io/distroless/cc-debian12:nonroot`:
- ✅ C library support
- ✅ CA certificates and timezone data
- ✅ Non-root user
- 📏 Size: ~35-45MB

Use when your application requires:
- Database drivers with C dependencies
- C libraries integration
- Native extensions

## Running the Container

### Basic Run

```bash
docker run -p 5051:5051 auth-service:latest
```

### With Environment Variables

```bash
docker run \
  -p 5051:5051 \
  -e PORT=5051 \
  -e GO_ENV=production \
  auth-service:latest
```

### With Volume Mounts (if needed)

```bash
docker run \
  -p 5051:5051 \
  -v $(pwd)/config:/config:ro \
  auth-service:latest
```

## Health Check

The container includes a health check that runs every 30 seconds:

```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/app", "--health"] || exit 1
```

Make sure your application supports the `--health` flag for proper health checking.

## Security Features

1. **Non-root execution**: Runs as user 65532 (nobody)
2. **Minimal attack surface**: Distroless base with no shell
3. **Read-only filesystem**: No write access to system files
4. **Reproducible builds**: Version information embedded

## Optimization Features

1. **Layer caching**: Dependencies downloaded separately from source
2. **BuildKit cache mounts**: Persistent caching across builds
3. **Binary stripping**: `-ldflags "-s -w"` removes debug info
4. **Path trimming**: `-trimpath` for reproducible paths
5. **Static linking**: No external dependencies

## Troubleshooting

### Build Issues

1. **CGO errors**: Use `Dockerfile.cgo` instead
2. **Network timeouts**: Check if behind corporate firewall
3. **Cache issues**: Use `--no-cache` flag to rebuild

### Runtime Issues

1. **Permission errors**: Ensure files are readable by user 65532
2. **Port binding**: Check if port is already in use
3. **Health check fails**: Implement `--health` flag in your application

### Size Issues

1. **Larger than expected**: Check if unnecessary files are being copied
2. **Missing dependencies**: May need to switch to distroless/cc for CGO

## Development Workflow

1. **Local development**:
   ```bash
   go run cmd/server/main.go
   ```

2. **Build and test**:
   ```bash
   ./build.sh single
   ./build.sh test
   ```

3. **Multi-platform release**:
   ```bash
   ./build.sh multi
   ```

## Integration with CI/CD

### GitHub Actions Example

```yaml
- name: Build and push Docker image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: cmd/server/Dockerfile
    platforms: linux/amd64,linux/arm64
    push: true
    tags: |
      your-registry/auth-service:latest
      your-registry/auth-service:${{ github.sha }}
    build-args: |
      GO_VERSION=1.23
      SERVICE_NAME=auth-service
      PORT=5051
    cache-from: type=gha
    cache-to: type=gha,mode=max
```