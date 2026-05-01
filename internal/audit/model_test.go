package audit

import (
	"testing"
	"time"

	"histkit/internal/sanitize"
)

func TestRecordValidate(t *testing.T) {
	record := Record{
		RunID:       "run_001",
		StartedAt:   time.Date(2026, 5, 1, 13, 0, 0, 0, time.UTC),
		CompletedAt: time.Date(2026, 5, 1, 13, 0, 2, 0, time.UTC),
		Shell:       "bash",
		RuleNames:   []string{"secret-token", "trivial-cd"},
		CountsByAction: map[sanitize.ActionType]int{
			sanitize.ActionRedact:     1,
			sanitize.ActionQuarantine: 1,
		},
		CountsByConfidence: map[sanitize.Confidence]int{
			sanitize.ConfidenceHigh:   1,
			sanitize.ConfidenceMedium: 1,
		},
		BackupID: "b_20260501T130000Z_001",
		Apply:    true,
	}

	if err := record.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestRecordValidateRejectsInvalidInputs(t *testing.T) {
	now := time.Date(2026, 5, 1, 13, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		record Record
	}{
		{name: "missing run id", record: Record{StartedAt: now, CompletedAt: now, Shell: "bash", RuleNames: []string{"rule"}}},
		{name: "missing started time", record: Record{RunID: "run", CompletedAt: now, Shell: "bash", RuleNames: []string{"rule"}}},
		{name: "missing completed time", record: Record{RunID: "run", StartedAt: now, Shell: "bash", RuleNames: []string{"rule"}}},
		{name: "completed before started", record: Record{RunID: "run", StartedAt: now, CompletedAt: now.Add(-time.Second), Shell: "bash", RuleNames: []string{"rule"}}},
		{name: "missing shell", record: Record{RunID: "run", StartedAt: now, CompletedAt: now, RuleNames: []string{"rule"}}},
		{name: "missing rules", record: Record{RunID: "run", StartedAt: now, CompletedAt: now, Shell: "bash"}},
		{name: "empty rule", record: Record{RunID: "run", StartedAt: now, CompletedAt: now, Shell: "bash", RuleNames: []string{" "}}},
		{name: "duplicate rule", record: Record{RunID: "run", StartedAt: now, CompletedAt: now, Shell: "bash", RuleNames: []string{"rule", "rule"}}},
		{name: "invalid action", record: Record{RunID: "run", StartedAt: now, CompletedAt: now, Shell: "bash", RuleNames: []string{"rule"}, CountsByAction: map[sanitize.ActionType]int{"bad": 1}}},
		{name: "negative action count", record: Record{RunID: "run", StartedAt: now, CompletedAt: now, Shell: "bash", RuleNames: []string{"rule"}, CountsByAction: map[sanitize.ActionType]int{sanitize.ActionDelete: -1}}},
		{name: "invalid confidence", record: Record{RunID: "run", StartedAt: now, CompletedAt: now, Shell: "bash", RuleNames: []string{"rule"}, CountsByConfidence: map[sanitize.Confidence]int{"bad": 1}}},
		{name: "negative confidence count", record: Record{RunID: "run", StartedAt: now, CompletedAt: now, Shell: "bash", RuleNames: []string{"rule"}, CountsByConfidence: map[sanitize.Confidence]int{sanitize.ConfidenceHigh: -1}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.record.Validate(); err == nil {
				t.Fatal("Validate returned nil error")
			}
		})
	}
}

func TestRenderLineIsDeterministic(t *testing.T) {
	record := Record{
		RunID:       "run_001",
		StartedAt:   time.Date(2026, 5, 1, 13, 0, 0, 0, time.UTC),
		CompletedAt: time.Date(2026, 5, 1, 13, 0, 2, 0, time.UTC),
		Shell:       "zsh",
		RuleNames:   []string{"trivial-cd", "secret-token"},
		CountsByAction: map[sanitize.ActionType]int{
			sanitize.ActionRedact:     2,
			sanitize.ActionQuarantine: 1,
		},
		CountsByConfidence: map[sanitize.Confidence]int{
			sanitize.ConfidenceHigh:   2,
			sanitize.ConfidenceMedium: 1,
		},
		Apply: true,
	}

	got := RenderLine(record)
	want := "run_id=run_001 started_at=2026-05-01T13:00:00Z completed_at=2026-05-01T13:00:02Z shell=zsh apply=true backup_id=- rules=secret-token,trivial-cd action_counts=keep:0,delete:0,redact:2,quarantine:1 confidence_counts=low:0,medium:1,high:2"
	if got != want {
		t.Fatalf("RenderLine = %q, want %q", got, want)
	}
}
