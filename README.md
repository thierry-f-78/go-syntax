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

2. **If Initialization Rule (`if-init`)**
   - **Description**: Flags `if` statements with initializations. Such
     statements can be uncommon, unreadable, and disrupt the flow of the
     code.

## Ignoring Rules

You can ignore specific rules for certain lines of code using the
`//nolint` comment directive. Here's an example:

```go
x := 10 //nolint:short-var-decl
```

In this example, the `short-var-decl` rule is ignored for this line.

## Command Line Usage

To run the linter, use the following command:

```sh
go-syntax -path <directory-path>
```

### Options

- `-path`: Specify the directory path to analyze. Defaults to the
  current directory (`.`).
- `-v`: Enable verbose output.
- `-exit-code`: Set the exit code when issues are found. Defaults to `1`.
- `-c`: Enable or disable color output. Defaults to `true`.

### Example

```sh
go-syntax -path ./my-go-project -v
```

This command will analyze all Go files in the `./my-go-project`
directory and output verbose results.

## Code

The linter is implemented with specific rules to ensure code quality:

- **ShortVarDeclRule**: Detects short variable declarations.
- **IfInitRule**: Detects `if` statements with initializations.

The main function walks through the specified directory, lints each Go
file, and outputs any issues found.
