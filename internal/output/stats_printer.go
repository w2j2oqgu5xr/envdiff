package output

import (
	"fmt"
	"io"
	"os"

	"github.com/your/envdiff/internal/diff"
)

// StatsPrinter renders per-environment statistics to a writer.
type StatsPrinter struct {
	w io.Writer
}

// NewStatsPrinter creates a StatsPrinter writing to w.
// If w is nil, os.Stdout is used.
func NewStatsPrinter(w io.Writer) *StatsPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &StatsPrinter{w: w}
}

// Print writes a formatted stats table for each environment.
func (sp *StatsPrinter) Print(stats []diff.EnvStats) {
	if len(stats) == 0 {
		fmt.Fprintln(sp.w, cyan("No environment statistics available."))
		return
	}

	fmt.Fprintln(sp.w, bold("Environment Coverage Report"))
	fmt.Fprintln(sp.w, bold("─────────────────────────────────────────────────"))
	fmt.Fprintf(sp.w, "%-20s %6s %7s %9s %7s %8s\n",
		bold("ENV"), bold("TOTAL"), bold("MATCH"), bold("MISMATCH"), bold("MISS"), bold("COVER%"))
	fmt.Fprintln(sp.w, "─────────────────────────────────────────────────")

	for _, s := range stats {
		coverStr := fmt.Sprintf("%.1f%%", s.Coverage)
		var coverColored string
		switch {
		case s.Coverage >= 90.0:
			coverColored = green(coverStr)
		case s.Coverage >= 70.0:
			coverColored = yellow(coverStr)
		default:
			coverColored = red(coverStr)
		}

		fmt.Fprintf(sp.w, "%-20s %6d %7d %9d %7d %8s\n",
			cyan(s.EnvName),
			s.Total,
			s.Matched,
			s.Mismatched,
			s.Missing,
			coverColored,
		)
	}
	fmt.Fprintln(sp.w, "─────────────────────────────────────────────────")
}
