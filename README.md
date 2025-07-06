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

3. **If Initialization Rule (`if-init`)**
   - **Description**: Flags `if` statements with initializations. Such
     statements can be uncommon, unreadable, and disrupt the flow of the
     code.
   - **Detects**: `if err := someFunc(); err != nil`

## Ignoring Rules

You can ignore specific rules for certain lines of code using the
`//nolint` comment directive. Here's an example:

```go
x := 10 //nolint:short-var-decl
var a = 33 //nolint:var-no-type
```

In this example, the `short-var-decl` and `var-no-type` rules are ignored for these lines.

## Command Line Usage

To run the linter, use the following command:

```sh
go-syntax [paths...]
```

### Options

- `-v`: Enable verbose output.
- `-exit-code`: Set the exit code when issues are found. Defaults to `1`.
- `-c`: Enable or disable color output. Defaults to `true`.

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
```

The linter supports Go-style path patterns like `./...` for recursive analysis.

## Code

The linter is implemented with specific rules to ensure code quality:

- **ShortVarDeclRule**: Detects short variable declarations (`:=`).
- **VarNoTypeRule**: Detects variable declarations without explicit type.
- **IfInitRule**: Detects `if` statements with initializations.

The main function walks through the specified directory, lints each Go
file, and outputs any issues found.
