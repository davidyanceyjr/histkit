package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"histkit/internal/config"
)

func TestExecuteCleanDryRunOutputsPreviewWithoutMutatingHistory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	historyPath := filepath.Join(home, ".bash_history")
	original := "pwd\nmysql --password hunter2\necho hi\n"
	if err := os.WriteFile(historyPath, []byte(original), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"clean"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	output := stdout.String()
	for _, want := range []string{
		"source: shell=bash",
		"dry-run preview:",
		"counts by action:",
		"original: mysql --password hunter2",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got %q", want, output)
		}
	}

	gotHistory, err := os.ReadFile(historyPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(gotHistory) != original {
		t.Fatalf("history changed during dry run: got %q want %q", string(gotHistory), original)
	}
}

func TestExecuteCleanApplyRewritesHistoryCreatesBackupAndAudit(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	historyPath := filepath.Join(home, ".bash_history")
	if err := os.WriteFile(historyPath, []byte("pwd\nmysql --password hunter2\necho hi\n"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"clean", "--apply"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	if got := stdout.String(); !strings.Contains(got, "clean apply: shell=bash") {
		t.Fatalf("unexpected apply output: %q", got)
	}

	historyContent, err := os.ReadFile(historyPath)
	if err != nil {
		t.Fatalf("ReadFile(history) returned error: %v", err)
	}
	if got, want := string(historyContent), "mysql [REDACTED]\necho hi\n"; got != want {
		t.Fatalf("rewritten history = %q, want %q", got, want)
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}

	backupMatches, err := filepath.Glob(filepath.Join(paths.StateDir, "backups", "*", ".bash_history"))
	if err != nil {
		t.Fatalf("Glob returned error: %v", err)
	}
	if len(backupMatches) != 1 {
		t.Fatalf("backup match count = %d, want 1", len(backupMatches))
	}

	backupContent, err := os.ReadFile(backupMatches[0])
	if err != nil {
		t.Fatalf("ReadFile(backup) returned error: %v", err)
	}
	if got, want := string(backupContent), "pwd\nmysql --password hunter2\necho hi\n"; got != want {
		t.Fatalf("backup content = %q, want %q", got, want)
	}

	auditContent, err := os.ReadFile(paths.AuditLog)
	if err != nil {
		t.Fatalf("ReadFile(audit log) returned error: %v", err)
	}
	for _, want := range []string{
		"apply=true",
		"shell=bash",
		"backup_id=b_",
		"rules=inline-password-flag,pwd-command",
	} {
		if !strings.Contains(string(auditContent), want) {
			t.Fatalf("expected audit log to contain %q, got %q", want, string(auditContent))
		}
	}
}

func TestExecuteCleanApplyRequiresBackupHistory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	historyPath := filepath.Join(home, ".bash_history")
	if err := os.WriteFile(historyPath, []byte("pwd\n"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	configPath := filepath.Join(home, "histkit.toml")
	if err := os.WriteFile(configPath, []byte("[general]\nbackup_history = false\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(config) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"clean", "--apply", "--config", configPath}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error when backup_history is false")
	}
	if !strings.Contains(err.Error(), "backup_history=true") {
		t.Fatalf("unexpected error: %v", err)
	}
}
