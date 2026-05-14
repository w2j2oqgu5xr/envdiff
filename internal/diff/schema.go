package diff

import "sort"

// SchemaResult holds the result of comparing a key's presence across envs
// against an expected schema (set of required keys).
type SchemaResult struct {
	Key     string
	Present map[string]bool // env name -> whether key exists
}

// ValidateSchema checks each env against the provided schema keys and reports
// which keys are present or absent per environment.
func ValidateSchema(envs map[string]map[string]string, schemaKeys []string) []SchemaResult {
	if len(envs) == 0 || len(schemaKeys) == 0 {
		return nil
	}

	results := make([]SchemaResult, 0, len(schemaKeys))

	for _, key := range schemaKeys {
		presence := make(map[string]bool, len(envs))
		for envName, kv := range envs {
			_, exists := kv[key]
			presence[envName] = exists
		}
		results = append(results, SchemaResult{
			Key:     key,
			Present: presence,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}

// SchemaFullyCovered returns true if every env has every schema key.
func SchemaFullyCovered(results []SchemaResult) bool {
	for _, r := range results {
		for _, present := range r.Present {
			if !present {
				return false
			}
		}
	}
	return true
}
