package policy

import (
	"testing"
)

func TestDefaultPolicy(t *testing.T) {
	policy := DefaultPolicy()

	if policy.Version != "1" {
		t.Errorf("Expected version 1, got %s", policy.Version)
	}

	if len(policy.Allowlist.Commands) == 0 {
		t.Error("Expected at least one allowlisted command")
	}

	if policy.Limits.MaxIterations <= 0 {
		t.Error("Expected positive max_iterations")
	}
}

func TestPolicyValidation(t *testing.T) {
	policy := DefaultPolicy()

	if err := policy.Validate(); err != nil {
		t.Errorf("Default policy should be valid: %v", err)
	}

	// Test invalid policy
	invalidPolicy := Policy{
		Allowlist: Allowlist{
			Commands: [][]string{},
		},
		Limits: Limits{
			MaxIterations: 0,
		},
	}

	if err := invalidPolicy.Validate(); err == nil {
		t.Error("Invalid policy should fail validation")
	}
}

func TestIsPathAllowed(t *testing.T) {
	policy := DefaultPolicy()

	// Test allowed path
	if !policy.IsPathAllowed("/workspace/src/main.go") {
		t.Error("Workspace path should be allowed")
	}

	// Test denied path
	if policy.IsPathAllowed("/.sentinel/config.yaml") {
		t.Error("Denied path should not be allowed")
	}

	if policy.IsPathAllowed("/.git/config") {
		t.Error("Git path should not be allowed")
	}
}

func TestIsCommandAllowed(t *testing.T) {
	policy := DefaultPolicy()

	// Test allowed command
	if !policy.IsCommandAllowed("go", []string{"build"}) {
		t.Error("go build should be allowed")
	}

	// Test disallowed command
	if policy.IsCommandAllowed("rm", []string{"-rf", "/"}) {
		t.Error("rm -rf should not be allowed")
	}

	// Test partial match (should fail)
	if policy.IsCommandAllowed("go", []string{"build", "--unsafe"}) {
		t.Error("go build --unsafe should not be allowed")
	}
}
