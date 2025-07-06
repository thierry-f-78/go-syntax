package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/thierry-f-78/go-syntax/pkg/linter"
	"github.com/thierry-f-78/go-syntax/pkg/types"
)

func main() {
	var l *linter.Linter
	var files []string
	var err error
	var issues []types.Issue

	var path = flag.String("path", ".", "Path to analyze")
	var verbose *bool = flag.Bool("v", false, "Verbose output")
	var exitCode = flag.Int("exit-code", 1, "Exit code when issues are found")
	var color *bool = flag.Bool("c", true, "Color output")

	var red = "\033[31m"
	var blue = "\033[34m"
	var reset = "\033[0m"

	flag.Parse()

	if !*color {
		red = ""
		blue = ""
		reset = ""
	}

	l = linter.New()

	err = filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor/") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		os.Exit(1)
	}

	issues = l.Lint(files)

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].File != issues[j].File {
			return issues[i].File < issues[j].File // File name alpha sort
		}
		return issues[i].Line > issues[j].Line // File line dec
	})

	for _, issue := range issues {
		fmt.Printf("%s%s:%d:%d: [%s] %s%s\n",
			red, issue.File, issue.Line, issue.Column,
			issue.Rule, issue.Message, reset,
		)
		if *verbose {
			fmt.Printf("  %s%s%s\n", blue, issue.Description, reset)
			fmt.Printf("\n")
		}
	}

	if *verbose {
		fmt.Printf("Analyzed %d files\n", len(files))
	}

	if len(issues) > 0 {
		os.Exit(*exitCode)
	}
}
