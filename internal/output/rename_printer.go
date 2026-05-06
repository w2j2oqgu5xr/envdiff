package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// RenamePrinter prints detected rename suggestions to a writer.
type RenamePrinter struct {
	w io.Writer
}

// NewRenamePrinter creates a RenamePrinter writing to the given writer.
// If w is nil, os.Stdout is used.
func NewRenamePrinter(w io.Writer) *RenamePrinter {
	if w == nil {
		w = os.Stdout
	}
	return &RenamePrinter{w: w}
}

// Print outputs rename suggestions grouped by environment.
func (p *RenamePrinter) Print(results []diff.RenameResult) {
	if len(results) == 0 {
		fmt.Fprintln(p.w, green("✔ No rename issues detected."))
		return
	}

	// Group by env name.
	grouped := make(map[string][]diff.RenameResult)
	for _, r := range results {
		grouped[r.EnvName] = append(grouped[r.EnvName], r)
	}

	envNames := make([]string, 0, len(grouped))
	for name := range grouped {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	fmt.Fprintln(p.w, bold(yellow("⚠ Possible renamed keys detected:")))
	for _, envName := range envNames {
		fmt.Fprintf(p.w, "  %s\n", cyan(envName))
		for _, r := range grouped[envName] {
			fmt.Fprintf(p.w, "    %s → %s  (old key present; new key missing)\n",
				red(r.OldKey), green(r.NewKey))
		}
	}
}
