package diff

import (
	"sort"

	"github.com/user/envdiff/internal/config"
)

// SortOrder defines how diff results should be ordered.
type SortOrder string

const (
	SortByKey      SortOrder = "key"
	SortByStatus   SortOrder = "status"
	SortByEnvCount SortOrder = "envcount"
)

// SortResults returns a new slice of Results sorted by the given order.
func SortResults(results []config.Result, order SortOrder) []config.Result {
	copied := make([]config.Result, len(results))
	copy(copied, results)

	switch order {
	case SortByStatus:
		sort.SliceStable(copied, func(i, j int) bool {
			si := statusRank(copied[i].Status)
			sj := statusRank(copied[j].Status)
			if si != sj {
				return si < sj
			}
			return copied[i].Key < copied[j].Key
		})
	case SortByEnvCount:
		sort.SliceStable(copied, func(i, j int) bool {
			ci := len(copied[i].Values)
			cj := len(copied[j].Values)
			if ci != cj {
				return ci < cj
			}
			return copied[i].Key < copied[j].Key
		})
	default: // SortByKey
		sort.SliceStable(copied, func(i, j int) bool {
			return copied[i].Key < copied[j].Key
		})
	}

	return copied
}

func statusRank(s config.Status) int {
	switch s {
	case config.StatusMissing:
		return 0
	case config.StatusMismatch:
		return 1
	case config.StatusMatch:
		return 2
	default:
		return 3
	}
}
