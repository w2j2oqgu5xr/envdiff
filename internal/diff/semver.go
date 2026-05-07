package diff

import (
	"fmt"
	"regexp"
	"sort"
)

// semverPattern matches basic semver strings like v1.2.3 or 1.2.3
var semverPattern = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(-[\w.]+)?$`)

// SemverMismatch represents a key whose value looks like a version string
// but differs across environments.
type SemverMismatch struct {
	Key    string
	Values map[string]string // env name -> value
}

// DetectSemverMismatches scans all envs for keys whose values look like
// semantic version strings and reports those that differ across environments.
func DetectSemverMismatches(envs map[string]map[string]string) []SemverMismatch {
	if len(envs) == 0 {
		return nil
	}

	// Collect all keys across all envs
	keySet := map[string]struct{}{}
	for _, kv := range envs {
		for k := range kv {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var results []SemverMismatch

	for _, key := range keys {
		values := map[string]string{}
		for envName, kv := range envs {
			if v, ok := kv[key]; ok && semverPattern.MatchString(v) {
				values[envName] = v
			}
		}

		if len(values) < 2 {
			continue
		}

		if hasSemverDiversity(values) {
			results = append(results, SemverMismatch{
				Key:    key,
				Values: values,
			})
		}
	}

	return results
}

func hasSemverDiversity(values map[string]string) bool {
	var first string
	for _, v := range values {
		if first == "" {
			first = v
			continue
		}
		if v != first {
			return true
		}
	}
		return false
}

// FormatSemverKey returns a display label for a semver mismatch key.
func FormatSemverKey(m SemverMismatch) string {
	return fmt.Sprintf("[semver] %s", m.Key)
}
