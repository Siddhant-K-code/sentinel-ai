package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/Siddhant-K-code/sentinel-ai/internal/policy"
	"github.com/Siddhant-K-code/sentinel-ai/internal/security"
	"github.com/Siddhant-K-code/sentinel-ai/internal/deadcode"
	"github.com/Siddhant-K-code/sentinel-ai/internal/tools"
	"github.com/Siddhant-K-code/sentinel-ai/internal/logging"
)

// Options defines engine configuration options
type Options struct {
	Repo      string
	AgentPath string
	Policy    policy.Policy
	LogPath   string
}

// ScanOpts defines scan operation options
type ScanOpts struct {
	Security bool
	DeadCode bool
}

// ScanResult represents the result of a scan operation
type ScanResult struct {
	SARIF     []byte
	Plan      Plan
	ExitCode  int
	Summary   string
}

// Plan represents a planned set of changes
type Plan struct {
	Steps    []Step `json:"steps"`
	Metadata Metadata `json:"metadata"`
}

// Step represents a single step in the plan
type Step struct {
	Name         string   `json:"name"`
	Why          string   `json:"why"`
	BudgetTokens int      `json:"budget_tokens"`
	Tools        []string `json:"tools"`
	StopAfter    string   `json:"stop_after,omitempty"`
}

// Metadata contains plan metadata
type Metadata struct {
	CreatedAt    time.Time `json:"created_at"`
	SuccessCriteria []string `json:"success_criteria"`
	TotalTokens  int       `json:"total_tokens"`
}

// ApplyResult represents the result of applying patches
type ApplyResult struct {
	Success bool
	Error   string
	Files   []string
}

// PROptions defines PR creation options
type PROptions struct {
	Title    string
	Body     string
	Draft    bool
	PlanPath string
}

// PRResult represents the result of creating a PR
type PRResult struct {
	URL string
}

// Engine represents the main sentinel-ai engine
type Engine struct {
	options     Options
	policy      policy.Policy
	runner      *tools.Runner
	scanner     *security.Scanner
	detector    *deadcode.Detector
	auditLogger *logging.AuditLogger
}

// New creates a new engine instance
func New(ctx context.Context, opts Options) (*Engine, error) {
	// Create tool runner
	runner := tools.NewRunner(opts.Policy.Allowlist.Commands, time.Duration(opts.Policy.Modes["default"].MaxRuntimeSec)*time.Second)

	// Create audit logger
	auditLogger, err := logging.NewAuditLogger(opts.LogPath, opts.Policy.Logging.PIIRedaction)
	if err != nil {
		return nil, err
	}

	// Create security scanner
	scanner := security.NewScanner(runner, opts.Repo)

	// Create dead code detector
	detector := deadcode.NewDetector(runner, opts.Repo)

	return &Engine{
		options:     opts,
		policy:      opts.Policy,
		runner:      runner,
		scanner:     scanner,
		detector:    detector,
		auditLogger: auditLogger,
	}, nil
}

// Scan performs security and dead-code scanning
func (e *Engine) Scan(ctx context.Context, opts ScanOpts) (*ScanResult, error) {
	start := time.Now()
	var sarifData []byte
	var plan Plan
	exitCode := 0
	summary := "Scan completed successfully"

	// Run security scanning if requested
	if opts.Security {
		e.auditLogger.LogToolCall("security", "scanner", []string{"security", "scan"}, 0, "started", nil)

		securityResults, err := e.scanner.Scan(ctx)
		if err != nil {
			e.auditLogger.LogToolCall("security", "scanner", []string{"security", "scan"}, time.Since(start), "error", err)
			return nil, err
		}

		// Generate SARIF output
		sarifData, err = e.scanner.GenerateSARIF(securityResults)
		if err != nil {
			e.auditLogger.LogToolCall("security", "scanner", []string{"sarif", "generate"}, time.Since(start), "error", err)
			return nil, err
		}

		// Count findings
		totalFindings := 0
		for _, result := range securityResults {
			totalFindings += len(result.Findings)
		}

		e.auditLogger.LogScanResult("security", totalFindings, time.Since(start))

		if totalFindings > 0 {
			exitCode = 10 // Security findings present
			summary = fmt.Sprintf("Found %d security findings", totalFindings)
		}
	}

	// Run dead code detection if requested
	if opts.DeadCode {
		e.auditLogger.LogToolCall("deadcode", "detector", []string{"deadcode", "detect"}, 0, "started", nil)

		deadCodeResult, err := e.detector.Detect(ctx)
		if err != nil {
			e.auditLogger.LogToolCall("deadcode", "detector", []string{"deadcode", "detect"}, time.Since(start), "error", err)
			return nil, err
		}

		e.auditLogger.LogScanResult("deadcode", len(deadCodeResult.Symbols), time.Since(start))

		if len(deadCodeResult.Symbols) > 0 {
			if exitCode == 0 {
				exitCode = 11 // Dead code found
			}
			summary = fmt.Sprintf("%s; Found %d dead code symbols", summary, len(deadCodeResult.Symbols))
		}
	}

	// Create a basic plan
	plan = Plan{
		Steps: []Step{
			{
				Name:         "analysis",
				Why:          "Analyze codebase for security issues and dead code",
				BudgetTokens: 1000,
				Tools:        []string{"scanner", "detector"},
			},
		},
		Metadata: Metadata{
			CreatedAt:       time.Now(),
			SuccessCriteria: []string{"build passes", "no high/critical findings"},
			TotalTokens:     1000,
		},
	}

	// If no SARIF data was generated, create empty SARIF
	if sarifData == nil {
		sarifData = []byte(`{"version":"2.1.0","runs":[]}`)
	}

	return &ScanResult{
		SARIF:    sarifData,
		Plan:     plan,
		ExitCode: exitCode,
		Summary:  summary,
	}, nil
}

// Apply applies patches from a plan
func (e *Engine) Apply(ctx context.Context, plan Plan, approveLevel string) (*ApplyResult, error) {
	// TODO: Implement patch application logic
	// This would include:
	// 1. Validate plan
	// 2. Check approval level
	// 3. Apply patches with safety checks
	// 4. Run tests
	// 5. Commit changes

	return &ApplyResult{
		Success: true,
		Files:   []string{},
	}, nil
}

// CreatePR creates a pull request
func (e *Engine) CreatePR(ctx context.Context, opts PROptions) (*PRResult, error) {
	// TODO: Implement PR creation logic
	// This would include:
	// 1. Read plan file
	// 2. Apply patches to working directory
	// 3. Commit changes
	// 4. Push branch
	// 5. Create PR via GitHub CLI

	return &PRResult{
		URL: "https://github.com/example/repo/pull/123",
	}, nil
}
