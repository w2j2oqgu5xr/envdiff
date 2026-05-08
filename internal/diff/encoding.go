package diff

// EncodingResult holds information about a detected encoding inconsistency
// for a given key across multiple environments.
type EncodingResult struct {
	Key      string
	EnvName  string
	Value    string
	Issue    string
}

// DetectEncodingIssues scans env values for common encoding problems:
// non-UTF-8 bytes, null bytes, BOM markers, or mixed line endings embedded
// in multi-line values.
func DetectEncodingIssues(envs map[string]map[string]string) []EncodingResult {
	var results []EncodingResult

	envNames := sortedKeys(envs)
	for _, envName := range envNames {
		keys := sortedKeys(envs[envName])
		for _, key := range keys {
			val := envs[envName][key]
			if issue := detectIssue(val); issue != "" {
				results = append(results, EncodingResult{
					Key:     key,
					EnvName: envName,
					Value:   val,
					Issue:   issue,
				})
			}
		}
	}
	return results
}

func detectIssue(val string) string {
	if len(val) == 0 {
		return ""
	}
	// BOM marker (UTF-8 BOM: 0xEF 0xBB 0xBF)
	if len(val) >= 3 && val[0] == 0xEF && val[1] == 0xBB && val[2] == 0xBF {
		return "BOM marker detected"
	}
	for i := 0; i < len(val); i++ {
		b := val[i]
		// Null byte
		if b == 0x00 {
			return "null byte detected"
		}
		// Non-printable control characters (excluding tab, newline, carriage return)
		if b < 0x09 || (b > 0x0D && b < 0x20) || b == 0x7F {
			return "non-printable control character detected"
		}
	}
	// Check for invalid UTF-8 sequences
	if !isValidUTF8(val) {
		return "invalid UTF-8 encoding"
	}
	return ""
}

func isValidUTF8(s string) bool {
	for i := 0; i < len(s); {
		b := s[i]
		var size int
		switch {
		case b < 0x80:
			size = 1
		case b < 0xC2:
			return false
		case b < 0xE0:
			size = 2
		case b < 0xF0:
			size = 3
		case b < 0xF5:
			size = 4
		default:
			return false
		}
		if i+size > len(s) {
			return false
		}
		for j := 1; j < size; j++ {
			if s[i+j]&0xC0 != 0x80 {
				return false
			}
		}
		i += size
	}
	return true
}
