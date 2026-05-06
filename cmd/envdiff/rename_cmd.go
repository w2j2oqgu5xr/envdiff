package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

// runRename loads the given env files and checks for keys that match
// any old→new rename pair provided via --rename OLD:NEW flags.
//
// Usage: envdiff rename --rename OLD:NEW [--rename ...] name1=file1 name2=file2 ...
func runRename(args []string, renameFlags []string) error {
	if len(args) < 1 {
		return fmt.Errorf("rename: at least one env file required")
	}

	renames := make(diff.RenameMap)
	for _, flag := range renameFlags {
		parts := strings.SplitN(flag, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid rename flag %q: expected OLD:NEW format", flag)
		}
		renames[parts[0]] = parts[1]
	}

	if len(renames) == 0 {
		return fmt.Errorf("rename: at least one --rename OLD:NEW flag is required")
	}

	envs := make(map[string]map[string]string)
	for _, arg := range args {
		name, path, err := splitNamePath(arg)
		if err != nil {
			return fmt.Errorf("rename: %w", err)
		}
		parsed, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("rename: failed to parse %q: %w", path, err)
		}
		envs[name] = parsed
	}

	results := diff.DetectRenames(envs, renames)
	p := output.NewRenamePrinter(os.Stdout)
	p.Print(results)

	if len(results) > 0 {
		os.Exit(1)
	}
	return nil
}
