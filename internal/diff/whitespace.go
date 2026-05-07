package diff

import (
	"strings"
)

// WhitespaceDiff represents a key whose value differs only in whitespace between envs.
type WhitespaceDiff struct {
	Key      string
	EnvName  string
	RawValue string
	Trimmed  string
}

// DetectWhitespaceDiffs scans all envs for keys where values are equal after
// trimming but differ in their raw form, indicating accidental whitespace.
func DetectWhitespaceDiffs(envs map[string]map[string]string) []WhitespaceDiff {
	if len(envs) == 0 {
		return nil
	}

	// Collect all keys across all envs.
	keys := map[string]struct{}{}
	for _, kv := range envs {
		for k := range kv {
			keys[k] = struct{}{}
		}
	}

	var results []WhitespaceDiff

	for key := range keys {
		// Gather raw values per env for this key.
		type entry struct {
			env string
			raw string
		}
		var entries []entry
		for envName, kv := range envs {
			if raw, ok := kv[key]; ok {
				entries = append(entries, entry{envName, raw})
			}
		}

		if len(entries) < 2 {
			continue
		}

		// Check if any raw value has leading/trailing whitespace.
		for _, e := range entries {
			trimmed := strings.TrimSpace(e.raw)
			if trimmed != e.raw {
				results = append(results, WhitespaceDiff{
					Key:      key,
					EnvName:  e.env,
					RawValue: e.raw,
					Trimmed:  trimmed,
				})
			}
		}
	}

	return results
}
