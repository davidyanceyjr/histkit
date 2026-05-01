package backup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteRecordAndLoadRecord(t *testing.T) {
	tempDir := t.TempDir()
	record := Record{
		ID:         "b_20260501T140000Z_001",
		SourceFile: filepath.Join(tempDir, ".bash_history"),
		BackupPath: filepath.Join(tempDir, "backups", "b_20260501T140000Z_001", ".bash_history"),
		CreatedAt:  time.Date(2026, 5, 1, 14, 0, 0, 0, time.UTC),
		Checksum:   "sha256:abc",
	}

	if err := WriteRecord(record, filepath.Join(tempDir, "backups")); err != nil {
		t.Fatalf("WriteRecord returned error: %v", err)
	}

	loaded, err := LoadRecord(filepath.Join(tempDir, "backups", record.ID, recordFileName))
	if err != nil {
		t.Fatalf("LoadRecord returned error: %v", err)
	}
	if loaded != record {
		t.Fatalf("loaded record = %#v, want %#v", loaded, record)
	}
}

func TestListRecordsSortsNewestFirst(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")

	older := Record{
		ID:         "b_20260501T130000Z_001",
		SourceFile: "/home/tester/.bash_history",
		BackupPath: filepath.Join(backupDir, "b_20260501T130000Z_001", ".bash_history"),
		CreatedAt:  time.Date(2026, 5, 1, 13, 0, 0, 0, time.UTC),
		Checksum:   "sha256:older",
	}
	newer := Record{
		ID:         "b_20260501T140000Z_001",
		SourceFile: "/home/tester/.zsh_history",
		BackupPath: filepath.Join(backupDir, "b_20260501T140000Z_001", ".zsh_history"),
		CreatedAt:  time.Date(2026, 5, 1, 14, 0, 0, 0, time.UTC),
		Checksum:   "sha256:newer",
	}

	for _, record := range []Record{older, newer} {
		if err := WriteRecord(record, backupDir); err != nil {
			t.Fatalf("WriteRecord returned error: %v", err)
		}
	}

	records, err := ListRecords(backupDir)
	if err != nil {
		t.Fatalf("ListRecords returned error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("len(records) = %d, want 2", len(records))
	}
	if records[0].ID != newer.ID || records[1].ID != older.ID {
		t.Fatalf("record order = [%s %s], want [%s %s]", records[0].ID, records[1].ID, newer.ID, older.ID)
	}
}

func TestRestoreRewritesTargetFromBackup(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, ".bash_history")
	backupDir := filepath.Join(tempDir, "backups")

	if err := os.WriteFile(sourcePath, []byte("current\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(source) returned error: %v", err)
	}

	record, err := Create(sourcePath, backupDir, time.Date(2026, 5, 1, 14, 5, 0, 0, time.UTC), 1)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if err := os.WriteFile(sourcePath, []byte("modified\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(modified source) returned error: %v", err)
	}

	if err := Restore(record); err != nil {
		t.Fatalf("Restore returned error: %v", err)
	}

	data, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) != "current\n" {
		t.Fatalf("restored contents = %q, want %q", string(data), "current\n")
	}
}

func TestRestoreRejectsChecksumMismatch(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, ".bash_history")
	backupDir := filepath.Join(tempDir, "backups")

	if err := os.WriteFile(sourcePath, []byte("current\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(source) returned error: %v", err)
	}

	record, err := Create(sourcePath, backupDir, time.Date(2026, 5, 1, 14, 6, 0, 0, time.UTC), 1)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if err := os.WriteFile(record.BackupPath, []byte("tampered\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(backup) returned error: %v", err)
	}

	err = Restore(record)
	if err == nil {
		t.Fatal("expected checksum mismatch error")
	}
	if !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("unexpected error: %v", err)
	}

	data, readErr := os.ReadFile(sourcePath)
	if readErr != nil {
		t.Fatalf("ReadFile(source) returned error: %v", readErr)
	}
	if string(data) != "current\n" {
		t.Fatalf("source contents changed on failed restore: got %q want %q", string(data), "current\n")
	}
}
