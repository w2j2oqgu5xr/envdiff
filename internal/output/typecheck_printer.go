package output

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// TypecheckPrinter renders type mismatch results to an output stream.
type TypecheckPrinter struct {
	w io.Writer
}

// NewTypecheckPrinter creates a TypecheckPrinter writing to w.
// If w is nil, os.Stdout is used.
func NewTypecheckPrinter(w io.Writer) *TypecheckPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &TypecheckPrinter{w: w}
}

// Print writes type mismatch results. If no mismatches are found, a success
// message is printed instead.
func (p *TypecheckPrinter) Print(mismatches []diff.TypeMismatch) {
	if len(mismatches) == 0 {
		fmt.Fprintln(p.w, green("✔ No type mismatches detected."))
		return
	}

	fmt.Fprintf(p.w, bold("Type Mismatches (%d):\n"), len(mismatches))
	fmt.Fprintln(p.w, strings.Repeat("─", 48))

	for _, m := range mismatches {
		fmt.Fprintf(p.w, "  %s %s\n", cyan("KEY:"), bold(m.Key))
		for _, envName := range sortedEnvNames(m.Types) {
			t := m.Types[envName]
			colored := colorizeType(t)
			fmt.Fprintf(p.w, "    %-20s %s\n", envName+":", colored)
		}
		fmt.Fprintln(p.w)
	}
}

func sortedEnvNames(types map[string]diff.ValueType) []string {
	names := make([]string, 0, len(types))
	for k := range types {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

func colorizeType(t diff.ValueType) string {
	switch t {
	case diff.TypeBool:
		return yellow(string(t))
	case diff.TypeInt, diff.TypeFloat:
		return cyan(string(t))
	case diff.TypeURL:
		return green(string(t))
	case diff.TypeEmpty:
		return red(string(t))
	default:
		return string(t)
	}
}
