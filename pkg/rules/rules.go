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
