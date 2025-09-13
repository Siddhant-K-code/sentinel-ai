#!/bin/bash

# Sentinel-AI CLI Usage Examples
# This script demonstrates various usage patterns for the sentinel-ai CLI

set -e

echo "=== Sentinel-AI CLI Usage Examples ==="

# Build the CLI
echo "Building CLI..."
go build -o bin/sentinel-ai ./cmd/sentinel-ai

# Create output directory
mkdir -p out

echo ""
echo "=== 1. Basic Security Scan ==="
./bin/sentinel-ai scan \
  --repo . \
  --security \
  --sarif ./out/security.sarif \
  --plan ./out/security-plan.json \
  --log ./out/security.log

echo "Security scan completed. Results:"
echo "- SARIF: $(wc -l < out/security.sarif) lines"
echo "- Plan: $(wc -l < out/security-plan.json) lines"
echo "- Log: $(wc -l < out/security.log) lines"

echo ""
echo "=== 2. Dead Code Detection ==="
./bin/sentinel-ai scan \
  --repo . \
  --dead-code \
  --sarif ./out/deadcode.sarif \
  --plan ./out/deadcode-plan.json \
  --log ./out/deadcode.log

echo "Dead code detection completed. Results:"
echo "- SARIF: $(wc -l < out/deadcode.sarif) lines"
echo "- Plan: $(wc -l < out/deadcode-plan.json) lines"
echo "- Log: $(wc -l < out/deadcode.log) lines"

echo ""
echo "=== 3. Combined Analysis ==="
./bin/sentinel-ai scan \
  --repo . \
  --security \
  --dead-code \
  --sarif ./out/combined.sarif \
  --plan ./out/combined-plan.json \
  --log ./out/combined.log

echo "Combined analysis completed. Results:"
echo "- SARIF: $(wc -l < out/combined.sarif) lines"
echo "- Plan: $(wc -l < out/combined-plan.json) lines"
echo "- Log: $(wc -l < out/combined.log) lines"

echo ""
echo "=== 4. Apply Patches (Dry Run) ==="
if [ -f "out/combined-plan.json" ]; then
    echo "Applying patches from plan..."
    ./bin/sentinel-ai apply \
      --plan ./out/combined-plan.json \
      --approve-level low
else
    echo "No plan file found, skipping patch application"
fi

echo ""
echo "=== 5. Create Pull Request (Dry Run) ==="
if [ -f "out/combined-plan.json" ]; then
    echo "Creating pull request..."
    ./bin/sentinel-ai pr \
      --title "Security fixes and dead code removal" \
      --body "Automated fixes from sentinel-ai analysis" \
      --plan ./out/combined-plan.json \
      --draft
else
    echo "No plan file found, skipping PR creation"
fi

echo ""
echo "=== 6. Custom Policy Example ==="
./bin/sentinel-ai scan \
  --repo . \
  --policy ./examples/policy.yaml \
  --agent ./examples/AGENT.md \
  --security \
  --dead-code \
  --sarif ./out/custom.sarif \
  --plan ./out/custom-plan.json \
  --log ./out/custom.log

echo "Custom policy scan completed. Results:"
echo "- SARIF: $(wc -l < out/custom.sarif) lines"
echo "- Plan: $(wc -l < out/custom-plan.json) lines"
echo "- Log: $(wc -l < out/custom.log) lines"

echo ""
echo "=== 7. Audit Log Analysis ==="
echo "Recent audit log entries:"
tail -5 out/combined.log | jq -r '.ts + " " + .step + " " + .event + " " + .status'

echo ""
echo "=== Usage Examples Complete ==="
echo "Check the 'out/' directory for all generated files:"
ls -la out/
