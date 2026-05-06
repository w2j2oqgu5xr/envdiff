package diff

import (
	"regexp"
	"strings"
)

// RedactResult holds a key and its redacted value for a single env.
type RedactResult struct {
	Key     string
	EnvName string
	Original string
	Redacted string
}

var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|secret|token|api_?key|auth|credential|private_?key|access_?key)`),
}

// IsSensitive returns true if the key name looks like it holds a secret.
func IsSensitive(key string) bool {
	for _, re := range sensitivePatterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// RedactValue replaces most characters of a value with asterisks,
// preserving up to 2 leading characters for identification.
func RedactValue(value string) string {
	if len(value) == 0 {
		return ""
	}
	if len(value) <= 3 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-2)
}

// RedactEnvs scans all envs and returns redaction info for sensitive keys.
func RedactEnvs(envs map[string]map[string]string) []RedactResult {
	var results []RedactResult
	for envName, kv := range envs {
		for key, val := range kv {
			if IsSensitive(key) {
				results = append(results, RedactResult{
					Key:      key,
					EnvName:  envName,
					Original: val,
					Redacted: RedactValue(val),
				})
			}
		}
	}
	return results
}
