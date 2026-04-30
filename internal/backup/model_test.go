package backup

import (
	"testing"
	"time"
)

func TestRecordValidate(t *testing.T) {
	record := Record{
		ID:         "b_20260430T120000Z_001",
		SourceFile: "/home/tester/.bash_history",
		BackupPath: "/home/tester/.local/share/histkit/backups/b_20260430T120000Z_001/.bash_history",
		CreatedAt:  time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC),
		Checksum:   "sha256:abc123",
	}

	if err := record.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestRecordValidateRequiresFields(t *testing.T) {
	tests := []struct {
		name   string
		record Record
	}{
		{name: "missing id", record: Record{SourceFile: "/tmp/a", BackupPath: "/tmp/b", CreatedAt: time.Now(), Checksum: "sum"}},
		{name: "missing source", record: Record{ID: "id", BackupPath: "/tmp/b", CreatedAt: time.Now(), Checksum: "sum"}},
		{name: "missing backup path", record: Record{ID: "id", SourceFile: "/tmp/a", CreatedAt: time.Now(), Checksum: "sum"}},
		{name: "missing created at", record: Record{ID: "id", SourceFile: "/tmp/a", BackupPath: "/tmp/b", Checksum: "sum"}},
		{name: "missing checksum", record: Record{ID: "id", SourceFile: "/tmp/a", BackupPath: "/tmp/b", CreatedAt: time.Now()}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.record.Validate(); err == nil {
				t.Fatal("Validate returned nil error")
			}
		})
	}
}

func TestBackupID(t *testing.T) {
	got := BackupID(time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC), 7)
	if got != "b_20260430T120000Z_007" {
		t.Fatalf("BackupID = %q, want b_20260430T120000Z_007", got)
	}
}

func TestBuildRecord(t *testing.T) {
	record, err := BuildRecord(
		"/home/tester/.bash_history",
		"/home/tester/.local/share/histkit/backups",
		"sha256:abc123",
		time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC),
		1,
	)
	if err != nil {
		t.Fatalf("BuildRecord returned error: %v", err)
	}

	if record.ID != "b_20260430T120000Z_001" {
		t.Fatalf("record.ID = %q, want b_20260430T120000Z_001", record.ID)
	}
	if record.BackupPath != "/home/tester/.local/share/histkit/backups/b_20260430T120000Z_001/.bash_history" {
		t.Fatalf("record.BackupPath = %q, want derived backup path", record.BackupPath)
	}
}

func TestBuildRecordRejectsInvalidInputs(t *testing.T) {
	if _, err := BuildRecord("", "/tmp/backups", "sha256:abc123", time.Now(), 1); err == nil {
		t.Fatal("BuildRecord returned nil error for empty source file")
	}
	if _, err := BuildRecord("/tmp/source", "", "sha256:abc123", time.Now(), 1); err == nil {
		t.Fatal("BuildRecord returned nil error for empty backup dir")
	}
	if _, err := BuildRecord("/tmp/source", "/tmp/backups", "", time.Now(), 1); err == nil {
		t.Fatal("BuildRecord returned nil error for empty checksum")
	}
	if _, err := BuildRecord("/tmp/source", "/tmp/backups", "sha256:abc123", time.Time{}, 1); err == nil {
		t.Fatal("BuildRecord returned nil error for zero created time")
	}
	if _, err := BuildRecord("/tmp/source", "/tmp/backups", "sha256:abc123", time.Now(), 0); err == nil {
		t.Fatal("BuildRecord returned nil error for non-positive sequence")
	}
}
