package types

import (
	"go/ast"
	"go/token"
)

type Issue struct {
	File        string
	Line        int
	Column      int
	Message     string
	Description string
	Rule        string
}

type Rule interface {
	Name() string
	Check(fset *token.FileSet, file *ast.File) []Issue
}
