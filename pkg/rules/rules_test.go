package rules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/thierry-f-78/go-syntax/pkg/types"
)

func TestShortVarDeclRule(t *testing.T) {
	var tests []struct {
		name     string
		code     string
		expected int // number of issues expected
	}
	tests = []struct {
		name     string
		code     string
		expected int // number of issues expected
	}{
		{
			name: "simple short var decl - should detect",
			code: `package main
func main() {
	x := 42
}`,
			expected: 1,
		},
		{
			name: "multiple short var decl - should detect all",
			code: `package main
func main() {
	x := 42
	y := "test"
	z := true
}`,
			expected: 3,
		},
		{
			name: "regular assignment - should not detect",
			code: `package main
func main() {
	var x int
	x = 42
}`,
			expected: 0,
		},
		{
			name: "type switch - should not detect",
			code: `package main
func main() {
	var i interface{} = "test"
	switch v := i.(type) {
	case string:
		_ = v
	}
}`,
			expected: 0,
		},
		{
			name: "slice creation with short var - should detect",
			code: `package main
func main() {
	items := []int{1, 2, 3}
	_ = items
}`,
			expected: 1,
		},
		{
			name: "for range with short var - should detect",
			code: `package main
func main() {
	items := []int{1, 2, 3}
	for i, v := range items {
		_, _ = i, v
	}
}`,
			expected: 2, // items := and for i, v := range
		},
		{
			name: "function call with short var - should detect",
			code: `package main
func main() {
	result := someFunc()
	_ = result
}
func someFunc() string { return "" }`,
			expected: 1,
		},
	}

	var rule *ShortVarDeclRule
	rule = &ShortVarDeclRule{}

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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				for i, issue := range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			for _, issue := range issues {
				if issue.Rule != "short-var-decl" {
					t.Errorf("Expected rule 'short-var-decl', got %s", issue.Rule)
				}
			}
		})
	}
}

func TestVarNoTypeRule(t *testing.T) {
	var tests []struct {
		name     string
		code     string
		expected int
	}
	tests = []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "var without type - should detect",
			code: `package main
func main() {
	var a = 33
}`,
			expected: 1,
		},
		{
			name: "var with function call - should detect",
			code: `package main
import "strings"
func main() {
	var r = strings.Split("a,b", ",")
}`,
			expected: 1,
		},
		{
			name: "multiple var without type - should detect all",
			code: `package main
func main() {
	var a = 33
	var b = "test"
	var c = true
}`,
			expected: 3,
		},
		{
			name: "var with explicit type - should not detect",
			code: `package main
func main() {
	var a int = 33
	var b string = "test"
}`,
			expected: 0,
		},
		{
			name: "var without value - should not detect",
			code: `package main
func main() {
	var a int
	var b string
}`,
			expected: 0,
		},
		{
			name: "var with type and value - should not detect",
			code: `package main
import "strings"
func main() {
	var r []string = strings.Split("a,b", ",")
}`,
			expected: 0,
		},
		{
			name: "mixed var declarations - should detect only those without type",
			code: `package main
func main() {
	var a = 33        // should detect
	var b int = 42    // should not detect
	var c string      // should not detect
	var d = "test"    // should detect
}`,
			expected: 2,
		},
	}

	var rule *VarNoTypeRule
	rule = &VarNoTypeRule{}

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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				for i, issue := range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			for _, issue := range issues {
				if issue.Rule != "var-no-type" {
					t.Errorf("Expected rule 'var-no-type', got %s", issue.Rule)
				}
			}
		})
	}
}

func TestNamedReturnsRule(t *testing.T) {
	var tests []struct {
		name     string
		code     string
		expected int
	}
	tests = []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "named return parameters - should detect",
			code: `package main
func divide(a, b int) (result int, err error) {
	return a / b, nil
}`,
			expected: 2, // result and err
		},
		{
			name: "single named return - should detect",
			code: `package main
func getValue() (value int) {
	return 42
}`,
			expected: 1,
		},
		{
			name: "multiple named returns - should detect all",
			code: `package main
func process() (result int, err error) {
	return 42, nil
}`,
			expected: 2, // result and err are both named
		},
		{
			name: "unnamed return parameters - should not detect",
			code: `package main
func divide(a, b int) (int, error) {
	return a / b, nil
}`,
			expected: 0,
		},
		{
			name: "function with no return - should not detect",
			code: `package main
func doSomething() {
	println("hello")
}`,
			expected: 0,
		},
		{
			name: "method with named returns - should detect",
			code: `package main
type MyStruct struct{}
func (m MyStruct) calculate() (result int, err error) {
	return 42, nil
}`,
			expected: 2,
		},
		{
			name: "multiple functions - should detect all named returns",
			code: `package main
func first() (a int, b string) {
	return 1, "test"
}
func second() (int, string) {
	return 2, "ok"
}
func third() (result bool) {
	return true
}`,
			expected: 3, // a, b, result
		},
	}

	var rule *NamedReturnsRule
	rule = &NamedReturnsRule{}

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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				for i, issue := range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			for _, issue := range issues {
				if issue.Rule != "named-returns" {
					t.Errorf("Expected rule 'named-returns', got %s", issue.Rule)
				}
			}
		})
	}
}

func TestIfInitRule(t *testing.T) {
	var tests []struct {
		name     string
		code     string
		expected int
	}
	tests = []struct {
		name     string
		code     string
		expected int
	}{
		{
			name: "if with assignment - should detect",
			code: `package main
func main() {
	var err error
	if err = someFunc(); err != nil {
		return
	}
}
func someFunc() error { return nil }`,
			expected: 1,
		},
		{
			name: "if with short var decl - should detect",
			code: `package main
func main() {
	if err := someFunc(); err != nil {
		return
	}
}
func someFunc() error { return nil }`,
			expected: 1,
		},
		{
			name: "simple if without init - should not detect",
			code: `package main
func main() {
	var err error
	if err != nil {
		return
	}
}`,
			expected: 0,
		},
		{
			name: "if with assignment but no nil check - should detect",
			code: `package main
func main() {
	var count int
	if count = getCount(); count > 0 {
		return
	}
}
func getCount() int { return 5 }`,
			expected: 1,
		},
		{
			name: "multiple if with init - should detect all",
			code: `package main
func main() {
	var err error
	var result interface{}
	if err = someFunc(); err != nil {
		return
	}
	if result = getResult(); result != nil {
		return
	}
}
func someFunc() error { return nil }
func getResult() interface{} { return nil }`,
			expected: 2,
		},
		{
			name: "nested if with init - should detect",
			code: `package main
func main() {
	if true {
		if err := someFunc(); err != nil {
			return
		}
	}
}
func someFunc() error { return nil }`,
			expected: 1,
		},
	}

	var rule *IfInitRule
	rule = &IfInitRule{}

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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				for i, issue := range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			for _, issue := range issues {
				if issue.Rule != "if-init" {
					t.Errorf("Expected rule 'if-init', got %s", issue.Rule)
				}
			}
		})
	}
}
