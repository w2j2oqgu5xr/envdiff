package diff

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// ApplyIgnore filters out diff results whose keys match any of the ignore patterns.
// Patterns are matched as case-insensitive prefix, suffix, or exact match.
func ApplyIgnore(results []diff.Result, ignoreKeys []string) []diff.Result {
	if len(ignoreKeys) == 0 {
		return results
	}

	ignoreSet := make(map[string]struct{}, len(ignoreKeys))
	for _, k := range ignoreKeys {
		ignoreSet[strings.ToUpper(k)] = struct{}{}
	}

	filtered := results[:0:0]
	for _, r := range results {
		if _, skip := ignoreSet[strings.ToUpper(r.Key)]; !skip {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
