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
	x := 42 //nolint:short-var-decl
}`,
			expected: 0,
		},
		{
			name: "nolint_all_should_ignore_everything",
			code: `package main
func main() {
	x := 42 //nolint:all
	y := "test"
}`,
			expected: 1, // only y should be flagged
		},
		{
			name: "no_nolint_should_detect_everything",
			code: `package main
func main() {
	x := 42
	y := "test"
}`,
			expected: 2,
		},
		{
			name: "nolint_different_rule_should_not_ignore",
			code: `package main
func main() {
	x := 42 //nolint:if-init
}`,
			expected: 1, // should still be flagged for short-var-decl
		},
		{
			name: "multiple_rules_in_nolint",
			code: `package main
func main() {
	if err := someFunc(); err != nil { //nolint:short-var-decl,if-init
		return
	}
	x := 42 // nolint:short-var-decl,if-init
}
func someFunc() error { return nil }`,
			expected: 0, // tous ignorés
		},
	}

	type testCase struct {
		name     string
		code     string
		expected int
	}
	var tt testCase
	for _, tt = range tests {
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
			var rule types.Rule
			for _, rule = range linter.rules {
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
				var i int
				var issue types.Issue
				for i, issue = range allIssues {
					t.Logf("  Issue %d: %s [%s] at line %d", i+1, issue.Message, issue.Rule, issue.Line)
				}
				t.Logf("Issues after filtering:")
				for i, issue = range filteredIssues {
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
			name: "exact_rule_match_should_ignore",
			code: `package main
func main() {
	x := 42 //nolint:short-var-decl
}`,
			line:     3,
			ruleName: "short-var-decl",
			expected: true,
		},
		{
			name: "nolint_all_should_ignore",
			code: `package main
func main() {
	x := 42 //nolint:all
}`,
			line:     3,
			ruleName: "short-var-decl",
			expected: true,
		},
		{
			name: "different_rule_should_not_ignore",
			code: `package main
func main() {
	x := 42 //nolint:if-init
}`,
			line:     3,
			ruleName: "short-var-decl",
			expected: false,
		},
		{
			name: "no_comment_should_not_ignore",
			code: `package main
func main() {
	x := 42
}`,
			line:     3,
			ruleName: "short-var-decl",
			expected: false,
		},
		{
			name: "comment_different_line_should_not_ignore",
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

	type testCase struct {
		name     string
		code     string
		line     int
		ruleName string
		expected bool
	}
	var tt testCase
	for _, tt = range tests {
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
				t.Errorf("Expected %v, got %v for line %d and rule %s", tt.expected, result, tt.line, tt.ruleName)
			}
		})
	}
}

func TestFileNolintComments(t *testing.T) {
	var tests []struct {
		name     string
		code     string
		expected int // number of issues expected after file-level nolint filtering
	}
	tests = []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "file_nolint_all_should_ignore_everything",
			code: `//nolint
package main
func main() {
	x := 42
	y := "test"
	if err := someFunc(); err != nil {
		return
	}
}
func someFunc() error { return nil }`,
			expected: 0, // all issues ignored
		},
		{
			name: "file_nolint_with_space_should_ignore_everything",
			code: `// nolint
package main
func main() {
	x := 42
	y := "test"
}`,
			expected: 0, // all issues ignored
		},
		{
			name: "file_nolint_all_explicit_should_ignore_everything",
			code: `//nolint:all
package main
func main() {
	x := 42
	y := "test"
}`,
			expected: 0, // all issues ignored
		},
		{
			name: "file_nolint_specific_rule_should_ignore_only_that_rule",
			code: `//nolint:short-var-decl
package main
func main() {
	x := 42        // should be ignored
	var a = 33     // should be detected (var-no-type)
}`,
			expected: 1, // only var-no-type should be detected
		},
		{
			name: "file_nolint_multiple_rules_should_ignore_specified_rules",
			code: `//nolint:short-var-decl,var-no-type
package main
func main() {
	x := 42        // should be ignored (short-var-decl)
	var a = 33     // should be ignored (var-no-type)
	if err := someFunc(); err != nil { // should be detected (if-init)
		return
	}
}
func someFunc() error { return nil }`,
			expected: 1, // only if-init should be detected
		},
		{
			name: "no_file_nolint_should_detect_everything",
			code: `package main
func main() {
	x := 42
	y := "test"
}`,
			expected: 2, // both short-var-decl should be detected
		},
		{
			name: "file_nolint_after_package_should_not_work",
			code: `package main
//nolint
func main() {
	x := 42
	y := "test"
}`,
			expected: 2, // nolint after package declaration should not work for file-level
		},
		{
			name: "file_nolint_wrong_rule_should_not_ignore",
			code: `//nolint:if-init
package main
func main() {
	x := 42  // should be detected (short-var-decl not ignored)
}`,
			expected: 1, // short-var-decl should still be detected
		},
	}

	type testCase struct {
		name     string
		code     string
		expected int
	}
	var tt testCase
	for _, tt = range tests {
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
			var rule types.Rule
			for _, rule = range linter.rules {
				var issues []types.Issue
				issues = rule.Check(fset, file)
				allIssues = append(allIssues, issues...)
			}

			// Apply nolint filtering (both line-level and file-level)
			var filteredIssues []types.Issue
			filteredIssues = filterNolintIssues(allIssues, file, fset)

			if len(filteredIssues) != tt.expected {
				t.Errorf("Expected %d issues after file-level nolint filtering, got %d", tt.expected, len(filteredIssues))
				t.Logf("All issues before filtering:")
				var i int
				var issue types.Issue
				for i, issue = range allIssues {
					t.Logf("  Issue %d: %s [%s] at line %d", i+1, issue.Message, issue.Rule, issue.Line)
				}
				t.Logf("Issues after filtering:")
				for i, issue = range filteredIssues {
					t.Logf("  Issue %d: %s [%s] at line %d", i+1, issue.Message, issue.Rule, issue.Line)
				}
			}
		})
	}
}
