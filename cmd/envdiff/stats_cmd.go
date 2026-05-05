package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your/envdiff/internal/diff"
	"github.com/your/envdiff/internal/output"
	"github.com/your/envdiff/internal/parser"
)

// runStats parses the provided env files and prints per-environment
// coverage statistics. It is invoked when the -stats flag is set.
func runStats(args []string) {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: envdiff -stats name=file [name=file ...]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pairs := fs.Args()
	if len(pairs) < 2 {
		fmt.Fprintln(os.Stderr, "error: at least two name=file pairs are required for stats")
		fs.Usage()
		os.Exit(1)
	}

	envs := map[string]map[string]string{}
	for _, pair := range pairs {
		name, path, err := splitNamePath(pair)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		kv, err := parser.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			os.Exit(1)
		}
		envs[name] = kv
	}

	results := diff.Compare(envs)
	stats := diff.ComputeEnvStats(results)

	p := output.NewStatsPrinter(os.Stdout)
	p.Print(stats)
}
