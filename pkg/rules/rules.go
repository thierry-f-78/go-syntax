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
				for _, spec := range node.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						// Check if type is not specified but values are provided
						if valueSpec.Type == nil && len(valueSpec.Values) > 0 {
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
				for _, field := range node.Type.Results.List {
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
