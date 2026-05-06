package diff

import (
	"regexp"
	"strconv"
	"strings"
)

// ValueType represents the inferred type of an env value.
type ValueType string

const (
	TypeString  ValueType = "string"
	TypeInt     ValueType = "int"
	TypeFloat   ValueType = "float"
	TypeBool    ValueType = "bool"
	TypeURL     ValueType = "url"
	TypeEmpty   ValueType = "empty"
)

var urlPattern = regexp.MustCompile(`^https?://`)

// InferType returns the inferred ValueType for a given string value.
func InferType(v string) ValueType {
	if v == "" {
		return TypeEmpty
	}
	lower := strings.ToLower(v)
	if lower == "true" || lower == "false" {
		return TypeBool
	}
	if urlPattern.MatchString(v) {
		return TypeURL
	}
	if _, err := strconv.Atoi(v); err == nil {
		return TypeInt
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return TypeFloat
	}
	return TypeString
}

// TypeMismatch describes a key whose value types differ across environments.
type TypeMismatch struct {
	Key   string
	Types map[string]ValueType // env name -> inferred type
}

// DetectTypeMismatches scans all envs and returns keys where the inferred
// value type is not consistent across environments.
func DetectTypeMismatches(envs map[string]map[string]string) []TypeMismatch {
	keys := collectAllKeys(envs)
	var results []TypeMismatch

	for _, key := range keys {
		types := make(map[string]ValueType)
		for envName, kv := range envs {
			if val, ok := kv[key]; ok {
				types[envName] = InferType(val)
			}
		}
		if hasTypeDiversity(types) {
			results = append(results, TypeMismatch{Key: key, Types: types})
		}
	}
	return results
}

func collectAllKeys(envs map[string]map[string]string) []string {
	seen := make(map[string]struct{})
	for _, kv := range envs {
		for k := range kv {
			seen[k] = struct{}{}
		}
	}
	return sortedKeys(seen)
}

func hasTypeDiversity(types map[string]ValueType) bool {
	if len(types) < 2 {
		return false
	}
	var first ValueType
	for _, t := range types {
		if first == "" {
			first = t
			continue
		}
		if t != first {
			return true
		}
	}
	return false
}
