.PHONY: build test clean install lint

# Build the CLI
build:
	/usr/local/go/bin/go build -o bin/sentinel-ai ./cmd/sentinel-ai

# Run tests
test:
	/usr/local/go/bin/go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	/usr/local/go/bin/go mod tidy
	/usr/local/go/bin/go mod download

# Lint code
lint:
	/home/vscode/go/bin/golangci-lint run

# Run the CLI with example commands
demo-scan:
	./bin/sentinel-ai scan --security --dead-code --sarif ./out/scan.sarif --plan ./out/plan.json

demo-apply:
	./bin/sentinel-ai apply --plan ./out/plan.json --approve-level low

demo-pr:
	./bin/sentinel-ai pr --title "Security fixes and dead code removal" --plan ./out/plan.json --draft

# Create output directory
setup:
	mkdir -p out bin

# Full build and test
all: deps lint test build
