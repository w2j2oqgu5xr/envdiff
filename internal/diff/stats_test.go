package diff

import (
	"testing"
)

func makeStatsResults() []Result {
	return []Result{
		{
			Key: "DB_HOST",
			Envs: map[string]Status{
				"prod": StatusMatch,
				"staging": StatusMatch,
			},
		},
		{
			Key: "API_KEY",
			Envs: map[string]Status{
				"prod": StatusMatch,
				"staging": StatusMissing,
			},
		},
		{
			Key: "LOG_LEVEL",
			Envs: map[string]Status{
				"prod": StatusMismatch,
				"staging": StatusMismatch,
			},
		},
	}
}

func TestComputeEnvStats_BasicCounts(t *testing.T) {
	results := makeStatsResults()
	stats := ComputeEnvStats(results)

	if len(stats) != 2 {
		t.Fatalf("expected 2 env stats, got %d", len(stats))
	}

	prod := stats[0] // sorted: prod
	if prod.EnvName != "prod" {
		t.Errorf("expected prod, got %s", prod.EnvName)
	}
	if prod.Total != 3 {
		t.Errorf("prod total: want 3, got %d", prod.Total)
	}
	if prod.Matched != 2 {
		t.Errorf("prod matched: want 2, got %d", prod.Matched)
	}
	if prod.Mismatched != 1 {
		t.Errorf("prod mismatched: want 1, got %d", prod.Mismatched)
	}
	if prod.Missing != 0 {
		t.Errorf("prod missing: want 0, got %d", prod.Missing)
	}
}

func TestComputeEnvStats_Coverage(t *testing.T) {
	results := makeStatsResults()
	stats := ComputeEnvStats(results)

	var staging EnvStats
	for _, s := range stats {
		if s.EnvName == "staging" {
			staging = s
		}
	}

	// staging has 1 missing out of 3 => coverage = 2/3 * 100 ≈ 66.67
	if staging.Coverage < 66.0 || staging.Coverage > 67.0 {
		t.Errorf("staging coverage: want ~66.67, got %.2f", staging.Coverage)
	}
}

func TestComputeEnvStats_Empty(t *testing.T) {
	stats := ComputeEnvStats([]Result{})
	if len(stats) != 0 {
		t.Errorf("expected empty stats, got %d", len(stats))
	}
}

func TestComputeEnvStats_SortedByName(t *testing.T) {
	results := makeStatsResults()
	stats := ComputeEnvStats(results)
	if stats[0].EnvName != "prod" || stats[1].EnvName != "staging" {
		t.Errorf("expected sorted order prod,staging; got %s,%s", stats[0].EnvName, stats[1].EnvName)
	}
}
