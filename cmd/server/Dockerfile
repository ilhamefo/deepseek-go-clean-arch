# syntax=docker/dockerfile:1

# Build arguments
ARG GO_VERSION=1.23
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG SERVICE_NAME=auth-service

# Build stage
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /src

# Build arguments for cross-compilation
ARG TARGETOS TARGETARCH

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy source code
COPY . .

# Generate swagger docs (if swag is available, otherwise skip)
RUN if command -v swag >/dev/null 2>&1; then \
        swag init -g cmd/server/main.go -o docs; \
    else \
        echo "Swagger docs generation skipped - swag not available"; \
        mkdir -p docs && echo 'package docs' > docs/docs.go; \
    fi

# Build the application
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -buildvcs=false \
    -mod=readonly \
    -trimpath \
    -ldflags "-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev') -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /app/server \
    ./cmd/server

# Test if binary was built correctly
RUN test -f /app/server && echo "Binary built successfully" || exit 1

# Runtime stage - using distroless for CA certificates support (Garmin API needs HTTPS)
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

# Build arguments for labels
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

# Metadata labels
LABEL org.opencontainers.image.title="Auth Service Go" \
      org.opencontainers.image.description="Event Registration Auth Service with Garmin Integration" \
      org.opencontainers.image.source="https://github.com/ilhamefo/deepseek-go-clean-arch" \
      org.opencontainers.image.vendor="ilhamefo" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${GIT_COMMIT}" \
      org.opencontainers.image.created="${BUILD_DATE}"

# Build arguments for runtime
ARG SERVICE_NAME=auth-service
ARG PORT=5051

# Set environment variables
ENV PORT=${PORT} \
    GIN_MODE=release \
    GO_ENV=production

# Copy the binary from builder stage
COPY --from=builder /app/server /app

# Create necessary directories for logs if needed
USER 65532:65532

# Expose port
EXPOSE ${PORT}

# Health check endpoint (adjust if your app has different health endpoint)
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/app", "--health"] || exit 1

# Set the entrypoint
ENTRYPOINT ["/app"]
