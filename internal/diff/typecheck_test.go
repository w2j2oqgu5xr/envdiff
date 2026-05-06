package diff

import (
	"testing"
)

func TestInferType_Values(t *testing.T) {
	cases := []struct {
		input    string
		expected ValueType
	}{
		{"", TypeEmpty},
		{"true", TypeBool},
		{"false", TypeBool},
		{"TRUE", TypeBool},
		{"42", TypeInt},
		{"-7", TypeInt},
		{"3.14", TypeFloat},
		{"http://example.com", TypeURL},
		{"https://api.example.com/v1", TypeURL},
		{"hello world", TypeString},
		{"some_value", TypeString},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := InferType(tc.input)
			if got != tc.expected {
				t.Errorf("InferType(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestDetectTypeMismatches_NoMismatch(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"PORT": "8080", "DEBUG": "true"},
		"prod": {"PORT": "443", "DEBUG": "false"},
	}
	results := DetectTypeMismatches(envs)
	if len(results) != 0 {
		t.Errorf("expected no mismatches, got %d", len(results))
	}
}

func TestDetectTypeMismatches_DetectsMismatch(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"TIMEOUT": "30"},
		"prod": {"TIMEOUT": "thirty"},
	}
	results := DetectTypeMismatches(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(results))
	}
	if results[0].Key != "TIMEOUT" {
		t.Errorf("expected key TIMEOUT, got %s", results[0].Key)
	}
	if results[0].Types["dev"] != TypeInt {
		t.Errorf("expected dev=int, got %s", results[0].Types["dev"])
	}
	if results[0].Types["prod"] != TypeString {
		t.Errorf("expected prod=string, got %s", results[0].Types["prod"])
	}
}

func TestDetectTypeMismatches_EmptyEnvs(t *testing.T) {
	results := DetectTypeMismatches(map[string]map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected empty results for empty input")
	}
}

func TestDetectTypeMismatches_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"PORT": "8080", "URL": "http://localhost"},
	}
	results := DetectTypeMismatches(envs)
	if len(results) != 0 {
		t.Errorf("single env should never produce mismatches")
	}
}

func TestDetectTypeMismatches_MultipleMismatches(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"A": "123", "B": "true"},
		"prod": {"A": "abc", "B": "yes"},
	}
	results := DetectTypeMismatches(envs)
	if len(results) != 2 {
		t.Errorf("expected 2 mismatches, got %d", len(results))
	}
}
