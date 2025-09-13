package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// AuditLogger handles structured audit logging
type AuditLogger struct {
	file   *os.File
	redact bool
}

// LogEntry represents a single audit log entry
type LogEntry struct {
	Timestamp time.Time              `json:"ts"`
	Step      string                 `json:"step"`
	Event     string                 `json:"event"`
	Tool      string                 `json:"tool,omitempty"`
	Args      []string               `json:"args,omitempty"`
	Duration  int64                  `json:"duration_ms,omitempty"`
	Status    string                 `json:"status"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(logPath string, redactPII bool) (*AuditLogger, error) {
	var file *os.File
	var err error

	if logPath != "" {
		file, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		file = os.Stdout
	}

	return &AuditLogger{
		file:   file,
		redact: redactPII,
	}, nil
}

// LogToolCall logs a tool call
func (a *AuditLogger) LogToolCall(step, tool string, args []string, duration time.Duration, status string, err error) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Step:      step,
		Event:     "tool.run",
		Tool:      tool,
		Args:      a.maybeRedactArgs(args),
		Duration:  duration.Milliseconds(),
		Status:    status,
	}

	if err != nil {
		entry.Error = a.maybeRedactError(err.Error())
	}

	a.writeEntry(entry)
}

// LogLLMCall logs an LLM call
func (a *AuditLogger) LogLLMCall(step, model string, tokens int, duration time.Duration, status string, err error) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Step:      step,
		Event:     "llm.call",
		Tool:      model,
		Duration:  duration.Milliseconds(),
		Status:    status,
		Metadata: map[string]interface{}{
			"tokens": tokens,
		},
	}

	if err != nil {
		entry.Error = a.maybeRedactError(err.Error())
	}

	a.writeEntry(entry)
}

// LogScanResult logs a scan result
func (a *AuditLogger) LogScanResult(step string, findings int, duration time.Duration) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Step:      step,
		Event:     "scan.result",
		Duration:  duration.Milliseconds(),
		Status:    "ok",
		Metadata: map[string]interface{}{
			"findings": findings,
		},
	}

	a.writeEntry(entry)
}

// LogPatchApplication logs patch application
func (a *AuditLogger) LogPatchApplication(step string, files []string, success bool, err error) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Step:      step,
		Event:     "patch.apply",
		Status:    "ok",
		Metadata: map[string]interface{}{
			"files":   files,
			"success": success,
		},
	}

	if err != nil {
		entry.Error = a.maybeRedactError(err.Error())
		entry.Status = "error"
	}

	a.writeEntry(entry)
}

// LogPolicyViolation logs a policy violation
func (a *AuditLogger) LogPolicyViolation(step, violation string, details map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Step:      step,
		Event:     "policy.violation",
		Status:    "error",
		Error:     violation,
		Metadata:  details,
	}

	a.writeEntry(entry)
}

// writeEntry writes a log entry to the file
func (a *AuditLogger) writeEntry(entry LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry: %v\n", err)
		return
	}

	fmt.Fprintln(a.file, string(data))
}

// maybeRedactArgs redacts sensitive information from command arguments
func (a *AuditLogger) maybeRedactArgs(args []string) []string {
	if !a.redact {
		return args
	}

	redacted := make([]string, len(args))
	for i, arg := range args {
		// Redact potential API keys, tokens, etc.
		if len(arg) > 20 && (containsSensitivePattern(arg) || looksLikeToken(arg)) {
			redacted[i] = "[REDACTED]"
		} else {
			redacted[i] = arg
		}
	}
	return redacted
}

// maybeRedactError redacts sensitive information from error messages
func (a *AuditLogger) maybeRedactError(err string) string {
	if !a.redact {
		return err
	}

	// Simple redaction - in production, use more sophisticated patterns
	if len(err) > 100 {
		return "[REDACTED ERROR]"
	}
	return err
}

// containsSensitivePattern checks if a string contains sensitive patterns
func containsSensitivePattern(s string) bool {
	sensitive := []string{"key", "token", "secret", "password", "auth"}
	lower := strings.ToLower(s)
	for _, pattern := range sensitive {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

// looksLikeToken checks if a string looks like a token or key
func looksLikeToken(s string) bool {
	// Simple heuristic: long alphanumeric strings
	if len(s) < 20 {
		return false
	}

	hasAlpha := false
	hasNum := false
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			hasAlpha = true
		}
		if r >= '0' && r <= '9' {
			hasNum = true
		}
	}

	return hasAlpha && hasNum
}

// Close closes the audit logger
func (a *AuditLogger) Close() error {
	if a.file != nil && a.file != os.Stdout {
		return a.file.Close()
	}
	return nil
}
