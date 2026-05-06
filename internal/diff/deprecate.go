package diff

// DeprecatedResult holds information about a key that is flagged as deprecated
// in one or more environments.
type DeprecatedResult struct {
	Key     string
	Envs    []string
	Message string
}

// DetectDeprecated scans parsed environments for keys listed in the deprecation
// map and returns a result for each key found in at least one environment.
// The deprecationMap maps key names (case-insensitive) to a human-readable
// deprecation message (e.g. "use NEW_KEY instead").
func DetectDeprecated(envs map[string]map[string]string, deprecationMap map[string]string) []DeprecatedResult {
	if len(envs) == 0 || len(deprecationMap) == 0 {
		return nil
	}

	// Normalise deprecation map keys to uppercase.
	normalised := make(map[string]string, len(deprecationMap))
	for k, msg := range deprecationMap {
		normalised[strings.ToUpper(k)] = msg
	}

	// Collect which envs contain each deprecated key.
	found := make(map[string][]string) // upper key -> env names
	for envName, keys := range envs {
		for k := range keys {
			upper := strings.ToUpper(k)
			if _, ok := normalised[upper]; ok {
				found[upper] = append(found[upper], envName)
			}
		}
	}

	if len(found) == 0 {
		return nil
	}

	results := make([]DeprecatedResult, 0, len(found))
	for upperKey, envNames := range found {
		sort.Strings(envNames)
		results = append(results, DeprecatedResult{
			Key:     upperKey,
			Envs:    envNames,
			Message: normalised[upperKey],
		})
	}

	// Stable output order.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}
