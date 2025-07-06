package linter

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/thierry-f-78/go-syntax/pkg/types"
)

func TestNolintComments(t *testing.T) {
	var tests []struct {
		name     string
		code     string
		expected int // number of issues expected after nolint filtering
	}
	tests = []struct {
		name     string
		code     string
		expected int // number of issues expected after nolint filtering
	}{
		{
			name: "nolint_by_rule_name_should_ignore",
			code: `package main
func main() {
	x := 42 // nolint:short-var-decl
}`,
			expected: 0,
		},
		{
			name: "nolint:all_should_ignore_all",
			code: `package main
func main() {
	x := 42 // nolint:all
}`,
			expected: 0,
		},
		{
			name: "no_nolint_should_detect",
			code: `package main
func main() {
	x := 42
}`,
			expected: 1,
		},
		{
			name: "nolint_wrong_rule_should_detect",
			code: `package main
func main() {
	x := 42 // nolint:if-init
}`,
			expected: 1,
		},
		{
			name: "nolint_wrong_code_should_detect",
			code: `package main
func main() {
	x := 42 // nolint:GS002
}`,
			expected: 1,
		},
		{
			name: "mixed_nolint_and_detection",
			code: `package main
func main() {
	x := 42 // nolint:short-var-decl
	y := "test"
	z := true // nolint:short-var-decl
}`,
			expected: 1, // seulement y := "test"
		},
		{
			name: "if_init_with_nolint_by_rule_name",
			code: `package main
func main() {
	var err error
	if err = someFunc(); err != nil { // nolint:if-init
		return
	}
}
func someFunc() error { return nil }`,
			expected: 0,
		},
		{
			name: "multiple_rules_same_line_nolint:all",
			code: `package main
func main() {
	if x := getValue(); x != nil { // nolint:all
		return
	}
}
func getValue() interface{} { return nil }`,
			expected: 0, // devrait ignorer les deux règles (GS001 et GS002)
		},
		{
			name: "multiple_rules_same_line_specific_nolint",
			code: `package main
func main() {
	if x := getValue(); x != nil { // nolint:short-var-decl
		return
	}
}
func getValue() interface{} { return nil }`,
			expected: 1, // ignore GS001 mais détecte GS002
		},
		{
			name: "nolint_on_different_line_should_not_affect",
			code: `package main
func main() {
	// nolint:GS001
	x := 42
}`,
			expected: 1, // nolint pas sur la même ligne
		},
		{
			name: "multiple_nolint_formats",
			code: `package main
func main() {
	x := 42 // nolint:short-var-decl,if-init
}
func someFunc() error { return nil }`,
			expected: 0, // tous ignorés
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var linter *Linter
			linter = New()

			// Create temporary file with code
			var fset *token.FileSet
			fset = token.NewFileSet()
			var file *ast.File
			var err error
			file, err = parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			// Get issues from rules
			var allIssues []types.Issue
			for _, rule := range linter.rules {
				var issues []types.Issue
				issues = rule.Check(fset, file)
				allIssues = append(allIssues, issues...)
			}

			// Apply nolint filtering
			var filteredIssues []types.Issue
			filteredIssues = filterNolintIssues(allIssues, file, fset)

			if len(filteredIssues) != tt.expected {
				t.Errorf("Expected %d issues after nolint filtering, got %d", tt.expected, len(filteredIssues))
				t.Logf("All issues before filtering:")
				for i, issue := range allIssues {
					t.Logf("  Issue %d: %s [%s] at line %d", i+1, issue.Message, issue.Rule, issue.Line)
				}
				t.Logf("Issues after filtering:")
				for i, issue := range filteredIssues {
					t.Logf("  Issue %d: %s [%s] at line %d", i+1, issue.Message, issue.Rule, issue.Line)
				}
			}
		})
	}
}

func TestIsNolintComment(t *testing.T) {
	var tests []struct {
		name     string
		code     string
		line     int
		ruleName string
		expected bool
	}
	tests = []struct {
		name     string
		code     string
		line     int
		ruleName string
		expected bool
	}{
		{
			name: "exact_rule_name_match",
			code: `package main
func main() {
	x := 42 // nolint:short-var-decl
}`,
			line:     3,
			ruleName: "short-var-decl",
			expected: true,
		},
		{
			name: "nolint_all_match",
			code: `package main
func main() {
	x := 42 // nolint:all
}`,
			line:     3,
			ruleName: "short-var-decl",
			expected: true,
		},
		{
			name: "no_nolint_comment",
			code: `package main
func main() {
	x := 42
}`,
			line:     3,
			ruleName: "short-var-decl",
			expected: false,
		},
		{
			name: "wrong_rule_name",
			code: `package main
func main() {
	x := 42 // nolint:if-init
}`,
			line:     3,
			ruleName: "short-var-decl",
			expected: false,
		},
		{
			name: "wrong_code",
			code: `package main
func main() {
	x := 42 // nolint:GS002
}`,
			line:     3,
			ruleName: "short-var-decl",
			expected: false,
		},
		{
			name: "multiple_rules_in_comment",
			code: `package main
func main() {
	x := 42 // nolint:short-var-decl,if-init
}`,
			line:     3,
			ruleName: "if-init",
			expected: true,
		},
		{
			name: "comment_on_different_line",
			code: `package main
func main() {
	// nolint:short-var-decl
	x := 42
}`,
			line:     4,
			ruleName: "short-var-decl",
			expected: false, // comment pas sur la même ligne
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fset *token.FileSet
			fset = token.NewFileSet()
			var file *ast.File
			var err error
			file, err = parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			var result bool
			result = isNolintComment(file, tt.line, tt.ruleName, fset)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
