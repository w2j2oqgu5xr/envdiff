package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
)

// BaselinePrinter prints baseline comparison results to a writer.
type BaselinePrinter struct {
	w io.Writer
}

// NewBaselinePrinter creates a BaselinePrinter writing to stdout by default.
func NewBaselinePrinter(w io.Writer) *BaselinePrinter {
	if w == nil {
		w = os.Stdout
	}
	return &BaselinePrinter{w: w}
}

// Print renders each BaselineDiff with colorized sections.
func (p *BaselinePrinter) Print(diffs []diff.BaselineDiff) {
	if len(diffs) == 0 {
		fmt.Fprintln(p.w, cyan("No environments to compare against baseline."))
		return
	}
	for _, bd := range diffs {
		header := fmt.Sprintf("── %s vs %s ──", bold(bd.Baseline), bold(bd.Env))
		fmt.Fprintln(p.w, cyan(header))

		if len(bd.Missing) == 0 && len(bd.Extra) == 0 && len(bd.Changed) == 0 {
			fmt.Fprintln(p.w, green("  ✔ fully in sync"))
			continue
		}

		for _, k := range bd.Missing {
			fmt.Fprintf(p.w, "  %s %s\n", red("MISSING"), k)
		}
		for _, k := range bd.Extra {
			fmt.Fprintf(p.w, "  %s  %s\n", yellow("EXTRA"), k)
		}
		for _, k := range bd.Changed {
			fmt.Fprintf(p.w, "  %s  %s\n", yellow("CHANGED"), k)
		}
	}
}
