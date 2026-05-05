package output

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format controls how the output is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
)

// FormatResults serializes diff results into the requested format.
func FormatResults(results []diff.Result, format Format) (string, error) {
	switch format {
	case FormatCSV:
		return formatCSV(results), nil
	case FormatJSON:
		return formatJSON(results), nil
	case FormatText, "":
		return formatText(results), nil
	default:
		return "", fmt.Errorf("unknown format %q: expected text, csv, or json", format)
	}
}

func formatText(results []diff.Result) string {
	if len(results) == 0 {
		return "No differences found.\n"
	}
	var sb strings.Builder
	for _, r := range results {
		status := strings.ToUpper(string(r.Status))
		sb.WriteString(fmt.Sprintf("[%s] %s\n", status, r.Key))
		for env, val := range r.Values {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", env, val))
		}
	}
	return sb.String()
}

func formatCSV(results []diff.Result) string {
	var sb strings.Builder
	sb.WriteString("key,status,env,value\n")
	for _, r := range results {
		for env, val := range r.Values {
			sb.WriteString(fmt.Sprintf("%s,%s,%s,%s\n",
				csvEscape(r.Key),
				string(r.Status),
				csvEscape(env),
				csvEscape(val),
			))
		}
	}
	return sb.String()
}

func formatJSON(results []diff.Result) string {
	if len(results) == 0 {
		return "[]\n"
	}
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, r := range results {
		sb.WriteString(fmt.Sprintf("  {\"key\":%q,\"status\":%q,\"values\":{", r.Key, r.Status))
		j := 0
		for env, val := range r.Values {
			if j > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(fmt.Sprintf("%q:%q", env, val))
			j++
		}
		sb.WriteString("}}") 
		if i < len(results)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	return sb.String()
}

func csvEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}
