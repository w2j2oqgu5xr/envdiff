package diff

import "sort"

// EnvCountResult holds the key count information for a single environment.
type EnvCountResult struct {
	EnvName  string
	KeyCount int
	UniqueKeys []string
}

// CountEnvKeys returns the number of keys defined in each environment,
// along with the sorted list of key names.
func CountEnvKeys(envs map[string]map[string]string) []EnvCountResult {
	results := make([]EnvCountResult, 0, len(envs))

	for envName, keys := range envs {
		unique := make([]string, 0, len(keys))
		for k := range keys {
			unique = append(unique, k)
		}
		sort.Strings(unique)

		results = append(results, EnvCountResult{
			EnvName:    envName,
			KeyCount:   len(keys),
			UniqueKeys: unique,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].EnvName < results[j].EnvName
	})

	return results
}

// KeyCountDelta returns pairs of environments and the delta (difference) in
// key count between them, sorted by absolute delta descending.
type KeyCountDelta struct {
	EnvA      string
	EnvB      string
	Delta     int // EnvA.KeyCount - EnvB.KeyCount
	AbsDelta  int
}

func KeyCountDeltas(results []EnvCountResult) []KeyCountDelta {
	deltas := []KeyCountDelta{}
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			d := results[i].KeyCount - results[j].KeyCount
			abs := d
			if abs < 0 {
				abs = -abs
			}
			deltas = append(deltas, KeyCountDelta{
				EnvA:     results[i].EnvName,
				EnvB:     results[j].EnvName,
				Delta:    d,
				AbsDelta: abs,
			})
		}
	}
	sort.Slice(deltas, func(i, j int) bool {
		return deltas[i].AbsDelta > deltas[j].AbsDelta
	})
	return deltas
}
