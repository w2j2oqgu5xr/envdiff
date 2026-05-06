package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// DeprecatePrinter prints deprecated-key scan results to a writer.
type DeprecatePrinter struct {
	w io.Writer
}

// NewDeprecatePrinter returns a DeprecatePrinter writing to w.
// If w is nil it falls back to os.Stdout.
func NewDeprecatePrinter(w io.Writer) *DeprecatePrinter {
	if w == nil {
		w = os.Stdout
	}
	return &DeprecatePrinter{w: w}
}

// Print writes the deprecation results to the configured writer.
func (p *DeprecatePrinter) Print(results []diff.DeprecatedResult) {
	if len(results) == 0 {
		fmt.Fprintln(p.w, green("✔ No deprecated keys detected."))
		return
	}

	fmt.Fprintln(p.w, bold(yellow(fmt.Sprintf("⚠  %d deprecated key(s) found:", len(results)))))
	fmt.Fprintln(p.w)

	for _, r := range results {
		fmt.Fprintf(p.w, "  %s\n", bold(red(r.Key)))
		fmt.Fprintf(p.w, "    %s %s\n", cyan("hint:"), r.Message)
		fmt.Fprintf(p.w, "    %s %s\n", cyan("found in:"), strings.Join(r.Envs, ", "))
		fmt.Fprintln(p.w)
	}
}
