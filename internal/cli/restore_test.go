package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/davidyanceyjr/histkit/internal/backup"
	"github.com/davidyanceyjr/histkit/internal/config"
)

func TestExecuteRestoreHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if err := Execute([]string{"restore", "--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	output := stdout.String()
	assertHelpContains(t, output,
		"Usage:\n  histkit restore [--config <path>] [backup-id]",
		"Without a backup ID, restore lists the backups available under the histkit state directory.",
		"With a backup ID, restore replaces the source history file from that backup and appends an audit record.",
		"--config <path>   load a specific histkit config file before listing or restoring backups",
	)
}

func TestExecuteRestoreListsBackupsWhenNoIDProvided(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}

	sourcePath := filepath.Join(home, ".bash_history")
	if err := os.WriteFile(sourcePath, []byte("pwd\n"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if _, err := backup.Create(sourcePath, filepath.Join(paths.StateDir, "backups"), time.Date(2026, 5, 1, 15, 0, 0, 0, time.UTC), 1); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"restore"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "backup=b_20260501T150000Z_001") {
		t.Fatalf("unexpected output: %q", stdout.String())
	}
}

func TestExecuteRestoreRestoresSpecificBackupAndAudits(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}
	sourcePath := filepath.Join(home, ".bash_history")
	if err := os.WriteFile(sourcePath, []byte("original\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(original) returned error: %v", err)
	}

	record, err := backup.Create(sourcePath, filepath.Join(paths.StateDir, "backups"), time.Date(2026, 5, 1, 15, 5, 0, 0, time.UTC), 1)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := os.WriteFile(sourcePath, []byte("modified\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(modified) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"restore", record.ID}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "restore complete: backup="+record.ID) {
		t.Fatalf("unexpected output: %q", stdout.String())
	}

	data, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatalf("ReadFile(source) returned error: %v", err)
	}
	if string(data) != "original\n" {
		t.Fatalf("restored contents = %q, want %q", string(data), "original\n")
	}

	auditData, err := os.ReadFile(paths.AuditLog)
	if err != nil {
		t.Fatalf("ReadFile(audit) returned error: %v", err)
	}
	for _, want := range []string{"apply=false", "backup_id=" + record.ID, "rules=restore"} {
		if !strings.Contains(string(auditData), want) {
			t.Fatalf("expected audit log to contain %q, got %q", want, string(auditData))
		}
	}
}

func TestExecuteRestoreNoBackups(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"restore"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !strings.Contains(stdout.String(), "restore: no backups available") {
		t.Fatalf("unexpected output: %q", stdout.String())
	}
}

func TestExecuteRestoreConfigPathExpandsTilde(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configPath := filepath.Join(home, "histkit.toml")
	if err := os.WriteFile(configPath, []byte("[general]\nbackup_history = true\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(config) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"restore", "--config", "~/histkit.toml"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "restore: no backups available") {
		t.Fatalf("unexpected output: %q", stdout.String())
	}
}

func TestExecuteRestoreRejectsMissingConfigPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"restore", "--config", "~/missing.toml"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for missing config path")
	}
	if !strings.Contains(err.Error(), `restore: load config "`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteRestoreRejectsInvalidConfigTOML(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configPath := filepath.Join(home, "histkit.toml")
	if err := os.WriteFile(configPath, []byte("[general\nbad"), 0o600); err != nil {
		t.Fatalf("WriteFile(config) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"restore", "--config", configPath}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for invalid config TOML")
	}
	if !strings.Contains(err.Error(), `restore: load config "`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteRestoreMissingBackupIDLeavesHistoryUntouched(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	sourcePath := filepath.Join(home, ".bash_history")
	original := "current\n"
	if err := os.WriteFile(sourcePath, []byte(original), 0o600); err != nil {
		t.Fatalf("WriteFile(source) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"restore", "b_20260501T150000Z_001"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for missing backup ID")
	}
	if !strings.Contains(err.Error(), "find backup record") {
		t.Fatalf("unexpected error: %v", err)
	}

	data, readErr := os.ReadFile(sourcePath)
	if readErr != nil {
		t.Fatalf("ReadFile(source) returned error: %v", readErr)
	}
	if string(data) != original {
		t.Fatalf("source contents = %q, want %q", string(data), original)
	}
}

func TestExecuteRestoreReturnsErrorButKeepsRestoredFileWhenAuditAppendFails(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}
	sourcePath := filepath.Join(home, ".bash_history")
	if err := os.WriteFile(sourcePath, []byte("original\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(original) returned error: %v", err)
	}

	record, err := backup.Create(sourcePath, filepath.Join(paths.StateDir, "backups"), time.Date(2026, 5, 1, 15, 6, 0, 0, time.UTC), 1)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if err := os.WriteFile(sourcePath, []byte("modified\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(modified) returned error: %v", err)
	}

	if err := os.Mkdir(paths.AuditLog, 0o700); err != nil {
		t.Fatalf("Mkdir(audit log blocker) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err = Execute([]string{"restore", record.ID}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected audit append failure")
	}
	if !strings.Contains(err.Error(), "append audit log") {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatalf("ReadFile(source) returned error: %v", err)
	}
	if string(data) != "original\n" {
		t.Fatalf("restored contents = %q, want %q", string(data), "original\n")
	}
}
