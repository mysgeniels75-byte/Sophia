# Multi-stage Docker build for OSE Advisory Service
#
# Stage 1: Build the Go binary
# Stage 2: Create minimal runtime image
#
# Usage:
#   docker build -t ose-advisory:latest -f deployments/docker/Dockerfile .
#   docker run -p 50051:50051 -p 8080:8080 ose-advisory:latest

# ═══════════════════════════════════════════════════════════════════════
# STAGE 1: Build
# ═══════════════════════════════════════════════════════════════════════
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files and download dependencies (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with version information
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.GitCommit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -a -installsuffix cgo \
    -o advisory-server \
    cmd/advisory-server/main.go

# ═══════════════════════════════════════════════════════════════════════
# STAGE 2: Runtime
# ═══════════════════════════════════════════════════════════════════════
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 ose && \
    adduser -D -u 1000 -G ose ose

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/advisory-server .

# Create directory for templates (Week 6-8)
RUN mkdir -p /app/templates && chown -R ose:ose /app

# Switch to non-root user
USER ose

# Expose ports
EXPOSE 50051 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the server
ENTRYPOINT ["./advisory-server"]
