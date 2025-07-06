package linter

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/thierry-f-78/go-syntax/pkg/rules"
	"github.com/thierry-f-78/go-syntax/pkg/types"
)

type Linter struct {
	rules []types.Rule
}

func New() *Linter {
	return &Linter{
		rules: []types.Rule{
			&rules.ShortVarDeclRule{},
			&rules.VarNoTypeRule{},
			&rules.NamedReturnsRule{},
			&rules.IfInitRule{},
		},
	}
}

func (l *Linter) Lint(files []string) []types.Issue {
	var allIssues []types.Issue

	for _, file := range files {
		var issues []types.Issue
		issues = l.lintFile(file)
		allIssues = append(allIssues, issues...)
	}

	return allIssues
}

func (l *Linter) lintFile(filename string) []types.Issue {
	var fset *token.FileSet
	fset = token.NewFileSet()

	var src *ast.File
	var err error
	src, err = parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return []types.Issue{{
			File:    filename,
			Line:    1,
			Column:  1,
			Message: "Parse error: " + err.Error(),
			Rule:    "parse",
		}}
	}

	var issues []types.Issue

	for _, rule := range l.rules {
		var ruleIssues []types.Issue
		ruleIssues = rule.Check(fset, src)
		issues = append(issues, ruleIssues...)
	}

	return filterNolintIssues(issues, src, fset)
}

func filterNolintIssues(issues []types.Issue, file *ast.File, fset *token.FileSet) []types.Issue {
	var filtered []types.Issue

	for _, issue := range issues {
		if !isNolintComment(file, issue.Line, issue.Rule, fset) {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

func isNolintComment(file *ast.File, line int, ruleName string, fset *token.FileSet) bool {
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			var commentPos token.Position
			commentPos = fset.Position(comment.Pos())

			if commentPos.Line == line && strings.Contains(comment.Text, "nolint") {
				if strings.Contains(comment.Text, "nolint:all") {
					return true
				}
				if strings.Contains(comment.Text, ruleName) {
					return true
				}
			}
		}
	}
	return false
}
