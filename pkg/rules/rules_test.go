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

	var tt struct {
		name     string
		code     string
		expected int
	}
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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				var i int
				var issue types.Issue
				for i, issue = range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			var issue types.Issue
			for _, issue = range issues {
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
			name: "multiple var without type - should detect only ambiguous ones",
			code: `package main
func main() {
	var a = 33
	var b = "test"
	var c = true
}`,
			expected: 1, // only 'a' (b and c are unambiguous literals)
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
			name: "mixed var declarations - should detect only ambiguous ones",
			code: `package main
func main() {
	var a = 33        // should detect (int literal ambiguous)
	var b int = 42    // should not detect (explicit type)
	var c string      // should not detect (no value)
	var d = "test"    // should not detect (string literal unambiguous)
}`,
			expected: 1, // only 'a'
		},
		{
			name: "var with composite literal - should not detect",
			code: `package main
func main() {
	var a = []int{1, 2, 3}
	var b = map[string]int{"key": 1}
	var c = struct{ x int }{x: 1}
}`,
			expected: 0,
		},
		{
			name: "var with make/new - should not detect",
			code: `package main
func main() {
	var a = make([]int, 0)
	var b = new(int)
	var c = make(map[string]int)
}`,
			expected: 0,
		},
		{
			name: "var with type assertion - should not detect",
			code: `package main
func main() {
	var x interface{} = 42
	var a = x.(int)
}`,
			expected: 0,
		},
		{
			name: "var with string literal - should not detect",
			code: `package main
func main() {
	var a = "hello"
	var b = ` + "`world`" + `
}`,
			expected: 0,
		},
		{
			name: "var with bool literal - should not detect",
			code: `package main
func main() {
	var a = true
	var b = false
}`,
			expected: 0,
		},
	}

	var rule *VarNoTypeRule
	rule = &VarNoTypeRule{}

	var tt struct {
		name     string
		code     string
		expected int
	}
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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				var i int
				var issue types.Issue
				for i, issue = range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			var issue types.Issue
			for _, issue = range issues {
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

	var tt struct {
		name     string
		code     string
		expected int
	}
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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				var i int
				var issue types.Issue
				for i, issue = range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			var issue types.Issue
			for _, issue = range issues {
				if issue.Rule != "named-returns" {
					t.Errorf("Expected rule 'named-returns', got %s", issue.Rule)
				}
			}
		})
	}
}

func TestNakedReturnRule(t *testing.T) {
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
			name: "naked return with named parameters - should detect",
			code: `package main
import "fmt"
func divide(a, b int) (result int, err error) {
	if b == 0 {
		err = fmt.Errorf("division by zero")
		return
	}
	result = a / b
	return
}`,
			expected: 2, // two naked returns
		},
		{
			name: "single naked return - should detect",
			code: `package main
func getValue() (value int) {
	value = 42
	return
}`,
			expected: 1,
		},
		{
			name: "explicit return with named parameters - should not detect",
			code: `package main
import "fmt"
func divide(a, b int) (result int, err error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}`,
			expected: 0,
		},
		{
			name: "naked return without named parameters - should not detect",
			code: `package main
func doSomething() {
	println("hello")
	return
}`,
			expected: 0,
		},
		{
			name: "explicit return without named parameters - should not detect",
			code: `package main
func divide(a, b int) (int, error) {
	return a / b, nil
}`,
			expected: 0,
		},
		{
			name: "mixed returns - should detect only naked ones",
			code: `package main
import "fmt"
func process() (result int, err error) {
	someCondition := true
	if someCondition {
		return 0, fmt.Errorf("error")  // explicit - ok
	}
	result = 42
	return  // naked - should detect
}`,
			expected: 1,
		},
		{
			name: "multiple functions - should detect all naked returns",
			code: `package main
func first() (a int) {
	a = 1
	return  // naked - should detect
}
func second() (int) {
	return 2  // explicit - ok
}
func third() (result bool) {
	result = true
	return  // naked - should detect
}`,
			expected: 2, // first and third functions
		},
		{
			name: "method with naked return - should detect",
			code: `package main
type MyStruct struct{}
func (m MyStruct) getValue() (value int) {
	value = 42
	return
}`,
			expected: 1,
		},
	}

	var rule *NakedReturnRule
	rule = &NakedReturnRule{}

	var tt struct {
		name     string
		code     string
		expected int
	}
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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				var i int
				var issue types.Issue
				for i, issue = range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			var issue types.Issue
			for _, issue = range issues {
				if issue.Rule != "naked-return" {
					t.Errorf("Expected rule 'naked-return', got %s", issue.Rule)
				}
			}
		})
	}
}

func TestConstNoTypeRule(t *testing.T) {
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
			name: "const without type - should detect",
			code: `package main
const BufferSize = 1024`,
			expected: 1,
		},
		{
			name: "const string without type - should not detect (unambiguous literal)",
			code: `package main
const AppName = "myapp"`,
			expected: 0,
		},
		{
			name: "const float without type - should detect",
			code: `package main
const Pi = 3.14159`,
			expected: 1,
		},
		{
			name: "multiple const without type - should detect only ambiguous ones",
			code: `package main
const (
	BufferSize = 1024
	AppName = "myapp"
	Pi = 3.14159
)`,
			expected: 2, // BufferSize and Pi (AppName is string literal)
		},
		{
			name: "const with explicit type - should not detect",
			code: `package main
const BufferSize int = 1024
const AppName string = "myapp"
const Pi float64 = 3.14159`,
			expected: 0,
		},
		{
			name: "const without value - should not detect",
			code: `package main
const (
	Red = iota
	Green
	Blue
)`,
			expected: 1, // only Red has a value (iota)
		},
		{
			name: "const with type and value - should not detect",
			code: `package main
const (
	BufferSize int = 1024
	AppName string = "myapp"
	Enabled bool = true
)`,
			expected: 0,
		},
		{
			name: "mixed const declarations - should detect only ambiguous ones",
			code: `package main
const (
	Size = 100          // should detect (int literal ambiguous)
	Name string = "app" // should not detect (explicit type)
	Count int = 5       // should not detect (explicit type)
	Value = 42          // should detect (int literal ambiguous)
	Message = "hello"   // should not detect (string literal unambiguous)
	Debug = true        // should not detect (bool literal unambiguous)
)`,
			expected: 2, // Size and Value
		},
		{
			name: "const with complex expressions - should detect",
			code: `package main
const MaxUsers = 100 * 1024
const Timeout = 30 * time.Second`,
			expected: 2,
		},
		{
			name: "const bool without type - should not detect (unambiguous literal)",
			code: `package main
const Debug = true
const Enabled = false`,
			expected: 0,
		},
		{
			name: "const with string and bool literals - should not detect",
			code: `package main
const (
	Name = "test"
	Debug = true
	Version = "1.0"
	Enabled = false
)`,
			expected: 0,
		},
	}

	var rule *ConstNoTypeRule
	rule = &ConstNoTypeRule{}

	var tt struct {
		name     string
		code     string
		expected int
	}
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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				var i int
				var issue types.Issue
				for i, issue = range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			var issue types.Issue
			for _, issue = range issues {
				if issue.Rule != "const-no-type" {
					t.Errorf("Expected rule 'const-no-type', got %s", issue.Rule)
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

	var tt struct {
		name     string
		code     string
		expected int
	}
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

			var issues []types.Issue
			issues = rule.Check(fset, file)
			if len(issues) != tt.expected {
				t.Errorf("Expected %d issues, got %d", tt.expected, len(issues))
				var i int
				var issue types.Issue
				for i, issue = range issues {
					t.Logf("Issue %d: %s at line %d", i+1, issue.Message, issue.Line)
				}
			}

			// Verify all issues have correct code and rule name
			var issue types.Issue
			for _, issue = range issues {
				if issue.Rule != "if-init" {
					t.Errorf("Expected rule 'if-init', got %s", issue.Rule)
				}
			}
		})
	}
}
