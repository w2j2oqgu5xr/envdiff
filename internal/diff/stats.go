package diff

import "github.com/your/envdiff/internal/diff"

// EnvStats holds per-environment statistics derived from comparison results.
type EnvStats struct {
	EnvName    string
	Total      int
	Missing    int
	Mismatched int
	Matched    int
	Coverage   float64 // percentage of keys present (not missing)
}

// ComputeEnvStats calculates per-environment statistics from a slice of Results.
func ComputeEnvStats(results []Result) []EnvStats {
	type counts struct {
		total      int
		missing    int
		mismatched int
		matched    int
	}

	envMap := map[string]*counts{}

	for _, r := range results {
		for env, status := range r.Envs {
			if _, ok := envMap[env]; !ok {
				envMap[env] = &counts{}
			}
			c := envMap[env]
			c.total++
			switch status {
			case StatusMissing:
				c.missing++
			case StatusMismatch:
				c.mismatched++
			case StatusMatch:
				c.matched++
			}
		}
	}

	envNames := sortedKeys(envMap)
	stats := make([]EnvStats, 0, len(envNames))
	for _, name := range envNames {
		c := envMap[name]
		coverage := 0.0
		if c.total > 0 {
			coverage = float64(c.total-c.missing) / float64(c.total) * 100.0
		}
		stats = append(stats, EnvStats{
			EnvName:    name,
			Total:      c.total,
			Missing:    c.missing,
			Mismatched: c.mismatched,
			Matched:    c.matched,
			Coverage:   coverage,
		})
	}
	return stats
}
