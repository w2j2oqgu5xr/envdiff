package diff

import (
	"sort"
	"strings"
)

// PlaceholderResult holds information about a key whose value appears to be
// a placeholder (e.g. "TODO", "CHANGEME", "<your-value>") in one or more envs.
type PlaceholderResult struct {
	Key     string
	EnvName string
	Value   string
	Pattern string
}

// placeholderPatterns lists substrings (lowercased) that indicate a value is
// a placeholder and has not been filled in.
var placeholderPatterns = []string{
	"todo",
	"changeme",
	"change_me",
	"fixme",
	"your-",
	"your_",
	"<your",
	"<value",
	"<fill",
	"example",
	"placeholder",
	"replace_me",
	"replace-me",
	"xxxx",
}

// DetectPlaceholders scans all envs for keys whose values match known
// placeholder patterns and returns a slice of PlaceholderResult.
func DetectPlaceholders(envs map[string]map[string]string) []PlaceholderResult {
	var results []PlaceholderResult

	envNames := sortedKeys(envs)
	for _, envName := range envNames {
		kv := envs[envName]
		keys := sortedKeys(kv)
		for _, key := range keys {
			val := kv[key]
			if pattern, ok := isPlaceholder(val); ok {
				results = append(results, PlaceholderResult{
					Key:     key,
					EnvName: envName,
					Value:   val,
					Pattern: pattern,
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Key != results[j].Key {
			return results[i].Key < results[j].Key
		}
		return results[i].EnvName < results[j].EnvName
	})

	return results
}

// isPlaceholder returns the matched pattern and true if the value looks like
// an unfilled placeholder.
func isPlaceholder(val string) (string, bool) {
	if val == "" {
		return "", false
	}
	lower := strings.ToLower(val)
	for _, p := range placeholderPatterns {
		if strings.Contains(lower, p) {
			return p, true
		}
	}
	return "", false
}
