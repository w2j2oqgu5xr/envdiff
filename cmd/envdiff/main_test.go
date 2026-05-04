package main

import "testing"

func TestSplitNamePath_Valid(t *testing.T) {
	tests := []struct {
		input    string
		wantName string
		wantPath string
		wantOK   bool
	}{
		{"dev=.env.dev", "dev", ".env.dev", true},
		{"prod=/etc/app/.env", "prod", "/etc/app/.env", true},
		{"staging=.env", "staging", ".env", true},
	}
	for _, tt := range tests {
		name, path, ok := splitNamePath(tt.input)
		if ok != tt.wantOK {
			t.Errorf("splitNamePath(%q) ok=%v, want %v", tt.input, ok, tt.wantOK)
		}
		if name != tt.wantName {
			t.Errorf("splitNamePath(%q) name=%q, want %q", tt.input, name, tt.wantName)
		}
		if path != tt.wantPath {
			t.Errorf("splitNamePath(%q) path=%q, want %q", tt.input, path, tt.wantPath)
		}
	}
}

func TestSplitNamePath_Invalid(t *testing.T) {
	inputs := []string{".env.dev", "noequalsign", ""}
	for _, input := range inputs {
		_, _, ok := splitNamePath(input)
		if ok {
			t.Errorf("splitNamePath(%q) should return ok=false", input)
		}
	}
}
