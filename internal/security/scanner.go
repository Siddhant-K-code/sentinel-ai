package security

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Siddhant-K-code/sentinel-ai/internal/tools"
)

// Scanner performs security scanning using various tools
type Scanner struct {
	runner *tools.Runner
	workspace string
}

// ScanResult represents the result of a security scan
type ScanResult struct {
	Tool     string    `json:"tool"`
	Findings []Finding `json:"findings"`
	Duration time.Duration `json:"duration"`
	Error    string    `json:"error,omitempty"`
}

// Finding represents a security finding
type Finding struct {
	RuleID      string `json:"rule_id"`
	Message     string `json:"message"`
	Severity    string `json:"severity"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Description string `json:"description"`
	Confidence  string `json:"confidence"`
}

// NewScanner creates a new security scanner
func NewScanner(runner *tools.Runner, workspace string) *Scanner {
	return &Scanner{
		runner:    runner,
		workspace: workspace,
	}
}

// Scan performs security scanning using available tools
func (s *Scanner) Scan(ctx context.Context) ([]ScanResult, error) {
	var results []ScanResult

	// Try Semgrep first
	if semgrepResult := s.runSemgrep(ctx); semgrepResult != nil {
		results = append(results, *semgrepResult)
	}

	// Try CodeQL if available
	if codeqlResult := s.runCodeQL(ctx); codeqlResult != nil {
		results = append(results, *codeqlResult)
	}

	return results, nil
}

// runSemgrep runs Semgrep security scanning
func (s *Scanner) runSemgrep(ctx context.Context) *ScanResult {
	start := time.Now()

	// Check if semgrep is available
	if _, err := exec.LookPath("semgrep"); err != nil {
		return &ScanResult{
			Tool:     "semgrep",
			Findings: []Finding{},
			Duration: time.Since(start),
			Error:    "semgrep not found in PATH",
		}
	}

	// Run semgrep with auto config
	result := s.runner.Run(ctx, "semgrep", "--config", "auto", "--json", s.workspace)
	if result.Error != nil {
		return &ScanResult{
			Tool:     "semgrep",
			Findings: []Finding{},
			Duration: time.Since(start),
			Error:    result.Error.Error(),
		}
	}

	// Parse semgrep JSON output
	var semgrepOutput struct {
		Results []struct {
			CheckID  string `json:"check_id"`
			Path     string `json:"path"`
			Start    struct {
				Line   int `json:"line"`
				Column int `json:"column"`
			} `json:"start"`
			End struct {
				Line   int `json:"line"`
				Column int `json:"column"`
			} `json:"end"`
			Extra struct {
				Message     string `json:"message"`
				Severity    string `json:"severity"`
				Description string `json:"description"`
				Confidence  string `json:"confidence"`
			} `json:"extra"`
		} `json:"results"`
	}

	if err := json.Unmarshal(result.Stdout, &semgrepOutput); err != nil {
		return &ScanResult{
			Tool:     "semgrep",
			Findings: []Finding{},
			Duration: time.Since(start),
			Error:    fmt.Sprintf("failed to parse semgrep output: %v", err),
		}
	}

	// Convert to our Finding format
	var findings []Finding
	for _, r := range semgrepOutput.Results {
		// Make path relative to workspace
		relPath, _ := filepath.Rel(s.workspace, r.Path)
		if relPath == "." {
			relPath = filepath.Base(r.Path)
		}

		findings = append(findings, Finding{
			RuleID:      r.CheckID,
			Message:     r.Extra.Message,
			Severity:    strings.ToLower(r.Extra.Severity),
			File:        relPath,
			Line:        r.Start.Line,
			Column:      r.Start.Column,
			Description: r.Extra.Description,
			Confidence:  strings.ToLower(r.Extra.Confidence),
		})
	}

	return &ScanResult{
		Tool:     "semgrep",
		Findings: findings,
		Duration: time.Since(start),
	}
}

// runCodeQL runs CodeQL security scanning
func (s *Scanner) runCodeQL(ctx context.Context) *ScanResult {
	start := time.Now()

	// Check if codeql is available
	if _, err := exec.LookPath("codeql"); err != nil {
		return &ScanResult{
			Tool:     "codeql",
			Findings: []Finding{},
			Duration: time.Since(start),
			Error:    "codeql not found in PATH",
		}
	}

	// For now, return empty result as CodeQL requires database creation
	// In a full implementation, this would:
	// 1. Create a CodeQL database
	// 2. Run analysis queries
	// 3. Parse the SARIF output
	return &ScanResult{
		Tool:     "codeql",
		Findings: []Finding{},
		Duration: time.Since(start),
		Error:    "CodeQL integration not fully implemented",
	}
}

// GenerateSARIF converts scan results to SARIF format
func (s *Scanner) GenerateSARIF(results []ScanResult) ([]byte, error) {
	sarif := map[string]interface{}{
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"version": "2.1.0",
		"runs":    []map[string]interface{}{},
	}

	var runs []map[string]interface{}
	for _, result := range results {
		if result.Error != "" {
			continue // Skip failed scans
		}

		run := map[string]interface{}{
			"tool": map[string]interface{}{
				"driver": map[string]interface{}{
					"name":    result.Tool,
					"version": "1.0.0",
				},
			},
			"results": []map[string]interface{}{},
		}

		var sarifResults []map[string]interface{}
		for _, finding := range result.Findings {
			sarifResult := map[string]interface{}{
				"ruleId": finding.RuleID,
				"message": map[string]interface{}{
					"text": finding.Message,
				},
				"level": s.severityToLevel(finding.Severity),
				"locations": []map[string]interface{}{
					{
						"physicalLocation": map[string]interface{}{
							"artifactLocation": map[string]interface{}{
								"uri": finding.File,
							},
							"region": map[string]interface{}{
								"startLine":   finding.Line,
								"startColumn": finding.Column,
							},
						},
					},
				},
			}

			if finding.Description != "" {
				sarifResult["message"].(map[string]interface{})["text"] = finding.Description
			}

			sarifResults = append(sarifResults, sarifResult)
		}

		run["results"] = sarifResults
		runs = append(runs, run)
	}

	sarif["runs"] = runs

	return json.MarshalIndent(sarif, "", "  ")
}

// severityToLevel converts severity string to SARIF level
func (s *Scanner) severityToLevel(severity string) string {
	switch strings.ToLower(severity) {
	case "error", "critical", "high":
		return "error"
	case "warning", "medium":
		return "warning"
	case "info", "low":
		return "note"
	default:
		return "note"
	}
}
