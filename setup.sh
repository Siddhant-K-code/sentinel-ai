#!/bin/bash

# Sentinel-AI Setup Script
# This script sets up the development environment for sentinel-ai

set -e

echo "=== Sentinel-AI Setup ==="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    curl -L https://go.dev/dl/go1.21.5.linux-amd64.tar.gz | sudo tar -xzC /usr/local
    export PATH=$PATH:/usr/local/go/bin
else
    echo "Go is already installed: $(go version)"
fi

# Add Go to PATH permanently
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Install golangci-lint if not present
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
else
    echo "golangci-lint is already installed: $(golangci-lint --version)"
fi

# Set up the project
echo "Setting up project..."
export PATH=$PATH:/usr/local/go/bin:$(go env GOPATH)/bin

# Download dependencies
echo "Downloading dependencies..."
go mod tidy
go mod download

# Create directories
echo "Creating directories..."
mkdir -p out bin

# Run tests
echo "Running tests..."
go test ./...

# Build the project
echo "Building project..."
go build -o bin/sentinel-ai ./cmd/sentinel-ai

# Run linting
echo "Running linter..."
golangci-lint run

echo ""
echo "=== Setup Complete ==="
echo "You can now run:"
echo "  make build    - Build the project"
echo "  make test     - Run tests"
echo "  make lint     - Run linter"
echo "  make all      - Run everything"
echo "  ./bin/sentinel-ai --help - Test the CLI"
echo ""
echo "Note: You may need to restart your terminal or run 'source ~/.bashrc' to use the updated PATH."
