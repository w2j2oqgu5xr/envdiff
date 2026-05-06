package diff

import (
	"testing"
)

func TestComputeOverlap_SinglePair(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"A": "1", "B": "2", "C": "3"},
		"dev":  {"A": "1", "C": "x", "D": "4"},
	}
	results := ComputeOverlap(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(results))
	}
	r := results[0]
	if r.EnvA != "dev" || r.EnvB != "prod" {
		t.Errorf("unexpected pair names: %s / %s", r.EnvA, r.EnvB)
	}
	if len(r.SharedKeys) != 2 {
		t.Errorf("expected 2 shared keys, got %d: %v", len(r.SharedKeys), r.SharedKeys)
	}
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "D" {
		t.Errorf("expected OnlyInA=[D], got %v", r.OnlyInA)
	}
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "B" {
		t.Errorf("expected OnlyInB=[B], got %v", r.OnlyInB)
	}
}

func TestComputeOverlap_OverlapPercent(t *testing.T) {
	envs := map[string]map[string]string{
		"a": {"X": "1", "Y": "2"},
		"b": {"X": "1", "Y": "2"},
	}
	results := ComputeOverlap(envs)
	if results[0].OverlapPct != 100.0 {
		t.Errorf("expected 100%% overlap, got %.2f", results[0].OverlapPct)
	}
}

func TestComputeOverlap_NoSharedKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"a": {"FOO": "1"},
		"b": {"BAR": "2"},
	}
	results := ComputeOverlap(envs)
	if results[0].OverlapPct != 0.0 {
		t.Errorf("expected 0%% overlap, got %.2f", results[0].OverlapPct)
	}
	if len(results[0].SharedKeys) != 0 {
		t.Errorf("expected no shared keys")
	}
}

func TestComputeOverlap_MultiplePairs(t *testing.T) {
	envs := map[string]map[string]string{
		"a": {"K": "1"},
		"b": {"K": "1"},
		"c": {"K": "1"},
	}
	results := ComputeOverlap(envs)
	// 3 envs => 3 pairs
	if len(results) != 3 {
		t.Errorf("expected 3 pairs, got %d", len(results))
	}
}

func TestComputeOverlap_Empty(t *testing.T) {
	results := ComputeOverlap(map[string]map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected no results for empty input")
	}
}
