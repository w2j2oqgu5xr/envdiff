package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// TableWriter renders diff results as an aligned ASCII table.
type TableWriter struct {
	w io.Writer
}

// NewTableWriter returns a TableWriter that writes to w.
func NewTableWriter(w io.Writer) *TableWriter {
	return &TableWriter{w: w}
}

// Write renders results as a formatted table.
func (t *TableWriter) Write(results []diff.Result) {
	if len(results) == 0 {
		fmt.Fprintln(t.w, "(no results)")
		return
	}

	envNames := collectEnvNames(results)

	// Compute column widths.
	keyWidth := len("KEY")
	statusWidth := len("STATUS")
	for _, r := range results {
		if len(r.Key) > keyWidth {
			keyWidth = len(r.Key)
		}
		if len(string(r.Status)) > statusWidth {
			statusWidth = len(string(r.Status))
		}
	}

	envWidth := 12
	for _, e := range envNames {
		if len(e) > envWidth {
			envWidth = len(e)
		}
	}

	// Header.
	header := fmt.Sprintf("%-*s  %-*s", keyWidth, "KEY", statusWidth, "STATUS")
	for _, e := range envNames {
		header += fmt.Sprintf("  %-*s", envWidth, e)
	}
	fmt.Fprintln(t.w, bold(header))
	fmt.Fprintln(t.w, strings.Repeat("-", len(header)))

	// Rows.
	for _, r := range results {
		colorFn := statusColor(r.Status)
		row := fmt.Sprintf("%-*s  %-*s", keyWidth, r.Key, statusWidth, string(r.Status))
		for _, e := range envNames {
			v := r.Values[e]
			row += fmt.Sprintf("  %-*s", envWidth, v)
		}
		fmt.Fprintln(t.w, colorFn(row))
	}
}

func statusColor(s diff.Status) func(string) string {
	switch s {
	case diff.StatusMissing:
		return red
	case diff.StatusMismatch:
		return yellow
	case diff.StatusMatch:
		return green
	default:
		return func(v string) string { return v }
	}
}
