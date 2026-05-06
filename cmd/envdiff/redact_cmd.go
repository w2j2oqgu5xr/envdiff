package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

// runRedact implements the `envdiff redact` subcommand.
// It scans all provided env files for sensitive keys and prints redacted output.
func runRedact(args []string) {
	fs := flag.NewFlagSet("redact", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: envdiff redact name=file [name=file ...]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Scans env files for sensitive keys (passwords, tokens, secrets)")
		fmt.Fprintln(os.Stderr, "and prints their redacted values.")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	pairs := fs.Args()
	if len(pairs) < 1 {
		fmt.Fprintln(os.Stderr, "error: at least one name=file argument is required")
		fs.Usage()
		os.Exit(1)
	}

	envs := make(map[string]map[string]string)
	for _, pair := range pairs {
		name, path, err := splitNamePath(pair)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		kv, err := parser.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %q: %v\n", path, err)
			os.Exit(1)
		}
		envs[name] = kv
	}

	results := diff.RedactEnvs(envs)
	p := output.NewRedactPrinter(os.Stdout)
	p.Print(results)
}
