package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
)

// Summary holds aggregated statistics from a diff comparison.
type Summary struct {
	TotalKeys   int
	MatchedKeys int
	MissingKeys int
	MismatchKeys int
}

// BuildSummary computes a Summary from a slice of diff results.
func BuildSummary(results []diff.Result) Summary {
	s := Summary{}
	for _, r := range results {
		s.TotalKeys++
		switch r.Status {
		case diff.StatusMatch:
			s.MatchedKeys++
		case diff.StatusMissing:
			s.MissingKeys++
		case diff.StatusMismatch:
			s.MismatchKeys++
		}
	}
	return s
}

// PrintSummary writes a colorized summary block to the given writer.
// If w is nil, os.Stdout is used.
func PrintSummary(w io.Writer, s Summary) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintln(w, bold("\n── Summary ─────────────────────"))
	fmt.Fprintf(w, "  Total keys   : %s\n", cyan(fmt.Sprintf("%d", s.TotalKeys)))
	fmt.Fprintf(w, "  Matched      : %s\n", green(fmt.Sprintf("%d", s.MatchedKeys)))
	fmt.Fprintf(w, "  Mismatched   : %s\n", yellow(fmt.Sprintf("%d", s.MismatchKeys)))
	fmt.Fprintf(w, "  Missing      : %s\n", red(fmt.Sprintf("%d", s.MissingKeys)))
	fmt.Fprintln(w, bold("─────────────────────────────────"))
}
