package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// OverlapPrinter prints pairwise key-overlap reports.
type OverlapPrinter struct {
	w io.Writer
}

// NewOverlapPrinter creates an OverlapPrinter writing to w.
// If w is nil it falls back to os.Stdout.
func NewOverlapPrinter(w io.Writer) *OverlapPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &OverlapPrinter{w: w}
}

// Print renders the overlap results for all env pairs.
func (p *OverlapPrinter) Print(results []diff.OverlapResult) {
	if len(results) == 0 {
		fmt.Fprintln(p.w, cyan("No environment pairs to compare."))
		return
	}
	for _, r := range results {
		p.printPair(r)
	}
}

func (p *OverlapPrinter) printPair(r diff.OverlapResult) {
	header := fmt.Sprintf("── %s  ↔  %s", bold(r.EnvA), bold(r.EnvB))
	fmt.Fprintln(p.w, header)

	pctColor := green
	switch {
	case r.OverlapPct < 50:
		pctColor = red
	case r.OverlapPct < 80:
		pctColor = yellow
	}
	fmt.Fprintf(p.w, "   Overlap: %s\n", pctColor(fmt.Sprintf("%.1f%%", r.OverlapPct)))

	if len(r.SharedKeys) > 0 {
		fmt.Fprintf(p.w, "   Shared (%d): %s\n",
			len(r.SharedKeys), cyan(strings.Join(r.SharedKeys, ", ")))
	}
	if len(r.OnlyInA) > 0 {
		fmt.Fprintf(p.w, "   Only in %s (%d): %s\n",
			r.EnvA, len(r.OnlyInA), yellow(strings.Join(r.OnlyInA, ", ")))
	}
	if len(r.OnlyInB) > 0 {
		fmt.Fprintf(p.w, "   Only in %s (%d): %s\n",
			r.EnvB, len(r.OnlyInB), yellow(strings.Join(r.OnlyInB, ", ")))
	}
	fmt.Fprintln(p.w)
}
