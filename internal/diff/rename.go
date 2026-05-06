package diff

import "sort"

// RenameMap maps old key names to new key names.
type RenameMap map[string]string

// RenameResult describes the outcome of a single key rename check.
type RenameResult struct {
	OldKey  string
	NewKey  string
	EnvName string
	Found   bool
}

// DetectRenames checks each environment for keys that appear under an old name
// but are missing under the new name, indicating a likely rename issue.
func DetectRenames(envs map[string]map[string]string, renames RenameMap) []RenameResult {
	var results []RenameResult

	// Collect env names in sorted order for deterministic output.
	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	for _, envName := range envNames {
		keys := envs[envName]
		for oldKey, newKey := range renames {
			_, hasOld := keys[oldKey]
			_, hasNew := keys[newKey]
			if hasOld && !hasNew {
				results = append(results, RenameResult{
					OldKey:  oldKey,
					NewKey:  newKey,
					EnvName: envName,
					Found:   true,
				})
			}
		}
	}

	// Sort results for deterministic output.
	sort.Slice(results, func(i, j int) bool {
		if results[i].EnvName != results[j].EnvName {
			return results[i].EnvName < results[j].EnvName
		}
		return results[i].OldKey < results[j].OldKey
	})

	return results
}
