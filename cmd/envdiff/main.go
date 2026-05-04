package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

func main() {
	showMatches := flag.Bool("show-matches", false, "also display keys that match across all envs")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <name=file> [name=file ...]\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n  envdiff dev=.env.dev prod=.env.prod\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "error: at least two env files required")
		flag.Usage()
		os.Exit(1)
	}

	envs := make(map[string]map[string]string, len(args))
	for _, arg := range args {
		name, path, ok := splitNamePath(arg)
		if !ok {
			fmt.Fprintf(os.Stderr, "error: invalid argument %q, expected name=path\n", arg)
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
	p := output.NewPrinter(output.Options{ShowMatches: *showMatches})
	p.Print(results)

	for _, r := range results {
		if r.Status != diff.StatusMatch {
			os.Exit(1)
		}
	}
}

func splitNamePath(arg string) (name, path string, ok bool) {
	for i, c := range arg {
		if c == '=' {
			return arg[:i], arg[i+1:], true
		}
	}
	return "", "", false
}
