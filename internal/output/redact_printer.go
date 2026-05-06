package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// RedactPrinter prints a report of sensitive keys found across envs.
type RedactPrinter struct {
	w io.Writer
}

// NewRedactPrinter creates a RedactPrinter writing to w (defaults to os.Stdout).
func NewRedactPrinter(w io.Writer) *RedactPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &RedactPrinter{w: w}
}

// Print outputs the redaction report for the given results.
func (p *RedactPrinter) Print(results []diff.RedactResult) {
	if len(results) == 0 {
		fmt.Fprintln(p.w, cyan("No sensitive keys detected."))
		return
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].EnvName != results[j].EnvName {
			return results[i].EnvName < results[j].EnvName
		}
		return results[i].Key < results[j].Key
	})

	fmt.Fprintln(p.w, bold("Sensitive Keys Detected:"))
	fmt.Fprintln(p.w, "")

	currentEnv := ""
	for _, r := range results {
		if r.EnvName != currentEnv {
			currentEnv = r.EnvName
			fmt.Fprintf(p.w, "  %s\n", bold(cyan(r.EnvName)))
		}
		fmt.Fprintf(p.w, "    %s = %s\n", yellow(r.Key), red(r.Redacted))
	}
	fmt.Fprintln(p.w, "")
	fmt.Fprintf(p.w, "  %s sensitive key(s) found across %s env(s).\n",
		bold(fmt.Sprintf("%d", len(results))),
		bold(fmt.Sprintf("%d", countUniqueEnvs(results))),
	)
}

func countUniqueEnvs(results []diff.RedactResult) int {
	seen := make(map[string]struct{})
	for _, r := range results {
		seen[r.EnvName] = struct{}{}
	}
	return len(seen)
}
