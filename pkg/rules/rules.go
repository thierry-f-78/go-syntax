package rules

import (
	"go/ast"
	"go/token"

	"github.com/thierry-f-78/go-syntax/pkg/types"
)

type ShortVarDeclRule struct{}

func (r *ShortVarDeclRule) Name() string {
	return "short-var-decl"
}

func (r *ShortVarDeclRule) Check(fset *token.FileSet, file *ast.File) []types.Issue {
	var issues []types.Issue

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.AssignStmt:
			if node.Tok == token.DEFINE {
				if !r.isInTypeSwitch(node, file) {
					var pos token.Position
					pos = fset.Position(node.Pos())
					issues = append(issues, types.Issue{
						File:        pos.Filename,
						Line:        pos.Line,
						Column:      pos.Column,
						Message:     "Short variable declaration ':=' is not allowed",
						Description: "Avoid ':=': unclear types make reviews harder, bugs likelier.",
						Rule:        r.Name(),
					})
				}
			}
		case *ast.RangeStmt:
			if node.Tok == token.DEFINE {
				var pos token.Position
				pos = fset.Position(node.Pos())
				issues = append(issues, types.Issue{
					File:        pos.Filename,
					Line:        pos.Line,
					Column:      pos.Column,
					Message:     "Short variable declaration ':=' is not allowed in range",
					Description: "Avoid ':=': unclear types make reviews harder, bugs likelier.",
					Rule:        r.Name(),
				})
			}
		}
		return true
	})

	return issues
}

func (r *ShortVarDeclRule) isInTypeSwitch(assign *ast.AssignStmt, file *ast.File) bool {
	var found bool

	ast.Inspect(file, func(n ast.Node) bool {
		var typeSwitch *ast.TypeSwitchStmt
		var ok bool
		typeSwitch, ok = n.(*ast.TypeSwitchStmt)
		if ok {
			if typeSwitch.Assign == assign {
				found = true
				return false
			}
		}
		return true
	})

	return found
}

// hasExplicitType checks if an expression contains an explicit type
func hasExplicitType(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.CompositeLit:
		// []int{1, 2}, map[string]int{}, struct{}{}
		return e.Type != nil
	case *ast.CallExpr:
		if ident, ok := e.Fun.(*ast.Ident); ok {
			// make([]int, 0), new(int)
			return ident.Name == "make" || ident.Name == "new"
		}
	case *ast.TypeAssertExpr:
		// x.(int)
		return true
	case *ast.FuncLit:
		// func(int) error { ... }
		return true
	}
	return false
}

// isUnambiguousLiteral checks if an expression is an unambiguous literal (string or bool)
func isUnambiguousLiteral(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		// String literals: "hello", `hello`
		// Bool literals: true, false (but these are actually *ast.Ident)
		return e.Kind == token.STRING
	case *ast.Ident:
		// Bool literals: true, false
		return e.Name == "true" || e.Name == "false"
	}
	return false
}

// isTypeExpression checks if an expression represents a type (like struct{}, int, etc.)
func isTypeExpression(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.StructType:
		// struct{}, struct{x int}
		return true
	case *ast.ArrayType:
		// []int, [5]int
		return true
	case *ast.MapType:
		// map[string]int
		return true
	case *ast.ChanType:
		// chan int, <-chan int
		return true
	case *ast.FuncType:
		// func(), func(int) string
		return true
	case *ast.InterfaceType:
		// interface{}, interface{Method()}
		return true
	case *ast.StarExpr:
		// *int, *struct{}
		return isTypeExpression(e.X)
	case *ast.SelectorExpr:
		// pkg.Type
		return true
	case *ast.Ident:
		// Basic types: int, string, bool, etc.
		// Note: This might include variable names too, but in the context of
		// "var x = TYPE", if it's a valid Go program, TYPE should be a type
		return true
	}
	return false
}

type VarNoTypeRule struct{}

func (r *VarNoTypeRule) Name() string {
	return "var-no-type"
}

func (r *VarNoTypeRule) Check(fset *token.FileSet, file *ast.File) []types.Issue {
	var issues []types.Issue

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			if node.Tok == token.VAR {
				var spec ast.Spec
				for _, spec = range node.Specs {
					var valueSpec *ast.ValueSpec
					var ok bool
					valueSpec, ok = spec.(*ast.ValueSpec)
					if ok {
						// Check if type is not specified but values are provided
						// Exception: allow when the value has an explicit type, is an unambiguous literal, or is a type expression
						if valueSpec.Type == nil && len(valueSpec.Values) > 0 && !hasExplicitType(valueSpec.Values[0]) && !isUnambiguousLiteral(valueSpec.Values[0]) && !isTypeExpression(valueSpec.Values[0]) {
							var pos token.Position
							pos = fset.Position(valueSpec.Pos())
							issues = append(issues, types.Issue{
								File:        pos.Filename,
								Line:        pos.Line,
								Column:      pos.Column,
								Message:     "Variable declaration without explicit type is not allowed",
								Description: "Avoid 'var x = value': unclear types make reviews harder, bugs likelier.",
								Rule:        r.Name(),
							})
						}
					}
				}
			}
		}
		return true
	})

	return issues
}

type NamedReturnsRule struct{}

func (r *NamedReturnsRule) Name() string {
	return "named-returns"
}

func (r *NamedReturnsRule) Check(fset *token.FileSet, file *ast.File) []types.Issue {
	var issues []types.Issue

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			if node.Type.Results != nil && len(node.Type.Results.List) > 0 {
				var field *ast.Field
				for _, field = range node.Type.Results.List {
					// Check if any return parameter has a name
					if len(field.Names) > 0 {
						var pos token.Position
						pos = fset.Position(field.Pos())
						issues = append(issues, types.Issue{
							File:        pos.Filename,
							Line:        pos.Line,
							Column:      pos.Column,
							Message:     "Named return parameters are not allowed",
							Description: "Avoid named returns: unclear what is returned, harder to review.",
							Rule:        r.Name(),
						})
					}
				}
			}
		}
		return true
	})

	return issues
}

type NakedReturnRule struct{}

func (r *NakedReturnRule) Name() string {
	return "naked-return"
}

func (r *NakedReturnRule) Check(fset *token.FileSet, file *ast.File) []types.Issue {
	var issues []types.Issue

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ReturnStmt:
			// Check if it's a naked return (no explicit values)
			if len(node.Results) == 0 {
				// Find the containing function to check if it has named returns
				var containingFunc *ast.FuncDecl
				ast.Inspect(file, func(fn ast.Node) bool {
					var funcDecl *ast.FuncDecl
					var ok bool
					funcDecl, ok = fn.(*ast.FuncDecl)
					if ok {
						// Check if the return statement is within this function
						if funcDecl.Pos() <= node.Pos() && node.Pos() <= funcDecl.End() {
							containingFunc = funcDecl
							return false
						}
					}
					return true
				})

				// If we found a containing function and it has named returns, flag it
				if containingFunc != nil && containingFunc.Type.Results != nil {
					var hasNamedReturns bool = false
					var field *ast.Field
					for _, field = range containingFunc.Type.Results.List {
						if len(field.Names) > 0 {
							hasNamedReturns = true
							break
						}
					}

					if hasNamedReturns {
						var pos token.Position
						pos = fset.Position(node.Pos())
						issues = append(issues, types.Issue{
							File:        pos.Filename,
							Line:        pos.Line,
							Column:      pos.Column,
							Message:     "Naked return is not allowed",
							Description: "Avoid naked returns: unclear what values are returned.",
							Rule:        r.Name(),
						})
					}
				}
			}
		}
		return true
	})

	return issues
}

type ConstNoTypeRule struct{}

func (r *ConstNoTypeRule) Name() string {
	return "const-no-type"
}

func (r *ConstNoTypeRule) Check(fset *token.FileSet, file *ast.File) []types.Issue {
	var issues []types.Issue

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			if node.Tok == token.CONST {
				var spec ast.Spec
				for _, spec = range node.Specs {
					var valueSpec *ast.ValueSpec
					var ok bool
					valueSpec, ok = spec.(*ast.ValueSpec)
					if ok {
						// Check if type is not specified but values are provided
						// Exception: allow when the value is an unambiguous literal
						if valueSpec.Type == nil && len(valueSpec.Values) > 0 && !isUnambiguousLiteral(valueSpec.Values[0]) {
							var pos token.Position
							pos = fset.Position(valueSpec.Pos())
							issues = append(issues, types.Issue{
								File:        pos.Filename,
								Line:        pos.Line,
								Column:      pos.Column,
								Message:     "Constant declaration without explicit type is not allowed",
								Description: "Avoid 'const x = value': unclear types make reviews harder, bugs likelier.",
								Rule:        r.Name(),
							})
						}
					}
				}
			}
		}
		return true
	})

	return issues
}

type IfInitRule struct{}

func (r *IfInitRule) Name() string {
	return "if-init"
}

func (r *IfInitRule) Check(fset *token.FileSet, file *ast.File) []types.Issue {
	var issues []types.Issue

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.IfStmt:
			if node.Init != nil {
				var pos token.Position
				pos = fset.Position(node.Pos())
				issues = append(issues, types.Issue{
					File:        pos.Filename,
					Line:        pos.Line,
					Column:      pos.Column,
					Message:     "If statement with initialization is not allowed.",
					Description: "Avoid 'if stmt; cond': uncommon, unreadable, breaks flow.",
					Rule:        r.Name(),
				})
			}
		}
		return true
	})

	return issues
}
