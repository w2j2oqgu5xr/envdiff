package diff

import "github.com/example/envdiff/internal/config"

// ApplyIgnore removes entries whose Key is listed in cfg.IgnoreKeys from a
// slice of Result values, returning a new filtered slice. If cfg.IgnoreKeys is
// empty the original slice is returned unchanged.
func ApplyIgnore(results []Result, cfg *config.Config) []Result {
	if len(cfg.IgnoreKeys) == 0 {
		return results
	}

	filtered := make([]Result, 0, len(results))
	for _, r := range results {
		if _, ignored := cfg.IgnoreKeys[r.Key]; !ignored {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
