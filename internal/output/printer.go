package output

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Options controls printer behavior.
type Options struct {
	ShowMatches bool
	NoColor     bool
}

// Printer renders diff results to stdout.
type Printer struct {
	opts Options
}

// NewPrinter creates a Printer with the given options.
func NewPrinter(opts Options) *Printer {
	return &Printer{opts: opts}
}

// Print renders all diff results with colorized output.
func (p *Printer) Print(results []diff.Result) {
	if len(results) == 0 {
		fmt.Println(green("All environments are identical."))
		return
	}

	envNames := collectEnvNames(results)
	printHeader(envNames)

	missing, mismatched, matched := 0, 0, 0

	for _, r := range results {
		switch r.Status {
		case diff.StatusMissing:
			printMissing(r.MissingIn, r.Key)
			missing++
		case diff.StatusMismatch:
			printMismatch(r.Key, r.Values)
			mismatched++
		case diff.StatusMatch:
			if p.opts.ShowMatches {
				printMatch(r.Key)
			}
			matched++
		}
	}

	p.printSummary(missing, mismatched, matched)
}

func (p *Printer) printSummary(missing, mismatched, matched int) {
	fmt.Println()
	fmt.Println(bold("Summary:"))
	fmt.Printf("  %s missing keys\n", red(fmt.Sprintf("%d", missing)))
	fmt.Printf("  %s mismatched values\n", yellow(fmt.Sprintf("%d", mismatched)))
	fmt.Printf("  %s matching keys\n", green(fmt.Sprintf("%d", matched)))
}

func collectEnvNames(results []diff.Result) []string {
	seen := map[string]struct{}{}
	for _, r := range results {
		for env := range r.Values {
			seen[env] = struct{}{}
		}
		if r.MissingIn != "" {
			seen[r.MissingIn] = struct{}{}
		}
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
