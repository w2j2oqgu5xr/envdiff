package diff

import (
	"fmt"
	"sort"
)

// BaselineDiff holds the result of comparing all envs against a single baseline env.
type BaselineDiff struct {
	Baseline string
	Env      string
	Missing  []string // keys in baseline but not in env
	Extra    []string // keys in env but not in baseline
	Changed  []string // keys present in both but with different values
}

// CompareToBaseline compares each non-baseline environment against the named
// baseline environment. It returns one BaselineDiff per compared env.
func CompareToBaseline(envs map[string]map[string]string, baseline string) ([]BaselineDiff, error) {
	baseEnv, ok := envs[baseline]
	if !ok {
		return nil, fmt.Errorf("baseline env %q not found", baseline)
	}

	var results []BaselineDiff

	names := make([]string, 0, len(envs))
	for name := range envs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if name == baseline {
			continue
		}
		env := envs[name]
		bd := BaselineDiff{
			Baseline: baseline,
			Env:      name,
		}
		for key, bVal := range baseEnv {
			if eVal, exists := env[key]; !exists {
				bd.Missing = append(bd.Missing, key)
			} else if eVal != bVal {
				bd.Changed = append(bd.Changed, key)
			}
		}
		for key := range env {
			if _, exists := baseEnv[key]; !exists {
				bd.Extra = append(bd.Extra, key)
			}
		}
		sort.Strings(bd.Missing)
		sort.Strings(bd.Extra)
		sort.Strings(bd.Changed)
		results = append(results, bd)
	}
	return results, nil
}
