package tools

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// Patcher handles patch application with safety checks
type Patcher struct {
	Workspace  string
	DenyPaths  []string
	DenyGlobs  []string
}

// NewPatcher creates a new patcher instance
func NewPatcher(workspace string, denyPaths, denyGlobs []string) *Patcher {
	return &Patcher{
		Workspace: workspace,
		DenyPaths: denyPaths,
		DenyGlobs: denyGlobs,
	}
}

// ApplyResult represents the result of applying a patch
type ApplyResult struct {
	Applied bool
	Files   []string
	Error   string
}

// ApplyPatch applies a unified diff with safety checks
func (p *Patcher) ApplyPatch(diff string) *ApplyResult {
	// Parse the unified diff
	patches, err := parseUnifiedDiff(diff)
	if err != nil {
		return &ApplyResult{
			Applied: false,
			Error:   fmt.Sprintf("failed to parse diff: %v", err),
		}
	}

	var appliedFiles []string
	var errors []string

	// Apply each patch with safety checks
	for _, patch := range patches {
		// Check if path is allowed
		if !p.isPathAllowed(patch.File) {
			errors = append(errors, fmt.Sprintf("path not allowed: %s", patch.File))
			continue
		}

		// Apply the patch
		if err := p.applySinglePatch(patch); err != nil {
			errors = append(errors, fmt.Sprintf("failed to apply patch to %s: %v", patch.File, err))
			continue
		}

		appliedFiles = append(appliedFiles, patch.File)
	}

	success := len(errors) == 0
	var errorMsg string
	if len(errors) > 0 {
		errorMsg = strings.Join(errors, "; ")
	}

	return &ApplyResult{
		Applied: success,
		Files:   appliedFiles,
		Error:   errorMsg,
	}
}

// isPathAllowed checks if a path is allowed by the security policy
func (p *Patcher) isPathAllowed(path string) bool {
	// Normalize path
	cleanPath := filepath.Clean(path)

	// Check deny paths
	for _, denyPath := range p.DenyPaths {
		if strings.HasPrefix(cleanPath, denyPath) {
			return false
		}
	}

	// Check deny globs
	for _, denyGlob := range p.DenyGlobs {
		if matched, _ := filepath.Match(denyGlob, cleanPath); matched {
			return false
		}
	}

	// Ensure path is within workspace
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return false
	}

	absWorkspace, err := filepath.Abs(p.Workspace)
	if err != nil {
		return false
	}

	return strings.HasPrefix(absPath, absWorkspace)
}

// applySinglePatch applies a single file patch
func (p *Patcher) applySinglePatch(patch FilePatch) error {
	// TODO: Implement actual patch application
	// This would involve:
	// 1. Reading the current file
	// 2. Applying the hunks in order
	// 3. Writing the modified file
	// 4. Validating the result

	return errors.New("patch application not implemented")
}

// FilePatch represents a patch for a single file
type FilePatch struct {
	File  string
	Hunks []Hunk
}

// Hunk represents a hunk of changes
type Hunk struct {
	OldStart int
	OldCount int
	NewStart int
	NewCount int
	Lines    []string
}

// parseUnifiedDiff parses a unified diff string
func parseUnifiedDiff(diff string) ([]FilePatch, error) {
	// TODO: Implement unified diff parsing
	// This is a simplified placeholder
	return []FilePatch{}, errors.New("diff parsing not implemented")
}
