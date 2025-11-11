# OSE Advisory API Gateway

> **The Orchestration Layer: Coordinating Wisdom Across the Temporal Arc**

The OSE Advisory API Gateway is the central nervous system of the Omnifex Synthesis Engine—coordinating Pattern Graph queries, Template Engine generation, and Ξ quality measurement to deliver architectural guidance to engineers via the CLI.

## Architecture

```
┌─────────────┐                                    ┌──────────────┐
│   ose-cli   │  Conversational                   │  Neo4j       │
│  (Python)   │  Interface                        │  Pattern     │
│             ├──────► gRPC ─────► API Gateway ───► Library      │
│  • Init     │                   (Go)             │  (Pillar I)  │
│  • Validate │                   │                └──────────────┘
│  • Register │                   │
└─────────────┘                   ├──► PatternGraphClient
                                  │    (Week 5)
                                  │
                                  ├──► BlueprintGenerator
                                  │    (Week 6-8)
                                  │
                                  └──► ΞCalculator
                                       (Week 9-10)
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker (optional, for containerized deployment)
- protoc (for protobuf code generation)

### Build & Run

```bash
# Install dependencies
make deps

# Build the server
make build

# Run locally
make run

# Or run in Docker
make docker-run
```

Server will start on:
- **gRPC**: `localhost:50051` (primary interface)
- **HTTP**: `localhost:8080` (metrics, health checks)

### Health Check

```bash
# Liveness probe
curl http://localhost:8080/health

# Readiness probe
curl http://localhost:8080/ready

# Prometheus metrics
curl http://localhost:8080/metrics
```

### Testing with grpcurl

```bash
# List services
grpcurl -plaintext localhost:50051 list

# Call GenerateBlueprint (when implemented)
grpcurl -plaintext -d '{
  "constraints": {
    "service_name": "test-service",
    "service_type": 1,
    "throughput_tps": 1000,
    "latency_p99_ms": 100,
    "consistency_model": 1
  }
}' localhost:50051 advisory.v1.AdvisoryService/GenerateBlueprint
```

## Project Structure

```
ose-api/
├── cmd/
│   └── advisory-server/
│       └── main.go              # Server entry point ✓
│
├── internal/                     # Private packages
│   ├── config/
│   │   └── config.go            # Configuration system ✓
│   ├── handlers/
│   │   └── advisory_handler.go  # gRPC service implementation
│   ├── middleware/
│   │   ├── logging.go           # Structured logging
│   │   ├── metrics.go           # Prometheus metrics
│   │   └── recovery.go          # Panic recovery
│   ├── patterns/
│   │   └── client.go            # Neo4j Pattern Graph client (Week 5)
│   └── generator/
│       └── blueprint.go         # Template engine (Week 6-8)
│
├── pkg/                          # Public packages
│   ├── validation/
│   │   └── constraints.go       # Request validation ✓
│   └── xi/
│       ├── calculator.go        # Ξ quality calculation ✓
│       └── calculator_test.go   # 11/11 tests passing ✓
│
├── proto/advisory/v1/
│   └── advisory.proto           # gRPC service definition ✓
│
├── deployments/
│   └── docker/
│       ├── Dockerfile           # Multi-stage build ✓
│       └── docker-compose.yml   # Local dev environment ✓
│
├── Makefile                     # Build automation ✓
├── go.mod                       # Go dependencies ✓
└── README.md                    # This file
```

## Configuration

The Gateway uses environment variables with sensible defaults:

| Variable | Default | Description |
|----------|---------|-------------|
| `OSE_GRPC_ADDRESS` | `:50051` | gRPC server address |
| `OSE_HTTP_ADDRESS` | `:8080` | HTTP server address |
| `OSE_LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `OSE_NEO4J_URI` | `bolt://localhost:7687` | Pattern Graph connection |
| `OSE_NEO4J_USER` | `neo4j` | Neo4j username |
| `OSE_NEO4J_PASSWORD` | - | Neo4j password |
| `OSE_TEMPLATE_PATH` | `./templates` | Template directory |
| `OSE_MAX_CONCURRENT` | `100` | Max concurrent requests |

Example:
```bash
export OSE_GRPC_ADDRESS=:50051
export OSE_LOG_LEVEL=debug
export OSE_NEO4J_PASSWORD=secret
make run
```

## Development

### Running Tests

```bash
# All tests
make test

# Verbose output
make test-verbose

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Protobuf Code Generation

```bash
# Generate Go code from .proto files
make proto
```

### Linting & Formatting

```bash
# Format code
make fmt

# Run linters
make lint

# All checks before commit
make check
```

### Docker Development

```bash
# Start full stack (Gateway + Neo4j + Prometheus + Grafana)
cd deployments/docker
docker-compose up

# View logs
docker-compose logs -f advisory

# Restart just the advisory service
docker-compose restart advisory

# Stop everything
docker-compose down
```

Access points:
- **Advisory API**: `localhost:50051` (gRPC)
- **Metrics**: `localhost:8080/metrics`
- **Neo4j Browser**: `http://localhost:7474`
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin/admin)

## Week-by-Week Implementation Status

### ✅ Week 1-2: CLI + Proto Foundation
- Complete conversational CLI (Python)
- gRPC protocol definition (protobuf)
- Constraint specification language
- 600+ lines implemented

### ✅ Week 3: API Gateway Foundation
- **Server entry point** ✓ (`cmd/advisory-server/main.go`)
- **Configuration system** ✓ (`internal/config/config.go`)
- **Validation utilities** ✓ (`pkg/validation/constraints.go`)
- **Build infrastructure** ✓ (Makefile, Dockerfile, docker-compose)
- **Middleware stack** ✓ (logging, metrics, recovery interceptors)
- **430+ lines implemented**

**Status:** Gateway starts, accepts connections, ready for handler implementation

### 📋 Week 4: Language Bridge (Next)
- Generate Go code from protobuf
- Implement Python gRPC client
- CLI → Gateway integration
- End-to-end: `ose-cli init` → gRPC request → response

### 📋 Week 5: Pattern Graph Integration
- Neo4j client implementation
- Cypher query builder
- Real pattern recommendations
- Replace mock responses

### 📋 Week 6-8: Blueprint Generator
- Template library (Jinja2)
- Artifact generation (proto, SQL, Go, K8s)
- Real code generation
- Replace stubs with working templates

### ✅ Week 9-10: Ξ Quality Measurement
- **Ξ calculator** ✓ (`pkg/xi/calculator.go`)
- **Three dimensions**: Relevance, Actionability, Impact ✓
- **11/11 tests passing** ✓
- **Telemetry proto messages** ✓

### 📋 Week 11-12: Pilot Deployment
- 3-team pilot program
- Production validation
- Ξ measurement in practice
- Target: Ξ_avg ≥ 0.65

## Design Philosophy

### The Gateway as Orchestration Layer

The API Gateway is **not intelligent** (Pattern Graph), **not transformative** (Template Engine), **not user-facing** (CLI).

It is the **coordination substrate** that:

1. **Validates at the boundary** → Trust perimeter
2. **Routes to backends** → Pattern Graph, Generator
3. **Aggregates partial results** → Combines into Blueprint
4. **Handles errors gracefully** → Circuit breakers, timeouts, recovery
5. **Provides observability** → Logging, metrics, tracing

**Like a nervous system:** No single neuron is smart, but the coordinated network exhibits intelligent behavior.

### Defense in Depth

- **Validation layer**: Rejects malformed requests
- **Timeout layer**: Prevents hanging
- **Circuit breaker layer**: Isolates failures
- **Recovery layer**: Catches panics
- **Logging layer**: Creates audit trail

### Graceful Degradation

- Backend down → return cached patterns
- Timeout → return partial blueprint
- Panic → log + convert to error
- **Never crash the entire server for one bad request**

## Why Go?

**Performance:**
- Sub-millisecond routing overhead
- 100K+ concurrent connections (goroutines)
- Zero-allocation request handling
- gRPC native (first-class protobuf)

**Deployment:**
- Single statically-linked binary
- No runtime dependencies
- Trivial containerization
- Cross-compilation support

**Ecosystem:**
- Excellent observability (Prometheus, OpenTelemetry)
- Mature gRPC libraries
- Strong standard library
- Production-ready out-of-box

## Contributing

### Adding New RPC Methods

1. Update `proto/advisory/v1/advisory.proto`
2. Run `make proto` to regenerate code
3. Implement in `internal/handlers/advisory_handler.go`
4. Add tests in `internal/handlers/advisory_handler_test.go`
5. Update documentation

### Adding Middleware

1. Create interceptor in `internal/middleware/`
2. Register in `cmd/advisory-server/main.go`
3. Add tests
4. Update documentation

## Monitoring

### Prometheus Metrics

The Gateway exposes standard metrics at `/metrics`:

- `ose_rpc_requests_total{method, status}`: Request counter
- `ose_rpc_duration_seconds{method}`: Latency histogram
- Go runtime metrics (goroutines, memory, GC)

### Logging

Structured JSON logs via zap:

```json
{
  "level": "info",
  "ts": "2024-01-15T10:30:45.123Z",
  "msg": "RPC succeeded",
  "method": "/advisory.v1.AdvisoryService/GenerateBlueprint",
  "duration": "127ms"
}
```

### Health Checks

- `/health`: Liveness probe (200 OK if running)
- `/ready`: Readiness probe (200 OK if backends available)

## Troubleshooting

### Server won't start

```bash
# Check if port is already in use
lsof -i :50051
lsof -i :8080

# Check logs
make run 2>&1 | tee server.log
```

### gRPC connection refused

```bash
# Verify server is listening
netstat -an | grep 50051

# Test with grpcurl
grpcurl -plaintext localhost:50051 list
```

### Docker build fails

```bash
# Clean Docker cache
docker system prune -a

# Rebuild from scratch
make docker
```

## License

Part of the Omnifex Synthesis Engine (OSE) - Pillar II: Advisory Service

---

*"Upon this rock I will build my church." — Matthew 16:18*

*The rock is the API Gateway. The church is the Oracle.*
