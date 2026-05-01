package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"histkit/internal/sanitize"
)

func TestAppendCreatesLogFileAndDirectory(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "state", "audit.log")

	record := Record{
		RunID:       "run_001",
		StartedAt:   time.Date(2026, 5, 1, 13, 0, 0, 0, time.UTC),
		CompletedAt: time.Date(2026, 5, 1, 13, 0, 3, 0, time.UTC),
		Shell:       "bash",
		RuleNames:   []string{"secret-token", "trivial-cd"},
		CountsByAction: map[sanitize.ActionType]int{
			sanitize.ActionRedact:     1,
			sanitize.ActionQuarantine: 2,
		},
		CountsByConfidence: map[sanitize.Confidence]int{
			sanitize.ConfidenceHigh:   2,
			sanitize.ConfidenceMedium: 1,
		},
		BackupID: "b_20260501T130000Z_001",
		Apply:    true,
	}

	if err := Append(logPath, record); err != nil {
		t.Fatalf("Append returned error: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	want := "run_id=run_001 started_at=2026-05-01T13:00:00Z completed_at=2026-05-01T13:00:03Z shell=bash apply=true backup_id=b_20260501T130000Z_001 rules=secret-token,trivial-cd action_counts=keep:0,delete:0,redact:1,quarantine:2 confidence_counts=low:0,medium:1,high:2\n"
	if string(data) != want {
		t.Fatalf("log contents = %q, want %q", string(data), want)
	}

	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Stat returned error: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("file mode = %o, want 600", info.Mode().Perm())
	}
}

func TestAppendAppendsSecondRecord(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "audit.log")

	first := Record{
		RunID:       "run_001",
		StartedAt:   time.Date(2026, 5, 1, 13, 1, 0, 0, time.UTC),
		CompletedAt: time.Date(2026, 5, 1, 13, 1, 1, 0, time.UTC),
		Shell:       "bash",
		RuleNames:   []string{"secret-token"},
		CountsByAction: map[sanitize.ActionType]int{
			sanitize.ActionRedact: 1,
		},
		CountsByConfidence: map[sanitize.Confidence]int{
			sanitize.ConfidenceHigh: 1,
		},
		Apply: false,
	}
	second := Record{
		RunID:       "run_002",
		StartedAt:   time.Date(2026, 5, 1, 13, 2, 0, 0, time.UTC),
		CompletedAt: time.Date(2026, 5, 1, 13, 2, 4, 0, time.UTC),
		Shell:       "zsh",
		RuleNames:   []string{"trivial-cd"},
		CountsByAction: map[sanitize.ActionType]int{
			sanitize.ActionDelete: 1,
		},
		CountsByConfidence: map[sanitize.Confidence]int{
			sanitize.ConfidenceLow: 1,
		},
		BackupID: "b_20260501T130200Z_001",
		Apply:    true,
	}

	if err := Append(logPath, first); err != nil {
		t.Fatalf("Append first returned error: %v", err)
	}
	if err := Append(logPath, second); err != nil {
		t.Fatalf("Append second returned error: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	want := "" +
		"run_id=run_001 started_at=2026-05-01T13:01:00Z completed_at=2026-05-01T13:01:01Z shell=bash apply=false backup_id=- rules=secret-token action_counts=keep:0,delete:0,redact:1,quarantine:0 confidence_counts=low:0,medium:0,high:1\n" +
		"run_id=run_002 started_at=2026-05-01T13:02:00Z completed_at=2026-05-01T13:02:04Z shell=zsh apply=true backup_id=b_20260501T130200Z_001 rules=trivial-cd action_counts=keep:0,delete:1,redact:0,quarantine:0 confidence_counts=low:1,medium:0,high:0\n"
	if string(data) != want {
		t.Fatalf("appended log contents = %q, want %q", string(data), want)
	}
}

func TestAppendRejectsInvalidRecord(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "audit.log")

	if err := Append(logPath, Record{}); err == nil {
		t.Fatal("Append returned nil error for invalid record")
	}
}

func TestAppendRejectsEmptyPath(t *testing.T) {
	record := Record{
		RunID:       "run_001",
		StartedAt:   time.Date(2026, 5, 1, 13, 0, 0, 0, time.UTC),
		CompletedAt: time.Date(2026, 5, 1, 13, 0, 1, 0, time.UTC),
		Shell:       "bash",
		RuleNames:   []string{"secret-token"},
		CountsByAction: map[sanitize.ActionType]int{
			sanitize.ActionRedact: 1,
		},
		CountsByConfidence: map[sanitize.Confidence]int{
			sanitize.ConfidenceHigh: 1,
		},
	}

	if err := Append("", record); err == nil {
		t.Fatal("Append returned nil error for empty path")
	}
}
