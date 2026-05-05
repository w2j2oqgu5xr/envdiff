package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// GroupedPrinter prints results organized by their status group.
type GroupedPrinter struct {
	w           io.Writer
	showMatches bool
}

// NewGroupedPrinter creates a GroupedPrinter writing to w.
func NewGroupedPrinter(w io.Writer, showMatches bool) *GroupedPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &GroupedPrinter{w: w, showMatches: showMatches}
}

// Print outputs results grouped by status with section headers.
func (p *GroupedPrinter) Print(results []diff.CompareResult) {
	groups := diff.GroupByStatus(results)
	summary := diff.SummarizeGroups(groups)

	if summary.Total == 0 {
		fmt.Fprintln(p.w, cyan("No results to display."))
		return
	}

	order := []string{"missing", "mismatch", "match"}
	for _, status := range order {
		group := groups[status]
		if len(group) == 0 {
			continue
		}
		if status == "match" && !p.showMatches {
			continue
		}
		p.printSection(status, group)
	}

	fmt.Fprintf(p.w, "\n%s  total=%d  match=%d  missing=%d  mismatch=%d\n",
		bold("Summary:"), summary.Total, summary.Match, summary.Missing, summary.Mismatch)
}

func (p *GroupedPrinter) printSection(status string, results []diff.CompareResult) {
	var header string
	switch status {
	case "missing":
		header = red("[MISSING]")
	case "mismatch":
		header = yellow("[MISMATCH]")
	case "match":
		header = green("[MATCH]")
	default:
		header = bold("[" + strings.ToUpper(status) + "]")
	}

	fmt.Fprintf(p.w, "\n%s (%d)\n", header, len(results))
	fmt.Fprintln(p.w, strings.Repeat("-", 40))

	for _, r := range results {
		fmt.Fprintf(p.w, "  %s\n", bold(r.Key))
		for env, val := range r.Values {
			if val == "" {
				fmt.Fprintf(p.w, "    %s: %s\n", cyan(env), red("<missing>"))
			} else {
				fmt.Fprintf(p.w, "    %s: %s\n", cyan(env), val)
			}
		}
	}
}
