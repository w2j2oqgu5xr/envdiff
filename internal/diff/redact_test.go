package diff

import (
	"testing"
)

func TestIsSensitive_MatchesSecretKeys(t *testing.T) {
	sensitive := []string{
		"PASSWORD", "db_password", "API_KEY", "api_key",
		"SECRET", "AUTH_TOKEN", "private_key", "ACCESS_KEY",
		"CREDENTIAL", "token",
	}
	for _, k := range sensitive {
		if !IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_IgnoresSafeKeys(t *testing.T) {
	safe := []string{"HOST", "PORT", "DEBUG", "APP_NAME", "LOG_LEVEL"}
	for _, k := range safe {
		if IsSensitive(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestRedactValue_ShortValue(t *testing.T) {
	if got := RedactValue("abc"); got != "***" {
		t.Errorf("got %q, want \"***\"", got)
	}
}

func TestRedactValue_Empty(t *testing.T) {
	if got := RedactValue(""); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestRedactValue_LongValue(t *testing.T) {
	got := RedactValue("supersecret")
	if got[:2] != "su" {
		t.Errorf("expected prefix 'su', got %q", got[:2])
	}
	if len(got) != len("supersecret") {
		t.Errorf("length mismatch: got %d, want %d", len(got), len("supersecret"))
	}
}

func TestRedactEnvs_ReturnsSensitiveKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"API_KEY": "abc123", "HOST": "localhost"},
		"dev":  {"PASSWORD": "hunter2", "PORT": "8080"},
	}
	results := RedactEnvs(envs)
	if len(results) != 2 {
		t.Fatalf("expected 2 redact results, got %d", len(results))
	}
	for _, r := range results {
		if r.Redacted == r.Original {
			t.Errorf("key %q was not redacted", r.Key)
		}
	}
}

func TestRedactEnvs_EmptyEnvs(t *testing.T) {
	results := RedactEnvs(map[string]map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
