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
			&rules.ConstNoTypeRule{},
			&rules.NamedReturnsRule{},
			&rules.NakedReturnRule{},
			&rules.IfInitRule{},
		},
	}
}

func (l *Linter) Lint(files []string) []types.Issue {
	var allIssues []types.Issue

	var file string
	for _, file = range files {
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

	var rule types.Rule
	for _, rule = range l.rules {
		var ruleIssues []types.Issue
		ruleIssues = rule.Check(fset, src)
		issues = append(issues, ruleIssues...)
	}

	return filterNolintIssues(issues, src, fset)
}

func filterNolintIssues(issues []types.Issue, file *ast.File, fset *token.FileSet) []types.Issue {
	var filtered []types.Issue

	var issue types.Issue
	for _, issue = range issues {
		if !isNolintComment(file, issue.Line, issue.Rule, fset) && !isFileNolintComment(file, issue.Rule) {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

func isNolintComment(file *ast.File, line int, ruleName string, fset *token.FileSet) bool {
	var commentGroup *ast.CommentGroup
	for _, commentGroup = range file.Comments {
		var comment *ast.Comment
		for _, comment = range commentGroup.List {
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

func isFileNolintComment(file *ast.File, ruleName string) bool {
	// Check only the first few comments in the file (file header)
	if len(file.Comments) == 0 {
		return false
	}

	// Look at the first comment group(s) which should be at the top of the file
	var commentGroup *ast.CommentGroup
	for _, commentGroup = range file.Comments {
		// Only check comments that are likely to be file headers (first 10 lines)
		if commentGroup.Pos() > file.Package+10 {
			break
		}

		var comment *ast.Comment
		for _, comment = range commentGroup.List {
			if strings.Contains(comment.Text, "nolint") {
				// Check for "//nolint" (disable all rules for file)
				if comment.Text == "//nolint" || comment.Text == "// nolint" {
					return true
				}
				// Check for "//nolint:all" (disable all rules for file)
				if strings.Contains(comment.Text, "nolint:all") {
					return true
				}
				// Check for specific rule in "//nolint:rule1,rule2"
				if strings.Contains(comment.Text, "nolint:") && strings.Contains(comment.Text, ruleName) {
					return true
				}
			}
		}
	}
	return false
}
