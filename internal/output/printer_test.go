package output

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrinter_EmptyResults(t *testing.T) {
	p := NewPrinter(Options{})
	out := captureStdout(func() {
		p.Print([]diff.Result{})
	})
	if out == "" {
		t.Error("expected some output for empty results")
	}
}

func TestPrinter_ShowsMissing(t *testing.T) {
	p := NewPrinter(Options{})
	results := []diff.Result{
		{Key: "SECRET", Status: diff.StatusMissing, MissingIn: "prod", Values: map[string]string{"dev": "abc"}},
	}
	out := captureStdout(func() {
		p.Print(results)
	})
	if out == "" {
		t.Error("expected output")
	}
}

func TestPrinter_ShowsMismatch(t *testing.T) {
	p := NewPrinter(Options{})
	results := []diff.Result{
		{Key: "DB_URL", Status: diff.StatusMismatch, Values: map[string]string{"dev": "localhost", "prod": "db.prod"}},
	}
	out := captureStdout(func() {
		p.Print(results)
	})
	if out == "" {
		t.Error("expected output")
	}
}

func TestPrinter_ShowMatchesWhenEnabled(t *testing.T) {
	p := NewPrinter(Options{ShowMatches: true})
	results := []diff.Result{
		{Key: "PORT", Status: diff.StatusMatch, Values: map[string]string{"dev": "8080", "prod": "8080"}},
	}
	out := captureStdout(func() {
		p.Print(results)
	})
	if out == "" {
		t.Error("expected output for matched key")
	}
}

func TestCollectEnvNames(t *testing.T) {
	results := []diff.Result{
		{Key: "X", Status: diff.StatusMismatch, Values: map[string]string{"alpha": "1", "beta": "2"}},
	}
	names := collectEnvNames(results)
	if len(names) != 2 {
		t.Errorf("expected 2 env names, got %d", len(names))
	}
}
