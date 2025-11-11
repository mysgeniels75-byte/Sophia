# Protocol Buffer Code Generation

## Overview

The OSE Advisory API uses gRPC with Protocol Buffers for type-safe, language-agnostic communication between the Python CLI and Go API Gateway.

## Prerequisites

To generate Go code from `.proto` definitions, you need:

1. **protoc** (Protocol Buffer Compiler)
   - Version 3.20+ recommended
   - Download: https://github.com/protocolbuffers/protobuf/releases

2. **Go protobuf plugins** (auto-installed via Makefile)
   - `protoc-gen-go` - generates message types
   - `protoc-gen-go-grpc` - generates gRPC service code

## Installation

### macOS
```bash
brew install protobuf
```

### Ubuntu/Debian
```bash
apt-get install -y protobuf-compiler
```

### Manual Installation (All Platforms)
```bash
# Download precompiled binary
VERSION=25.1
curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${VERSION}/protoc-${VERSION}-linux-x86_64.zip
unzip protoc-${VERSION}-linux-x86_64.zip -d /usr/local
```

## Generating Go Code

Once `protoc` is installed:

```bash
# From ose-api/ directory
make proto
```

This command:
1. Installs Go protobuf plugins if not present
2. Generates Go code from `proto/advisory/v1/advisory.proto`
3. Outputs to `api/proto/advisory/v1/*.pb.go`

## Generated Files

After running `make proto`, you should see:

```
api/proto/advisory/v1/
├── advisory.pb.go         # Message type definitions
└── advisory_grpc.pb.go    # gRPC service interfaces
```

These files are **generated code** and should not be manually edited.

## Git Handling

The generated `.pb.go` files are NOT committed to git (see `.gitignore`).

Each developer/CI environment must:
1. Install `protoc`
2. Run `make proto` after `git clone`
3. Rebuild when `.proto` files change

## Current Status

**Week 3 Implementation Note:**

The handler code (`internal/handlers/advisory_handler.go`) is structurally complete but requires protobuf generation to compile.

**To complete the build:**
1. Install `protoc` (see above)
2. Run `make proto`
3. Run `make build`

The server will then compile and start successfully.

## Troubleshooting

### Error: "protoc: command not found"
- Install protoc using one of the methods above
- Verify: `protoc --version`

### Error: "protoc-gen-go: program not found"
- The Makefile should auto-install this
- Manual install: `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0`

### Error: "imports not found"
- Ensure `google/protobuf/timestamp.proto` is available
- This is included with protoc installation

### Generated code has wrong import paths
- Check `option go_package` in the `.proto` file
- Should match: `github.com/mysgeniels75-byte/ose-api/api/proto/advisory/v1;advisoryv1`

## Alternative: Using Buf

[Buf](https://buf.build) is a modern protobuf toolchain that simplifies generation:

```bash
# Install buf
go install github.com/bufbuild/buf/cmd/buf@latest

# Generate code (if buf.yaml is configured)
buf generate
```

We may add Buf support in Week 4 for improved developer experience.

## Python Client Generation (Week 4)

The same `.proto` file will be used to generate Python gRPC stubs for the CLI:

```bash
python -m grpc_tools.protoc \
  --proto_path=proto \
  --python_out=ose-cli/gen \
  --grpc_python_out=ose-cli/gen \
  proto/advisory/v1/advisory.proto
```

This enables type-safe bidirectional communication: Python CLI ↔ Go API Gateway.

## References

- [Protocol Buffers Documentation](https://protobuf.dev/)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [protoc Installation Guide](https://grpc.io/docs/protoc-installation/)
