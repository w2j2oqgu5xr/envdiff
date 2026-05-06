package diff

import (
	"sort"
)

// DuplicateKey represents a key that appears more than once within a single env file.
type DuplicateKey struct {
	Key    string
	Count  int
	Values []string
}

// DuplicateResult holds duplicate key findings for one environment.
type DuplicateResult struct {
	EnvName    string
	Duplicates []DuplicateKey
}

// DetectDuplicates scans each environment's raw lines for keys that appear
// more than once and returns one DuplicateResult per environment that has
// at least one duplicate.
func DetectDuplicates(envs map[string]map[string]string, rawLines map[string][]string) []DuplicateResult {
	var results []DuplicateResult

	envNames := make([]string, 0, len(rawLines))
	for name := range rawLines {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	for _, name := range envNames {
		lines := rawLines[name]
		keyCounts := make(map[string][]string)

		for _, line := range lines {
			key, val, ok := splitRawLine(line)
			if !ok {
				continue
			}
			keyCounts[key] = append(keyCounts[key], val)
		}

		var dupes []DuplicateKey
		for key, vals := range keyCounts {
			if len(vals) > 1 {
				dupes = append(dupes, DuplicateKey{
					Key:    key,
					Count:  len(vals),
					Values: vals,
				})
			}
		}

		if len(dupes) == 0 {
			continue
		}

		sort.Slice(dupes, func(i, j int) bool {
			return dupes[i].Key < dupes[j].Key
		})

		results = append(results, DuplicateResult{
			EnvName:    name,
			Duplicates: dupes,
		})
	}

	return results
}

// splitRawLine parses a raw "KEY=VALUE" line, returning (key, value, ok).
// Lines that are blank or start with '#' are skipped.
func splitRawLine(line string) (string, string, bool) {
	if len(line) == 0 || line[0] == '#' {
		return "", "", false
	}
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			return line[:i], line[i+1:], true
		}
	}
	return "", "", false
}
