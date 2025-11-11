# OSE Advisory API

> **The Protocol of Architectural Consultation**

The OSE Advisory API defines the contract between engineers seeking architectural guidance (via the CLI) and the Advisory Service that provides recommendations drawn from the Pattern Knowledge Graph.

## Philosophy

Protocol Buffers provide the standard unit of measurement for architectural consultation—the typed, versioned, language-agnostic contract that ensures the CLI's questions and the Service's answers speak the same language across all implementations.

Just as the meter was standardized by a platinum-iridium bar in Paris, enabling global collaboration on construction, the `advisory.proto` standardizes architectural consultation, enabling the CLI (Python), API Gateway (Go), and future tools (any language) to collaborate seamlessly.

## Quick Start

```bash
# Generate Go code from proto definitions
cd ose-api
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/advisory/v1/advisory.proto

# Generate Python code for CLI
python -m grpc_tools.protoc -I. \
    --python_out=. --grpc_python_out=. \
    proto/advisory/v1/advisory.proto
```

## Week 2 Status: Protocol Foundation

**What's Defined:**
- ✅ Complete gRPC service definition
- ✅ Constraint specification language
- ✅ Blueprint generation request/response
- ✅ Pattern search interface
- ✅ Service validation messages
- ✅ Learning feedback registration
- ✅ Ξ quality tracking metrics

**What's Coming:**
- Week 3-4: Go API Gateway implementation
- Week 5-6: Pattern Graph Client integration
- Week 7-8: Blueprint Generator service
- Week 9-10: Validation and registration endpoints

## Service Definition

### AdvisoryService

The core service interface with four primary RPCs:

```protobuf
service AdvisoryService {
  rpc GenerateBlueprint(GenerateBlueprintRequest) returns (GenerateBlueprintResponse);
  rpc SearchPatterns(SearchPatternsRequest) returns (SearchPatternsResponse);
  rpc ValidateService(ValidateServiceRequest) returns (ValidateServiceResponse);
  rpc RegisterService(RegisterServiceRequest) returns (RegisterServiceResponse);
}
```

#### 1. GenerateBlueprint

**Purpose**: The primary advisory function—convert constraints into recommendations.

**Request**:
```protobuf
message GenerateBlueprintRequest {
  ServiceConstraints constraints = 1;
  repeated ArtifactType requested_artifacts = 2;
}
```

**Response**:
```protobuf
message GenerateBlueprintResponse {
  Blueprint blueprint = 1;                // Recommended patterns + artifacts
  AdvisoryMetrics metrics = 2;            // For Ξ tracking
}
```

**Week 1-2**: CLI sends mock requests
**Week 3-4**: API Gateway receives and processes
**Week 5-6**: Queries Pattern Graph for real recommendations
**Week 7-8**: Generates actual code artifacts

#### 2. SearchPatterns

**Purpose**: Natural language or constraint-based pattern discovery.

**Request**:
```protobuf
message SearchPatternsRequest {
  oneof query_type {
    string semantic_query = 1;            // "handle backpressure"
    ServiceConstraints constraints = 2;   // Constraint-based filtering
  }
  int32 top_k = 3;                       // Number of results
}
```

**Use Cases**:
- Engineer asks: "How do I prevent database connection exhaustion?"
- CLI searches: `semantic_query = "database connection exhaustion"`
- Service returns: Connection pooling patterns with confidence scores

#### 3. ValidateService

**Purpose**: Anti-pattern detection in existing services.

**Request**:
```protobuf
message ValidateServiceRequest {
  string service_path = 1;
  repeated string files_to_analyze = 2;
}
```

**Response**:
```protobuf
message ValidateServiceResponse {
  repeated AntiPatternDetection anti_patterns_found = 1;
  repeated PatternViolation pattern_violations = 2;
  double overall_health_score = 3;
}
```

**Week 10**: Static analysis implementation

#### 4. RegisterService

**Purpose**: Close the learning loop—report pattern application outcomes.

**Request**:
```protobuf
message RegisterServiceRequest {
  string service_name = 1;
  repeated string patterns_applied = 2;
  ServiceConstraints constraints = 3;
  MetricsSnapshot metrics_before = 4;
  MetricsSnapshot metrics_after = 5;
}
```

**Response**:
```protobuf
message RegisterServiceResponse {
  string registration_id = 1;
  repeated PatternConfidenceUpdate confidence_updates = 3;
}
```

**Impact**: Feeds Meta-Learning Orchestrator, updates pattern confidence scores.

## The Constraint Language

`ServiceConstraints` is the structured vocabulary for expressing architectural requirements:

```protobuf
message ServiceConstraints {
  string service_name = 1;
  ServiceType service_type = 2;
  int32 throughput_tps = 3;
  int32 latency_p99_ms = 4;
  ConsistencyModel consistency_model = 5;
  repeated IntegrationType integrations = 6;
  int32 team_size = 7;
  double data_volume_gb = 8;
  DeploymentTarget deployment_target = 9;
  repeated string excluded_patterns = 10;
}
```

### Design Principles

**1. Each Field = Architectural Degree of Freedom**

Every field represents a dimension in the design space:
- `service_type`: Determines concurrency model
- `consistency_model`: Affects data architecture
- `throughput_tps + latency_p99_ms`: Defines performance envelope
- `integrations`: Influences interface patterns

**2. Type Safety**

Enums prevent invalid values:
```protobuf
enum ServiceType {
  SERVICE_TYPE_API = 1;              // ✓ Valid
  SERVICE_TYPE_INVALID = 999;        // ✗ Compiler error
}
```

**3. Extensibility**

Field numbers are stable. New fields can be added without breaking existing clients:
```protobuf
message ServiceConstraints {
  // ... existing fields 1-10 ...
  string regulatory_compliance = 11;  // Add in v1.1 without breaking v1.0
}
```

## Blueprint Structure

The `Blueprint` message is the advisory output—the Oracle's answer:

```protobuf
message Blueprint {
  string service_name = 1;
  repeated RecommendedPattern patterns = 2;
  repeated Artifact artifacts = 3;
  double confidence_score = 4;
  repeated TradeOff trade_offs = 5;
  repeated AlternativePattern alternatives_considered = 6;
  google.protobuf.Timestamp generated_at = 7;
}
```

### RecommendedPattern

Each pattern includes justification:

```protobuf
message RecommendedPattern {
  string pattern_id = 1;              // "pattern-001"
  string pattern_name = 2;            // "Actor Mailbox Backpressure"
  string category = 3;                // "concurrency"
  double confidence_score = 4;        // 0.95 (95% confidence)
  string rationale = 5;               // "Handles high load with graceful degradation"
  int32 application_count = 6;        // 12 services using this
  repeated string related_patterns = 7; // Often used with pattern-007
}
```

### Artifact

Generated code and configuration:

```protobuf
message Artifact {
  ArtifactType type = 1;              // PROTO, SQL, KUBERNETES, etc.
  string path = 2;                    // "proto/inventory/v1/service.proto"
  string content = 3;                 // Actual file content
  string template_used = 4;           // "templates/go-grpc-service.tmpl"
}
```

### TradeOff

Explicit cost-benefit analysis:

```protobuf
message TradeOff {
  string decision = 1;                // "Chose eventual consistency"
  string benefit = 2;                 // "10x higher throughput"
  string cost = 3;                    // "Temporary data inconsistency"
  string mitigation = 4;              // "Vector clock conflict resolution"
}
```

**Philosophy**: Engineers should understand what they're trading, not just what they're getting.

## Ξ Quality Tracking

The `AdvisoryMetrics` message enables quality measurement:

```protobuf
message AdvisoryMetrics {
  int32 patterns_recommended = 1;
  int32 patterns_applied = 2;
  int32 total_artifacts_generated = 3;
  double predicted_latency_improvement = 4;
  google.protobuf.Timestamp advisory_timestamp = 5;
}
```

**Ξ Formula**:
```
Ξ(advisory) = Relevance × Actionability × Impact_Realization × Adoption_Depth

Where:
  Relevance = patterns_applied / patterns_recommended
  Actionability = artifacts_compiled / artifacts_generated
  Impact_Realization = actual_improvement / predicted_improvement
  Adoption_Depth = patterns_maintained_6mo / patterns_applied
```

## Proto Best Practices

### 1. Never Change Field Numbers

```protobuf
message Foo {
  string name = 1;        // NEVER change to "string name = 2"
  // int32 old_field = 2; // Mark deprecated, don't reuse number
}
```

### 2. Use Reserved for Deleted Fields

```protobuf
message Foo {
  reserved 2, 15, 9 to 11;
  reserved "old_field", "deprecated_field";
  string name = 1;
}
```

### 3. Add, Don't Remove

```protobuf
// ✓ Good: Add new field
message Foo {
  string name = 1;
  string new_field = 3;   // v1.1 addition
}

// ✗ Bad: Remove existing field
message Foo {
  // string name = 1;     // Breaking change!
}
```

### 4. Use Enums with Explicit Values

```protobuf
// ✓ Good
enum Status {
  STATUS_UNSPECIFIED = 0;  // Always have UNSPECIFIED = 0
  STATUS_ACTIVE = 1;
  STATUS_INACTIVE = 2;
}

// ✗ Bad
enum Status {
  ACTIVE = 1;              // Missing UNSPECIFIED
}
```

## Code Generation

### For Go (API Gateway)

```bash
# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/advisory/v1/advisory.proto
```

Generates:
- `advisory.pb.go`: Message types
- `advisory_grpc.pb.go`: Service interface

### For Python (CLI)

```bash
# Install tools
pip install grpcio-tools

# Generate
python -m grpc_tools.protoc -I. \
    --python_out=. --grpc_python_out=. \
    proto/advisory/v1/advisory.proto
```

Generates:
- `advisory_pb2.py`: Message types
- `advisory_pb2_grpc.py`: Service stubs

## Architecture

```
┌─────────────┐
│   ose-cli   │  (Python)
│   Client    │
└──────┬──────┘
       │ gRPC over HTTP/2
       │
┌──────▼──────────────────────┐
│   API Gateway (Go)          │
│   - Request validation      │
│   - Auth/authz              │
│   - Rate limiting           │
│   - Observability           │
└──────┬──────────────────────┘
       │
       ├─────► Pattern Graph Client ──► Neo4j (Pattern Library)
       │
       ├─────► Blueprint Generator ──► Template Engine
       │
       └─────► Validation Service ──► Static Analysis

```

## Roadmap

### Week 2: Protocol Definition ✅
- Complete proto specification
- Documentation
- Design rationale

### Week 3-4: Go Implementation
- Generate Go code
- Implement API Gateway
- Mock service handlers

### Week 5-6: Pattern Graph Integration
- Neo4j client
- Cypher query builder
- Pattern matching logic

### Week 7-8: Blueprint Generator
- Template engine (Jinja2)
- Artifact generation
- File writing service

### Week 9-10: Full Service Implementation
- Validation service
- Registration service
- Metrics collection

## Philosophy in Protocol

The proto definition embodies architectural philosophy:

**1. Constraints, Not Solutions**

Engineers specify constraints (throughput, latency, consistency). The service determines solutions (patterns, architectures). This separation enables the system to improve recommendations without CLI changes.

**2. Confidence Scores**

Every recommendation includes confidence. The system admits uncertainty. Engineers make informed decisions.

**3. Transparent Trade-offs**

The `TradeOff` message makes costs explicit. No hidden complexity. No "magic" solutions.

**4. Learning Loop**

`RegisterService` closes the loop. Every deployment teaches the system. Knowledge compounds.

---

*Part of the Omnifex Synthesis Engine (OSE) - Pillar II: Advisory Service*

*"The standard meter of architectural consultation—typed, versioned, and language-agnostic."*
