package backup

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "source.history")
	sourceContents := "echo test\nexport TOKEN=abc123\n"
	if err := os.WriteFile(sourcePath, []byte(sourceContents), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	record, err := Create(
		sourcePath,
		filepath.Join(tempDir, "backups"),
		time.Date(2026, 4, 30, 21, 30, 0, 0, time.UTC),
		1,
	)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if record.ID != "b_20260430T213000Z_001" {
		t.Fatalf("record.ID = %q, want b_20260430T213000Z_001", record.ID)
	}
	if record.SourceFile != sourcePath {
		t.Fatalf("record.SourceFile = %q, want %q", record.SourceFile, sourcePath)
	}

	data, err := os.ReadFile(record.BackupPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) != sourceContents {
		t.Fatalf("backup contents = %q, want %q", string(data), sourceContents)
	}

	checksum, err := ChecksumFile(record.BackupPath)
	if err != nil {
		t.Fatalf("ChecksumFile returned error: %v", err)
	}
	if checksum != record.Checksum {
		t.Fatalf("record.Checksum = %q, want %q", record.Checksum, checksum)
	}
}

func TestCreateCreatesBackupDirectory(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, ".zsh_history")
	if err := os.WriteFile(sourcePath, []byte("ls -la\n"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	backupDir := filepath.Join(tempDir, "nested", "backups")
	record, err := Create(sourcePath, backupDir, time.Date(2026, 4, 30, 21, 31, 0, 0, time.UTC), 2)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Dir(record.BackupPath)); err != nil {
		t.Fatalf("Stat returned error: %v", err)
	}
}

func TestCreateRejectsInvalidInputs(t *testing.T) {
	now := time.Now().UTC()

	if _, err := Create("", "/tmp/backups", now, 1); err == nil {
		t.Fatal("Create returned nil error for empty source path")
	}
	if _, err := Create("/tmp/source", "", now, 1); err == nil {
		t.Fatal("Create returned nil error for empty backup dir")
	}
	if _, err := Create("/tmp/source", "/tmp/backups", time.Time{}, 1); err == nil {
		t.Fatal("Create returned nil error for zero created time")
	}
	if _, err := Create("/tmp/source", "/tmp/backups", now, 0); err == nil {
		t.Fatal("Create returned nil error for non-positive sequence")
	}
	if _, err := Create("/tmp/does-not-exist", "/tmp/backups", now, 1); err == nil {
		t.Fatal("Create returned nil error for missing source file")
	}
}

func TestCreateRejectsExistingBackupPath(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "history")
	if err := os.WriteFile(sourcePath, []byte("pwd\n"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	backupDir := filepath.Join(tempDir, "backups")
	createdAt := time.Date(2026, 4, 30, 21, 32, 0, 0, time.UTC)
	record, err := BuildRecord(sourcePath, backupDir, "sha256:placeholder", createdAt, 1)
	if err != nil {
		t.Fatalf("BuildRecord returned error: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(record.BackupPath), 0o700); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(record.BackupPath, []byte("existing\n"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if _, err := Create(sourcePath, backupDir, createdAt, 1); err == nil {
		t.Fatal("Create returned nil error for existing backup path")
	}
}

func TestChecksumFile(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "sample")
	if err := os.WriteFile(path, []byte("abc"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	checksum, err := ChecksumFile(path)
	if err != nil {
		t.Fatalf("ChecksumFile returned error: %v", err)
	}

	const want = "sha256:ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if checksum != want {
		t.Fatalf("ChecksumFile = %q, want %q", checksum, want)
	}
}

func TestChecksumFileRejectsEmptyPath(t *testing.T) {
	if _, err := ChecksumFile(""); err == nil {
		t.Fatal("ChecksumFile returned nil error for empty path")
	}
}
