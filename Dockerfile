# Sentinel-AI CLI - Proof of Concept (POC) for personal use only
# Not intended for production environments

# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o sentinel-ai ./cmd/sentinel-ai

# Runtime stage
FROM alpine:latest

# Install required tools
RUN apk add --no-cache \
    git \
    bash \
    curl \
    ca-certificates

# Install Semgrep
RUN curl -sL https://github.com/returntocorp/semgrep/releases/latest/download/semgrep-v1.45.0-linux-x86_64.tar.gz | \
    tar -xz -C /usr/local/bin --strip-components=1 semgrep/semgrep

# Install GitHub CLI
RUN curl -sL https://github.com/cli/cli/releases/latest/download/gh_2.40.1_linux_amd64.tar.gz | \
    tar -xz -C /usr/local/bin --strip-components=1 gh_2.40.1_linux_amd64/bin/gh

COPY --from=builder /app/sentinel-ai /usr/local/bin/sentinel-ai

# Create non-root user
RUN adduser -D -s /bin/bash sentinel
USER sentinel
WORKDIR /workspace

ENTRYPOINT ["sentinel-ai"]
