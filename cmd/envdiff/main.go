package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

func main() {
	showMatches := flag.Bool("show-matches", false, "also display keys that match across all environments")
	format := flag.String("format", "text", "output format: text, csv, json")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] name1=file1 name2=file2 ...\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	envs := make(map[string]map[string]string, len(args))
	for _, arg := range args {
		name, path, err := splitNamePath(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		parsed, err := parser.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			os.Exit(1)
		}
		envs[name] = parsed
	}

	results := diff.Compare(envs)

	if output.Format(*format) == output.FormatText {
		p := output.NewPrinter(os.Stdout, *showMatches)
		p.Print(results)
		return
	}

	var filtered []diff.Result
	for _, r := range results {
		if *showMatches || r.Status != diff.StatusMatch {
			filtered = append(filtered, r)
		}
	}

	out, err := output.FormatResults(filtered, output.Format(*format))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(out)
}

func splitNamePath(arg string) (string, string, error) {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid argument %q: expected name=path", arg)
	}
	return parts[0], parts[1], nil
}
