package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestBaselinePrinter_EmptyDiffs(t *testing.T) {
	var buf bytes.Buffer
	p := NewBaselinePrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "No environments") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestBaselinePrinter_FullyInSync(t *testing.T) {
	var buf bytes.Buffer
	p := NewBaselinePrinter(&buf)
	p.Print([]diff.BaselineDiff{
		{Baseline: "prod", Env: "staging"},
	})
	out := buf.String()
	if !strings.Contains(out, "fully in sync") {
		t.Errorf("expected sync message, got: %s", out)
	}
}

func TestBaselinePrinter_ShowsMissing(t *testing.T) {
	var buf bytes.Buffer
	p := NewBaselinePrinter(&buf)
	p.Print([]diff.BaselineDiff{
		{Baseline: "prod", Env: "dev", Missing: []string{"SECRET"}},
	})
	out := buf.String()
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING label, got: %s", out)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET key, got: %s", out)
	}
}

func TestBaselinePrinter_ShowsExtra(t *testing.T) {
	var buf bytes.Buffer
	p := NewBaselinePrinter(&buf)
	p.Print([]diff.BaselineDiff{
		{Baseline: "prod", Env: "dev", Extra: []string{"DEBUG"}},
	})
	out := buf.String()
	if !strings.Contains(out, "EXTRA") {
		t.Errorf("expected EXTRA label, got: %s", out)
	}
	if !strings.Contains(out, "DEBUG") {
		t.Errorf("expected DEBUG key, got: %s", out)
	}
}

func TestBaselinePrinter_ShowsChanged(t *testing.T) {
	var buf bytes.Buffer
	p := NewBaselinePrinter(&buf)
	p.Print([]diff.BaselineDiff{
		{Baseline: "prod", Env: "dev", Changed: []string{"DB_HOST"}},
	})
	out := buf.String()
	if !strings.Contains(out, "CHANGED") {
		t.Errorf("expected CHANGED label, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST key, got: %s", out)
	}
}

func TestBaselinePrinter_HeaderContainsEnvNames(t *testing.T) {
	var buf bytes.Buffer
	p := NewBaselinePrinter(&buf)
	p.Print([]diff.BaselineDiff{
		{Baseline: "prod", Env: "staging", Missing: []string{"KEY"}},
	})
	out := buf.String()
	if !strings.Contains(out, "prod") || !strings.Contains(out, "staging") {
		t.Errorf("expected env names in header, got: %s", out)
	}
}
