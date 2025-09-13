package policy

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Policy represents the configuration policy for sentinel-ai
type Policy struct {
	Version string                 `yaml:"version" json:"version"`
	Modes   map[string]Mode        `yaml:"modes" json:"modes"`
	Models  ModelConfig            `yaml:"models" json:"models"`
	Limits  Limits                 `yaml:"limits" json:"limits"`
	Allowlist Allowlist            `yaml:"allowlist" json:"allowlist"`
	Patch    PatchConfig           `yaml:"patch" json:"patch"`
	Security SecurityConfig        `yaml:"security" json:"security"`
	Logging  LoggingConfig         `yaml:"logging" json:"logging"`
}

// Mode defines operational modes
type Mode struct {
	ReadOnly      bool `yaml:"read_only" json:"read_only"`
	Network       bool `yaml:"network" json:"network"`
	MaxRuntimeSec int  `yaml:"max_runtime_sec" json:"max_runtime_sec"`
	MaxTokens     int  `yaml:"max_tokens" json:"max_tokens"`
}

// ModelConfig defines LLM model configuration
type ModelConfig struct {
	PrimaryAlias   string `yaml:"primary_alias" json:"primary_alias"`
	SecondaryAlias string `yaml:"secondary_alias" json:"secondary_alias"`
}

// Limits defines various limits
type Limits struct {
	MaxFiles       int `yaml:"max_files" json:"max_files"`
	MaxFileBytes   int `yaml:"max_file_bytes" json:"max_file_bytes"`
	MaxPatchBytes  int `yaml:"max_patch_bytes" json:"max_patch_bytes"`
	MaxIterations  int `yaml:"max_iterations" json:"max_iterations"`
}

// Allowlist defines allowed commands
type Allowlist struct {
	Commands [][]string `yaml:"commands" json:"commands"`
}

// PatchConfig defines patch-related settings
type PatchConfig struct {
	RequireTests  bool   `yaml:"require_tests" json:"require_tests"`
	CommitStyle   string `yaml:"commit_style" json:"commit_style"`
	Signoff       bool   `yaml:"signoff" json:"signoff"`
}

// SecurityConfig defines security-related settings
type SecurityConfig struct {
	DenyPaths  []string `yaml:"deny_paths" json:"deny_paths"`
	DenyGlobs  []string `yaml:"deny_globs" json:"deny_globs"`
}

// LoggingConfig defines logging settings
type LoggingConfig struct {
	PIIRedaction bool `yaml:"pii_redaction" json:"pii_redaction"`
}

// DefaultPolicy returns a default policy configuration
func DefaultPolicy() Policy {
	return Policy{
		Version: "1",
		Modes: map[string]Mode{
			"default": {
				ReadOnly:      true,
				Network:       false,
				MaxRuntimeSec: 300,
				MaxTokens:     200000,
			},
			"apply": {
				ReadOnly:      false,
				Network:       false,
				MaxRuntimeSec: 300,
				MaxTokens:     200000,
			},
		},
		Models: ModelConfig{
			PrimaryAlias:   "gpt-4",
			SecondaryAlias: "claude-3-sonnet",
		},
		Limits: Limits{
			MaxFiles:       4000,
			MaxFileBytes:   800000,
			MaxPatchBytes:  200000,
			MaxIterations:  4,
		},
		Allowlist: Allowlist{
			Commands: [][]string{
				{"go", "build"},
				{"go", "test", "-cover"},
				{"cargo", "build"},
				{"cargo", "llvm-cov"},
				{"semgrep", "--config", "auto"},
				{"codeql", "database", "analyze"},
				{"gh", "pr", "create"},
			},
		},
		Patch: PatchConfig{
			RequireTests:  true,
			CommitStyle:   "conventional",
			Signoff:       true,
		},
		Security: SecurityConfig{
			DenyPaths: []string{"/.sentinel", "/AGENT.md", "/.git", "/etc", "/usr"},
			DenyGlobs: []string{"**/.sentinel/**", "**/.git/**"},
		},
		Logging: LoggingConfig{
			PIIRedaction: true,
		},
	}
}

// Load loads a policy from a file path
func Load(path string) (Policy, error) {
	if path == "" {
		return DefaultPolicy(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Policy{}, err
	}

	var policy Policy
	if err := yaml.Unmarshal(data, &policy); err != nil {
		return Policy{}, err
	}

	// Validate policy
	if err := policy.Validate(); err != nil {
		return Policy{}, err
	}

	return policy, nil
}

// Validate validates the policy configuration
func (p Policy) Validate() error {
	if len(p.Allowlist.Commands) == 0 {
		return errors.New("no allowlisted commands")
	}

	if p.Limits.MaxIterations <= 0 {
		return errors.New("max_iterations must be positive")
	}

	if p.Limits.MaxFiles <= 0 {
		return errors.New("max_files must be positive")
	}

	return nil
}

// IsPathAllowed checks if a path is allowed by security policy
func (p Policy) IsPathAllowed(path string) bool {
	// Check deny paths
	for _, denyPath := range p.Security.DenyPaths {
		if strings.HasPrefix(path, denyPath) {
			return false
		}
	}

	// Check deny globs
	for _, denyGlob := range p.Security.DenyGlobs {
		if matched, _ := filepath.Match(denyGlob, path); matched {
			return false
		}
	}

	return true
}

// IsCommandAllowed checks if a command is in the allowlist
func (p Policy) IsCommandAllowed(cmd string, args []string) bool {
	for _, allowed := range p.Allowlist.Commands {
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

// JSON returns the policy as JSON
func (p Policy) JSON() []byte {
	data, _ := json.Marshal(p)
	return data
}
