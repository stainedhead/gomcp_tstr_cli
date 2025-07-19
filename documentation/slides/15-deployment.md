# Deployment

## Binary Distribution

- **Single Binary**: Compiled as a statically-linked binary
- **Cross-Platform**: Builds for Linux, macOS, Windows
- **Minimal Dependencies**: No external runtime dependencies

## Configuration Management

- **Example Configurations**: Provided in examples/ directory
- **Environment Variables**: Support for environment-based config
- **Validation**: Configuration validation on startup

## Installation

```bash
# Build from source
git clone <repository-url>
cd mcp_tstr
go build -o mcp_tstr .

# Configure
cp examples/mcp_tstr.config mcp_tstr.config
cp examples/mcp.json mcp.json

# Run
./mcp_tstr --help
```

---
