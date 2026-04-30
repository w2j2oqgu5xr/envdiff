package diff

import "sort"

// KeyStatus represents the comparison status of a key across environments.
type KeyStatus int

const (
	StatusMatch    KeyStatus = iota // key exists in all envs with same value
	StatusMissing                   // key is absent in one or more envs
	StatusMismatch                  // key exists in all envs but values differ
)

// Entry holds the comparison result for a single key.
type Entry struct {
	Key    string
	Values map[string]string // env name -> value
	Status KeyStatus
}

// Result holds the full diff output for a set of environments.
type Result struct {
	EnvNames []string
	Entries  []Entry
}

// Compare takes a map of env name -> parsed key/value pairs and returns a Result.
func Compare(envs map[string]map[string]string) Result {
	allKeys := collectKeys(envs)
	envNames := sortedKeys(envs)

	entries := make([]Entry, 0, len(allKeys))
	for _, key := range allKeys {
		entry := Entry{
			Key:    key,
			Values: make(map[string]string, len(envNames)),
		}

		presentCount := 0
		uniqueVals := map[string]struct{}{}

		for _, name := range envNames {
			val, ok := envs[name][key]
			if ok {
				presentCount++
				uniqueVals[val] = struct{}{}
				entry.Values[name] = val
			} else {
				entry.Values[name] = ""
			}
		}

		switch {
		case presentCount < len(envNames):
			entry.Status = StatusMissing
		case len(uniqueVals) > 1:
			entry.Status = StatusMismatch
		default:
			entry.Status = StatusMatch
		}

		entries = append(entries, entry)
	}

	return Result{EnvNames: envNames, Entries: entries}
}

func collectKeys(envs map[string]map[string]string) []string {
	seen := map[string]struct{}{}
	for _, kv := range envs {
		for k := range kv {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
