# Sentinel-AI CLI

> [!CAUTION]
> This is a proof of concept (POC) for personal use only. Not intended for production environments.**

A hermetic CLI for security scanning and dead-code detection with LLM-powered triage and remediation.

## Features

- **Security Scanning**: SAST analysis with Semgrep and CodeQL
- **Dead Code Detection**: Static analysis to identify unused code
- **LLM-Powered Triage**: Intelligent prioritization and evidence gathering
- **Safe Patch Application**: Controlled patch application with safety checks
- **Policy-Driven**: Configurable security and operational policies
- **Audit Logging**: Comprehensive logging of all operations

### Limitations
- Not production-ready
- Limited security testing
- Basic error handling
- No warranty or support
- Use at your own risk

### Intended Use
- Personal projects
- Learning and experimentation
- Proof of concept development
- Educational purposes

**Do not use in production environments or with sensitive data.**

## Quick Start

### Installation

```bash
# Build from source
make build

# Or install directly
go install github.com/Siddhant-K-code/sentinel-ai/cmd/sentinel-ai@latest

# Or use Docker
docker build -t sentinel-ai .
```

### Basic Usage

```bash
# Scan for security issues and dead code
sentinel-ai scan --security --dead-code --sarif ./out/scan.sarif --plan ./out/plan.json

# Apply patches locally
sentinel-ai apply --plan ./out/plan.json --approve-level low

# Create a pull request
sentinel-ai pr --title "Security fixes and cleanup" --plan ./out/plan.json --draft

# Run with custom policy
sentinel-ai scan --policy ./examples/policy.yaml --agent ./examples/AGENT.md --security --dead-code
```

### Docker Usage

```bash
# Run in container
docker run -v $(pwd):/workspace sentinel-ai scan --repo /workspace --security --dead-code

# With custom policy
docker run -v $(pwd):/workspace -v $(pwd)/.sentinel:/workspace/.sentinel sentinel-ai scan --repo /workspace --security
```

### Configuration

Create a `.sentinel/policy.yaml` file in your repository:

```yaml
version: 1
modes:
  default:
    read_only: true
    network: false
    max_runtime_sec: 300
    max_tokens: 200000
models:
  primary_alias: gpt-4
  secondary_alias: claude-3-sonnet
limits:
  max_files: 4000
  max_file_bytes: 800000
  max_patch_bytes: 200000
  max_iterations: 4
allowlist:
  commands:
    - ["go", "build"]
    - ["go", "test", "-cover"]
    - ["semgrep", "--config", "auto"]
    - ["codeql", "database", "analyze"]
security:
  deny_paths: ["/.sentinel", "/AGENT.md", "/.git"]
  deny_globs: ["**/.sentinel/**", "**/.git/**"]
```

## Commands

### `scan`

Performs security scanning and dead-code detection:

```bash
sentinel-ai scan [flags]

Flags:
  --repo string      Repository path to scan (default ".")
  --agent string     Path to AGENT.md file (default "./AGENT.md")
  --policy string    Path to policy file (default "./.sentinel/policy.yaml")
  --sarif string     SARIF output file path
  --plan string      Plan output file path
  --log string       Log output file path
  --security         Enable security scanning
  --dead-code        Enable dead-code detection
```

### `apply`

Applies patches from a plan file:

```bash
sentinel-ai apply [flags]

Flags:
  --plan string        Path to plan file (required)
  --approve-level string  Approval level: low, medium, high (default "low")
  --policy string      Path to policy file (default "./.sentinel/policy.yaml")
```

### `pr`

Creates a pull request with proposed changes:

```bash
sentinel-ai pr [flags]

Flags:
  --title string     PR title (required)
  --body string      PR body (file path or text)
  --draft           Create as draft PR
  --plan string     Path to plan file (required)
  --policy string   Path to policy file (default "./.sentinel/policy.yaml")
```

## Exit Codes

- `0`: Success/no actionable findings
- `10`: Actionable security findings present
- `11`: Dead code candidates found
- `20`: Policy violation (attempted forbidden operation)
- `>100`: Internal error

## Environment Variables

```bash
SENTINEL_OPENAI_API_KEY=your_openai_key
SENTINEL_ANTHROPIC_API_KEY=your_anthropic_key
SENTINEL_MODEL_PRIMARY=gpt-4
SENTINEL_MODEL_SECONDARY=claude-3-sonnet
SENTINEL_NO_NETWORK=1        # Disable network access
SENTINEL_POLICY=/path/to/policy.yaml
```

## Examples

See the `examples/` directory for:
- `policy.yaml` - Comprehensive policy configuration
- `AGENT.md` - Project conventions template
- `usage.sh` - Complete usage examples

Run the examples:
```bash
# Run all examples
./examples/usage.sh

# Test with example Go code
sentinel-ai scan --repo ./example --security --dead-code
```

## Development

### Quick Setup

```bash
# Run the setup script (installs Go, golangci-lint, and builds the project)
./setup.sh

# Or setup manually
make setup
```

### Manual Setup

```bash
# Install Go (if not already installed)
curl -L https://go.dev/dl/go1.21.5.linux-amd64.tar.gz | sudo tar -xzC /usr/local
export PATH=$PATH:/usr/local/go/bin

# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
export PATH=$PATH:$(go env GOPATH)/bin

# Setup project
make setup
```

### Development Commands

```bash
# Build
make build

# Test
make test

# Lint
make lint

# Clean
make clean

# Run everything
make all

# Run examples
./examples/usage.sh
```

## Security

- **Read-only by default**: Scans don't modify code unless explicitly requested
- **Command allowlisting**: Only pre-approved commands can be executed
- **Path restrictions**: Cannot modify policy files or system directories
- **Network isolation**: No network access unless explicitly enabled
- **Audit logging**: All operations are logged for security review
