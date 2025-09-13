package tools

import (
	"context"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	allowlist := [][]string{
		{"go", "build"},
		{"go", "test", "-cover"},
		{"echo", "hello"},
	}

	runner := NewRunner(allowlist, 5*time.Second)

	// Test allowed command
	ctx := context.Background()
	result := runner.Run(ctx, "echo", "hello")
	if result.Error != nil {
		t.Errorf("Expected echo hello to succeed: %v", result.Error)
	}

	// Test disallowed command
	result = runner.Run(ctx, "rm", "-rf", "/")
	if result.Error == nil {
		t.Error("Expected rm -rf to be disallowed")
	}

	// Test partial match (should fail)
	result = runner.Run(ctx, "go", "build", "--unsafe")
	if result.Error == nil {
		t.Error("Expected go build --unsafe to be disallowed")
	}
}

func TestAllowed(t *testing.T) {
	allowlist := [][]string{
		{"go", "build"},
		{"go", "test", "-cover"},
	}

	runner := NewRunner(allowlist, 5*time.Second)

	// Test exact match
	if !runner.allowed("go", []string{"build"}) {
		t.Error("go build should be allowed")
	}

	// Test wrong command
	if runner.allowed("python", []string{"build"}) {
		t.Error("python build should not be allowed")
	}

	// Test wrong args
	if runner.allowed("go", []string{"run"}) {
		t.Error("go run should not be allowed")
	}

	// Test extra args
	if runner.allowed("go", []string{"build", "--unsafe"}) {
		t.Error("go build --unsafe should not be allowed")
	}
}
