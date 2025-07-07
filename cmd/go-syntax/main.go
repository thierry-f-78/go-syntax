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

// stringSlice implements flag.Value for multiple string flags
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// isExcluded checks if a file path matches any of the exclude patterns
func isExcluded(filePath string, excludePatterns []string) bool {
	var pattern string
	for _, pattern = range excludePatterns {
		var matched bool
		var err error
		matched, err = filepath.Match(pattern, filepath.Base(filePath))
		if err == nil && matched {
			return true
		}

		// Also try matching the full path
		matched, err = filepath.Match(pattern, filePath)
		if err == nil && matched {
			return true
		}

		// Check if pattern matches any part of the path
		if strings.Contains(filePath, pattern) {
			return true
		}
	}
	return false
}

func main() {
	var l *linter.Linter
	var files []string
	var err error
	var issues []types.Issue

	var verbose *bool = flag.Bool("v", false, "Verbose output")
	var exitCode *int = flag.Int("exit-code", 1, "Exit code when issues are found")
	var color *bool = flag.Bool("c", true, "Color output")

	var excludePatterns stringSlice
	flag.Var(&excludePatterns, "e", "Exclude files matching pattern (can be repeated)")

	var red string = "\033[31m"
	var blue string = "\033[34m"
	var reset string = "\033[0m"

	flag.Parse()

	if !*color {
		red = ""
		blue = ""
		reset = ""
	}

	// Use command line arguments as paths, default to "." if none provided
	var paths []string = flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	l = linter.New()

	// Process each path argument
	var path string
	for _, path = range paths {
		var walkPath string
		var recursive bool

		// Handle Go-style path patterns
		if strings.HasSuffix(path, "/...") {
			walkPath = strings.TrimSuffix(path, "/...")
			recursive = true
		} else if path == "./..." {
			walkPath = "."
			recursive = true
		} else {
			walkPath = path
			recursive = false
		}

		err = filepath.Walk(walkPath, func(currentPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// If not recursive, only process files in the exact directory
			if !recursive {
				var rel string
				rel, _ = filepath.Rel(walkPath, currentPath)
				if strings.Contains(rel, string(filepath.Separator)) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}

			if strings.HasSuffix(currentPath, ".go") && !strings.Contains(currentPath, "vendor/") {
				if !isExcluded(currentPath, excludePatterns) {
					files = append(files, currentPath)
				}
			}
			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking directory %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	issues = l.Lint(files)

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].File != issues[j].File {
			return issues[i].File < issues[j].File // File name alpha sort
		}
		return issues[i].Line > issues[j].Line // File line dec
	})

	var issue types.Issue
	for _, issue = range issues {
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
