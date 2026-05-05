package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// FormatResults serializes results into the requested format string.
// Supported formats: "text", "csv", "json", "table".
func FormatResults(results []diff.Result, format string) (string, error) {
	switch strings.ToLower(format) {
	case "text", "":
		return formatText(results), nil
	case "csv":
		return formatCSV(results), nil
	case "json":
		return formatJSON(results)
	case "table":
		var sb strings.Builder
		NewTableWriter(&sb).Write(results)
		return sb.String(), nil
	default:
		return "", fmt.Errorf("unknown format %q: must be text, csv, json, or table", format)
	}
}

func formatText(results []diff.Result) string {
	if len(results) == 0 {
		return ""
	}
	envNames := collectEnvNames(results)
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("[%s] %s", r.Status, r.Key))
		for _, e := range envNames {
			v := r.Values[e]
			if v == "" {
				v = "<missing>"
			}
			sb.WriteString(fmt.Sprintf(" %s=%s", e, v))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatCSV(results []diff.Result) string {
	if len(results) == 0 {
		return ""
	}
	envNames := collectEnvNames(results)
	var sb strings.Builder
	header := append([]string{"key", "status"}, envNames...)
	for i, h := range header {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(csvEscape(h))
	}
	sb.WriteString("\n")
	for _, r := range results {
		sb.WriteString(csvEscape(r.Key))
		sb.WriteString(",")
		sb.WriteString(csvEscape(string(r.Status)))
		for _, e := range envNames {
			sb.WriteString(",")
			sb.WriteString(csvEscape(r.Values[e]))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatJSON(results []diff.Result) (string, error) {
	type row struct {
		Key    string            `json:"key"`
		Status string            `json:"status"`
		Values map[string]string `json:"values"`
	}
	rows := make([]row, len(results))
	for i, r := range results {
		rows[i] = row{Key: r.Key, Status: string(r.Status), Values: r.Values}
	}
	b, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func csvEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}
