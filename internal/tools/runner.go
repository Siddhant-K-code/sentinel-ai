package tools

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

// Runner executes allowlisted commands
type Runner struct {
	Allow   [][]string
	Timeout time.Duration
}

// RunResult represents the result of running a command
type RunResult struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Duration time.Duration
	Error    error
}

// NewRunner creates a new command runner
func NewRunner(allowlist [][]string, timeout time.Duration) *Runner {
	return &Runner{
		Allow:   allowlist,
		Timeout: timeout,
	}
}

// Run executes a command if it's in the allowlist
func (r *Runner) Run(ctx context.Context, cmd string, args ...string) *RunResult {
	start := time.Now()

	// Check if command is allowed
	if !r.allowed(cmd, args) {
		return &RunResult{
			Error:    errors.New("command not allowlisted"),
			Duration: time.Since(start),
		}
	}

	// Create context with timeout
	runCtx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	// Execute command
	execCmd := exec.CommandContext(runCtx, cmd, args...)
	output, err := execCmd.CombinedOutput()

	result := &RunResult{
		Stdout:   output,
		Stderr:   nil, // CombinedOutput puts everything in stdout
		Duration: time.Since(start),
		Error:    err,
	}

	if execCmd.ProcessState != nil {
		result.ExitCode = execCmd.ProcessState.ExitCode()
	}

	return result
}

// allowed checks if a command and its arguments are in the allowlist
func (r *Runner) allowed(cmd string, args []string) bool {
	for _, allowed := range r.Allow {
		if len(allowed) == len(args)+1 && allowed[0] == cmd {
			matches := true
			for i := 1; i < len(allowed); i++ {
				if allowed[i] != args[i-1] {
					matches = false
					break
				}
			}
			if matches {
				return true
			}
		}
	}
	return false
}

// Build runs build commands
func (r *Runner) Build(ctx context.Context, target string) *RunResult {
	if target == "" {
		// Try to detect build system
		if _, err := exec.LookPath("go"); err == nil {
			return r.Run(ctx, "go", "build", "./...")
		}
		if _, err := exec.LookPath("cargo"); err == nil {
			return r.Run(ctx, "cargo", "build")
		}
		return &RunResult{
			Error: errors.New("no supported build system found"),
		}
	}
	return r.Run(ctx, "go", "build", target)
}

// Test runs test commands
func (r *Runner) Test(ctx context.Context, pattern string) *RunResult {
	args := []string{"test"}
	if pattern != "" {
		args = append(args, "-run", pattern)
	}
	args = append(args, "-cover", "./...")
	return r.Run(ctx, "go", args...)
}

// Coverage runs coverage analysis
func (r *Runner) Coverage(ctx context.Context) *RunResult {
	return r.Run(ctx, "go", "test", "-coverprofile=coverage.out", "./...")
}
