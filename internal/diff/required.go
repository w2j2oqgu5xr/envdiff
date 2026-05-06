package diff

import "sort"

// RequiredResult holds the validation result for a single required key.
type RequiredResult struct {
	Key     string
	Missing []string // env names where this key is absent
	Present []string // env names where this key exists
}

// CheckRequired verifies that all keys in requiredKeys exist in every
// provided environment. It returns one RequiredResult per required key.
func CheckRequired(envs map[string]map[string]string, requiredKeys []string) []RequiredResult {
	if len(envs) == 0 || len(requiredKeys) == 0 {
		return nil
	}

	// Collect sorted env names for deterministic output.
	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	results := make([]RequiredResult, 0, len(requiredKeys))

	for _, key := range requiredKeys {
		res := RequiredResult{Key: key}
		for _, name := range envNames {
			if _, ok := envs[name][key]; ok {
				res.Present = append(res.Present, name)
			} else {
				res.Missing = append(res.Missing, name)
			}
		}
		results = append(results, res)
	}

	return results
}

// AnyMissing returns true if at least one RequiredResult has missing envs.
func AnyMissing(results []RequiredResult) bool {
	for _, r := range results {
		if len(r.Missing) > 0 {
			return true
		}
	}
	return false
}
