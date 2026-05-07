package output

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/user/envdiff/internal/diff"
)

// PlaceholderPrinter renders placeholder detection results to stdout.
type PlaceholderPrinter struct {
	w io.Writer
}

// NewPlaceholderPrinter creates a PlaceholderPrinter writing to stdout.
func NewPlaceholderPrinter() *PlaceholderPrinter {
	return &PlaceholderPrinter{w: os.Stdout}
}

// Print writes the placeholder results in a human-readable table.
func (p *PlaceholderPrinter) Print(results []diff.PlaceholderResult) {
	if len(results) == 0 {
		fmt.Fprintln(p.w, green("✔ No placeholder values detected."))
		return
	}

	fmt.Fprintf(p.w, bold("Placeholder Values Detected (%d)\n"), len(results))
	fmt.Fprintln(p.w)

	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, bold("KEY")+'\t'+bold("ENV")+'\t'+bold("VALUE")+'\t'+bold("PATTERN"))
	fmt.Fprintln(tw, "---\t---\t-----\t-------")

	for _, r := range results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			yellow(r.Key),
			cyan(r.EnvName),
			red(r.Value),
			r.Pattern,
		)
	}
	tw.Flush()
}
