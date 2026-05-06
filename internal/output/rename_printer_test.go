package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestRenamePrinter_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	p := NewRenamePrinter(&buf)
	p.Print(nil)
	out := buf.String()
	if !strings.Contains(out, "No rename issues") {
		t.Errorf("expected no-issue message, got: %q", out)
	}
}

func TestRenamePrinter_ShowsHeader(t *testing.T) {
	var buf bytes.Buffer
	p := NewRenamePrinter(&buf)
	p.Print([]diff.RenameResult{
		{OldKey: "API_KEY", NewKey: "API_TOKEN", EnvName: "production", Found: true},
	})
	out := buf.String()
	if !strings.Contains(out, "renamed keys") {
		t.Errorf("expected header about renamed keys, got: %q", out)
	}
}

func TestRenamePrinter_ShowsEnvName(t *testing.T) {
	var buf bytes.Buffer
	p := NewRenamePrinter(&buf)
	p.Print([]diff.RenameResult{
		{OldKey: "OLD", NewKey: "NEW", EnvName: "staging", Found: true},
	})
	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected env name 'staging' in output, got: %q", out)
	}
}

func TestRenamePrinter_ShowsKeys(t *testing.T) {
	var buf bytes.Buffer
	p := NewRenamePrinter(&buf)
	p.Print([]diff.RenameResult{
		{OldKey: "DB_HOST", NewKey: "DATABASE_HOST", EnvName: "dev", Found: true},
	})
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected old key in output, got: %q", out)
	}
	if !strings.Contains(out, "DATABASE_HOST") {
		t.Errorf("expected new key in output, got: %q", out)
	}
}

func TestRenamePrinter_GroupsByEnv(t *testing.T) {
	var buf bytes.Buffer
	p := NewRenamePrinter(&buf)
	p.Print([]diff.RenameResult{
		{OldKey: "A", NewKey: "B", EnvName: "alpha", Found: true},
		{OldKey: "X", NewKey: "Y", EnvName: "beta", Found: true},
	})
	out := buf.String()
	alphaIdx := strings.Index(out, "alpha")
	betaIdx := strings.Index(out, "beta")
	if alphaIdx == -1 || betaIdx == -1 {
		t.Errorf("expected both env names in output, got: %q", out)
	}
	if alphaIdx > betaIdx {
		t.Errorf("expected 'alpha' to appear before 'beta'")
	}
}
