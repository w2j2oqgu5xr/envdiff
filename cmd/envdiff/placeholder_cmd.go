package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

// runPlaceholder loads the provided env files and reports any keys whose
// values look like unfilled placeholders (e.g. TODO, CHANGEME, <your-value>).
func runPlaceholder(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envdiff placeholder <name=file> [name=file ...]")
		os.Exit(1)
	}

	envs := make(map[string]map[string]string, len(args))
	for _, arg := range args {
		name, path, err := splitNamePath(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid argument %q: %v\n", arg, err)
			os.Exit(1)
		}
		kv, err := parser.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			os.Exit(1)
		}
		envs[name] = kv
	}

	results := diff.DetectPlaceholders(envs)
	p := output.NewPlaceholderPrinter()
	p.Print(results)

	if len(results) > 0 {
		os.Exit(1)
	}
}
