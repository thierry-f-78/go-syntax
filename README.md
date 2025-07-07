# Go-Syntax Linter

## Installation

To install the Go-Syntax linter, use the following command:

```sh
go install github.com/thierry-f-78/go-syntax/cmd/go-syntax@latest
```

## Purpose

I prefer explicit typing and clear code structure, which is why I
dislike implicit typing and statements within `if` initializations. This
linter helps detect such patterns in Go code, ensuring better
readability and maintainability.

## Detection Rules

The linter includes rules to detect and flag specific code patterns:

1. **Short Variable Declaration Rule (`short-var-decl`)**
   - **Description**: Avoids the use of short variable declarations
     (`:=`) outside of type switches. Implicit types can make code
     reviews harder and bugs more likely.
   - **Detects**: `x := 42`, `for i := range items`, `for i, v := range items`

2. **Variable Without Type Rule (`var-no-type`)**
   - **Description**: Flags variable declarations without explicit type
     when a value is provided. Implicit types can make code reviews
     harder and bugs more likely.
   - **Detects**: `var a = 33`, `var r = strings.Split("a,b", ",")`
   - **Exception**: Allows declarations where the type is explicit in the value: `var x = []int{1, 2}`, `var a = make([]int, 0)`, `var b = x.(int)`, or unambiguous literals: `var s = "hello"`, `var b = true`

3. **Constant Without Type Rule (`const-no-type`)**
   - **Description**: Flags constant declarations without explicit type
     when a value is provided. Implicit types can make code reviews
     harder and bugs more likely.
   - **Detects**: `const BufferSize = 1024`, `const Pi = 3.14159`
   - **Exception**: Allows unambiguous literals: `const Name = "app"`, `const Debug = true`

4. **Named Returns Rule (`named-returns`)**
   - **Description**: Flags functions with named return parameters.
     Named returns make it unclear what values are returned and can
     lead to confusion during code reviews.
   - **Detects**: `func divide(a, b int) (result int, err error)`

5. **Naked Return Rule (`naked-return`)**
   - **Description**: Flags naked returns in functions with named return
     parameters. Naked returns make it unclear what values are being
     returned without checking the function signature.
   - **Detects**: `return` (without explicit values) in functions with named returns

6. **If Initialization Rule (`if-init`)**
   - **Description**: Flags `if` statements with initializations. Such
     statements can be uncommon, unreadable, and disrupt the flow of the
     code.
   - **Detects**: `if err := someFunc(); err != nil`

## Ignoring Rules

You can ignore specific rules using the `//nolint` comment directive in two ways:

### Line-Level Ignoring

Ignore rules for specific lines by adding comments at the end of the line:

```go
x := 10 //nolint:short-var-decl
var a = 33 //nolint:var-no-type
const BufferSize = 1024 //nolint:const-no-type
func recover() (err error) { //nolint:named-returns
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    return //nolint:naked-return
}
```

### File-Level Ignoring

Ignore rules for entire files by adding comments at the top of the file (before the `package` declaration):

```go
//nolint:short-var-decl,var-no-type
package main

func main() {
    // These violations will be ignored for the entire file
    x := 42
    var a = 33

    // This will still be detected (if-init not in nolint list)
    if err := someFunc(); err != nil {
        return
    }
}
```

You can also disable all rules for a file:

```go
//nolint
package main

func main() {
    // All violations ignored for this file
    x := 42
    var a = 33
    if err := someFunc(); err != nil {
        return
    }
}
```

The `named-returns` and `naked-return` rules are commonly ignored together for panic recovery patterns.

## Command Line Usage

To run the linter, use the following command:

```sh
go-syntax [paths...]
```

### Options

- `-v`: Enable verbose output.
- `-exit-code`: Set the exit code when issues are found. Defaults to `1`.
- `-c`: Enable or disable color output. Defaults to `true`.
- `-e <pattern>`: Exclude files matching pattern. Can be repeated multiple times.

### Examples

```sh
# Analyze current directory
go-syntax

# Analyze specific directory recursively
go-syntax ./...

# Analyze multiple paths
go-syntax ./cmd/... ./pkg/...

# Analyze with verbose output
go-syntax -v ./...

# Exclude specific files
go-syntax -e "generated.go" ./...

# Exclude files by pattern
go-syntax -e "*_test.go" ./...

# Multiple exclude patterns
go-syntax -e "generated.go" -e "*_mock.go" -e "vendor/*" ./...

# Exclude with wildcard patterns
go-syntax -e "*.pb.go" -e "wire_gen.go" ./...
```

The linter supports Go-style path patterns like `./...` for recursive analysis.

## File Exclusion

The `-e` flag allows you to exclude files from analysis using patterns:

### Pattern Types

1. **Exact filename**: `-e "generated.go"`
2. **Wildcard patterns**: `-e "*.pb.go"` (protocol buffer files)
3. **Path patterns**: `-e "vendor/*"` (anything in vendor directory)
4. **Multiple patterns**: `-e "*.pb.go" -e "*_mock.go" -e "generated.go"`

### Common Use Cases

```sh
# Exclude generated files
go-syntax -e "*.pb.go" -e "*_gen.go" ./...

# Exclude test files
go-syntax -e "*_test.go" ./...

# Exclude vendor and generated files
go-syntax -e "vendor/*" -e "*.pb.go" -e "wire_gen.go" ./...

# Exclude specific directories
go-syntax -e "testdata/*" -e "examples/*" ./...
```

## Code

The linter is implemented with specific rules to ensure code quality:

- **ShortVarDeclRule**: Detects short variable declarations (`:=`).
- **VarNoTypeRule**: Detects variable declarations without explicit type.
- **ConstNoTypeRule**: Detects constant declarations without explicit type.
- **NamedReturnsRule**: Detects functions with named return parameters.
- **NakedReturnRule**: Detects naked returns in functions with named parameters.
- **IfInitRule**: Detects `if` statements with initializations.

The main function walks through the specified directory, lints each Go
file, and outputs any issues found.
