package deadcode

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Siddhant-K-code/sentinel-ai/internal/tools"
)

// Detector finds dead code in Go projects
type Detector struct {
	runner    *tools.Runner
	workspace string
}

// DeadCodeResult represents the result of dead code detection
type DeadCodeResult struct {
	Symbols   []Symbol `json:"symbols"`
	Duration  time.Duration `json:"duration"`
	Error     string   `json:"error,omitempty"`
}

// Symbol represents a potentially dead code symbol
type Symbol struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"` // func, var, const, type
	Package     string `json:"package"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Exported    bool   `json:"exported"`
	References  int    `json:"references"`
	LastTouch   string `json:"last_touch,omitempty"`
	Risk        string `json:"risk"` // low, medium, high
	Description string `json:"description"`
}

// NewDetector creates a new dead code detector
func NewDetector(runner *tools.Runner, workspace string) *Detector {
	return &Detector{
		runner:    runner,
		workspace: workspace,
	}
}

// Detect finds dead code in the workspace
func (d *Detector) Detect(ctx context.Context) (*DeadCodeResult, error) {
	start := time.Now()

	// Get coverage information (optional)
	coverage, err := d.getCoverage(ctx)
	if err != nil {
		coverage = make(map[string]float64) // Use empty coverage map
	}

	// Find all Go files
	goFiles, err := d.findGoFiles()
	if err != nil {
		return &DeadCodeResult{
			Duration: time.Since(start),
			Error:    fmt.Sprintf("failed to find Go files: %v", err),
		}, nil
	}

	// Analyze each file for dead code
	var symbols []Symbol
	for _, file := range goFiles {
		fileSymbols, err := d.analyzeFile(file, coverage)
		if err != nil {
			continue // Skip files with errors
		}
		symbols = append(symbols, fileSymbols...)
	}

	// Filter out symbols that are actually used
	deadSymbols := d.filterDeadSymbols(symbols)

	return &DeadCodeResult{
		Symbols:  deadSymbols,
		Duration: time.Since(start),
	}, nil
}

// getCoverage gets test coverage information
func (d *Detector) getCoverage(ctx context.Context) (map[string]float64, error) {
	// Run coverage analysis
	result := d.runner.Coverage(ctx)
	if result.Error != nil {
		return nil, result.Error
	}

	// Parse coverage output (simplified)
	// In a real implementation, this would parse the coverage.out file
	coverage := make(map[string]float64)

	// For now, return empty coverage map
	// Real implementation would parse go test -coverprofile output
	return coverage, nil
}

// findGoFiles finds all Go files in the workspace
func (d *Detector) findGoFiles() ([]string, error) {
	var files []string

	// Simple file walker for Go files
	err := filepath.Walk(d.workspace, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// Skip test files and vendor directories
			if !strings.Contains(path, "_test.go") &&
			   !strings.Contains(path, "/vendor/") &&
			   !strings.Contains(path, "/.git/") {
				files = append(files, path)
			}
		}

		return nil
	})

	return files, err
}

// analyzeFile analyzes a Go file for potentially dead symbols
func (d *Detector) analyzeFile(filePath string, coverage map[string]float64) ([]Symbol, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var symbols []Symbol
	packageName := node.Name.Name

	// Collect all function calls in the file to check for references
	callMap := make(map[string]bool)
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				callMap[ident.Name] = true
			}
		}
		return true
	})

	// Analyze functions
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name != nil && !x.Name.IsExported() && x.Name.Name != "main" {
				// Check if this function is called anywhere in the file
				isCalled := callMap[x.Name.Name]
				references := 0
				if isCalled {
					references = 1
				}

				symbols = append(symbols, Symbol{
					Name:       x.Name.Name,
					Kind:       "func",
					Package:    packageName,
					File:       filepath.Base(filePath),
					Line:       fset.Position(x.Pos()).Line,
					Exported:   x.Name.IsExported(),
					References: references,
					Risk:       d.calculateRisk(x.Name.Name, "func", x.Name.IsExported()),
					Description: fmt.Sprintf("Function %s in package %s", x.Name.Name, packageName),
				})
			}
		case *ast.GenDecl:
			for _, spec := range x.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if !name.IsExported() {
							symbols = append(symbols, Symbol{
								Name:       name.Name,
								Kind:       "var",
								Package:    packageName,
								File:       filepath.Base(filePath),
								Line:       fset.Position(name.Pos()).Line,
								Exported:   name.IsExported(),
								References: 0,
								Risk:       d.calculateRisk(name.Name, "var", name.IsExported()),
								Description: fmt.Sprintf("Variable %s in package %s", name.Name, packageName),
							})
						}
					}
				case *ast.TypeSpec:
					if !s.Name.IsExported() {
						symbols = append(symbols, Symbol{
							Name:       s.Name.Name,
							Kind:       "type",
							Package:    packageName,
							File:       filepath.Base(filePath),
							Line:       fset.Position(s.Pos()).Line,
							Exported:   s.Name.IsExported(),
							References: 0,
							Risk:       d.calculateRisk(s.Name.Name, "type", s.Name.IsExported()),
							Description: fmt.Sprintf("Type %s in package %s", s.Name.Name, packageName),
						})
					}
				}
			}
		}
		return true
	})

	return symbols, nil
}

// calculateRisk calculates the risk level for removing a symbol
func (d *Detector) calculateRisk(name, kind string, exported bool) string {
	// High risk for exported symbols
	if exported {
		return "high"
	}

	// Medium risk for common patterns that might be used
	commonPatterns := []string{"init", "main", "New", "Create", "Build"}
	for _, pattern := range commonPatterns {
		if strings.Contains(name, pattern) {
			return "medium"
		}
	}

	// Low risk for private symbols with no obvious external use
	return "low"
}

// filterDeadSymbols filters out symbols that are actually used
func (d *Detector) filterDeadSymbols(symbols []Symbol) []Symbol {
	var deadSymbols []Symbol

	for _, symbol := range symbols {
		// For now, consider all unexported symbols with 0 references as dead
		// In a real implementation, this would do proper reference analysis
		if !symbol.Exported && symbol.References == 0 {
			deadSymbols = append(deadSymbols, symbol)
		}
	}

	return deadSymbols
}
