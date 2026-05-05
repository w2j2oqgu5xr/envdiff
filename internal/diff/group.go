package diff

import "sort"

// GroupByStatus organizes CompareResults into a map keyed by status string.
func GroupByStatus(results []CompareResult) map[string][]CompareResult {
	groups := map[string][]CompareResult{
		"match":    {},
		"missing":  {},
		"mismatch": {},
	}
	for _, r := range results {
		key := string(r.Status)
		if _, ok := groups[key]; ok {
			groups[key] = append(groups[key], r)
		}
	}
	return groups
}

// GroupByKey organizes CompareResults into a map keyed by the env variable key.
func GroupByKey(results []CompareResult) map[string][]CompareResult {
	groups := make(map[string][]CompareResult)
	for _, r := range results {
		groups[r.Key] = append(groups[r.Key], r)
	}
	return groups
}

// GroupSummary holds counts for each status group.
type GroupSummary struct {
	Total    int
	Match    int
	Missing  int
	Mismatch int
}

// SummarizeGroups returns a GroupSummary from a grouped result map.
func SummarizeGroups(groups map[string][]CompareResult) GroupSummary {
	s := GroupSummary{}
	s.Match = len(groups["match"])
	s.Missing = len(groups["missing"])
	s.Mismatch = len(groups["mismatch"])
	s.Total = s.Match + s.Missing + s.Mismatch
	return s
}

// SortedGroupKeys returns the status keys in a consistent order.
func SortedGroupKeys(groups map[string][]CompareResult) []string {
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
