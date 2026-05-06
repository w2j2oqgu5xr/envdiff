package diff

import "sort"

// UnusedResult holds a key that appears in some envs but not all,
// suggesting it may be unused or forgotten in certain environments.
type UnusedResult struct {
	Key        string
	PresentIn  []string
	AbsentFrom []string
}

// DetectUnused finds keys that exist in at least one env but are missing
// from one or more others. A key must appear in fewer than all envs to
// be considered potentially unused.
func DetectUnused(envs map[string]map[string]string) []UnusedResult {
	if len(envs) < 2 {
		return nil
	}

	// Collect all keys across all envs
	keyEnvs := make(map[string][]string)
	for envName, kv := range envs {
		for k := range kv {
			keyEnvs[k] = append(keyEnvs[k], envName)
		}
	}

	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	var results []UnusedResult
	for key, present := range keyEnvs {
		if len(present) == len(envs) {
			continue // key is in all envs, not unused
		}

		presentSet := make(map[string]bool, len(present))
		for _, e := range present {
			presentSet[e] = true
		}

		var absent []string
		for _, name := range envNames {
			if !presentSet[name] {
				absent = append(absent, name)
			}
		}

		sort.Strings(present)
		results = append(results, UnusedResult{
			Key:        key,
			PresentIn:  present,
			AbsentFrom: absent,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}
