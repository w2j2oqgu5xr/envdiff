package diff

import "sort"

// OverlapResult holds the overlap analysis between two environments.
type OverlapResult struct {
	EnvA       string
	EnvB       string
	SharedKeys []string
	OnlyInA    []string
	OnlyInB    []string
	OverlapPct float64 // percentage of keys shared out of the union
}

// ComputeOverlap calculates the key overlap between every pair of environments.
func ComputeOverlap(envs map[string]map[string]string) []OverlapResult {
	names := make([]string, 0, len(envs))
	for name := range envs {
		names = append(names, name)
	}
	sort.Strings(names)

	var results []OverlapResult
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			a, b := names[i], names[j]
			results = append(results, computePair(a, b, envs[a], envs[b]))
		}
	}
	return results
}

func computePair(nameA, nameB string, a, b map[string]string) OverlapResult {
	setA := keySet(a)
	setB := keySet(b)

	var shared, onlyA, onlyB []string

	for k := range setA {
		if setB[k] {
			shared = append(shared, k)
		} else {
			onlyA = append(onlyA, k)
		}
	}
	for k := range setB {
		if !setA[k] {
			onlyB = append(onlyB, k)
		}
	}

	sort.Strings(shared)
	sort.Strings(onlyA)
	sort.Strings(onlyB)

	unionSize := len(shared) + len(onlyA) + len(onlyB)
	var pct float64
	if unionSize > 0 {
		pct = float64(len(shared)) / float64(unionSize) * 100
	}

	return OverlapResult{
		EnvA:       nameA,
		EnvB:       nameB,
		SharedKeys: shared,
		OnlyInA:    onlyA,
		OnlyInB:    onlyB,
		OverlapPct: pct,
	}
}

func keySet(m map[string]string) map[string]bool {
	s := make(map[string]bool, len(m))
	for k := range m {
		s[k] = true
	}
	return s
}
