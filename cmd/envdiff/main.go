package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/config"
	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

func main() {
	cfg := config.DefaultConfig()

	flag.StringVar(&cfg.Format, "format", cfg.Format, "output format: text, csv, json")
	flag.BoolVar(&cfg.ShowMatches, "show-matches", cfg.ShowMatches, "include matching keys in output")
	flag.BoolVar(&cfg.ShowSummary, "summary", cfg.ShowSummary, "print a summary after the diff")
	var ignoreRaw string
	flag.StringVar(&ignoreRaw, "ignore", "", "comma-separated list of keys to ignore")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff [flags] name1=file1 name2=file2 ...")
		os.Exit(1)
	}

	if ignoreRaw != "" {
		for _, k := range strings.Split(ignoreRaw, ",") {
			cfg.IgnoreKeys = append(cfg.IgnoreKeys, strings.TrimSpace(k))
		}
	}

	for _, arg := range args {
		name, path, err := splitNamePath(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid argument %q: %v\n", arg, err)
			os.Exit(1)
		}
		cfg.Envs = append(cfg.Envs, config.EnvEntry{Name: name, Path: path})
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	envMaps := make(map[string]map[string]string, len(cfg.Envs))
	for _, e := range cfg.Envs {
		m, err := parser.ParseFile(e.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", e.Path, err)
			os.Exit(1)
		}
		envMaps[e.Name] = m
	}

	results := diff.Compare(envMaps)
	results = diff.ApplyIgnore(results, cfg.IgnoreKeys)

	formatted, err := output.FormatResults(results, cfg.Format, cfg.ShowMatches)
	if err != nil {
		fmt.Fprintf(os.Stderr, "format error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(formatted)

	if cfg.ShowSummary {
		summary := output.BuildSummary(results)
		output.PrintSummary(os.Stdout, summary)
	}
}

// splitNamePath parses a "name=path" argument into its name and file path
// components. Returns an error if the argument is not in the expected format
// or if either component is empty.
func splitNamePath(arg string) (string, string, error) {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected format name=path")
	}
	return parts[0], parts[1], nil
}
