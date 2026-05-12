package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davidyanceyjr/histkit/internal/config"
)

func TestExecuteCleanHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if err := Execute([]string{"clean", "--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	output := stdout.String()
	assertHelpContains(t, output,
		"Usage:\n  histkit clean [--apply] [--dry-run] [--shell <shell>] [--config <path>]",
		"Without --apply, clean runs in preview mode and prints the planned actions without changing history files.",
		"--apply rewrites matching history sources, creates backups, and appends an audit record.",
		"--dry-run         render the cleanup preview explicitly without changing files",
		"--apply           rewrite matching history files; requires backup_history=true",
	)
}

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
	if got, want := string(historyContent), "mysql --password [REDACTED]\necho hi\n"; got != want {
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

func TestExecuteCleanRejectsApplyAndDryRunTogether(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Execute([]string{"clean", "--apply", "--dry-run"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for mutually exclusive flags")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("unexpected error: %v", err)
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

func TestExecuteCleanApplyNoMatchesDoesNotCreateBackupOrAudit(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	historyPath := filepath.Join(home, ".bash_history")
	original := "echo hi\nprintf 'done\\n'\n"
	if err := os.WriteFile(historyPath, []byte(original), 0o600); err != nil {
		t.Fatalf("WriteFile(history) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"clean", "--apply"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "no matching entries") {
		t.Fatalf("unexpected output: %q", stdout.String())
	}

	historyContent, err := os.ReadFile(historyPath)
	if err != nil {
		t.Fatalf("ReadFile(history) returned error: %v", err)
	}
	if string(historyContent) != original {
		t.Fatalf("history content = %q, want %q", string(historyContent), original)
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}

	backupMatches, err := filepath.Glob(filepath.Join(paths.StateDir, "backups", "*"))
	if err != nil {
		t.Fatalf("Glob returned error: %v", err)
	}
	if len(backupMatches) != 0 {
		t.Fatalf("backup match count = %d, want 0", len(backupMatches))
	}
	if _, err := os.Stat(paths.AuditLog); !os.IsNotExist(err) {
		t.Fatalf("audit log stat error = %v, want not exists", err)
	}
}

func TestExecuteCleanConfigPathExpandsTilde(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configPath := filepath.Join(home, "histkit.toml")
	if err := os.WriteFile(configPath, []byte("[general]\nbackup_history = true\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(config) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"clean", "--config", "~/histkit.toml"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "clean: no history sources detected") {
		t.Fatalf("unexpected output: %q", stdout.String())
	}
}

func TestExecuteCleanRejectsMissingConfigPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"clean", "--config", "~/missing.toml"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for missing config path")
	}
	if !strings.Contains(err.Error(), `clean: load config "`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteCleanRejectsInvalidConfigTOML(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configPath := filepath.Join(home, "histkit.toml")
	if err := os.WriteFile(configPath, []byte("[general\nbad"), 0o600); err != nil {
		t.Fatalf("WriteFile(config) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"clean", "--config", configPath}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for invalid config TOML")
	}
	if !strings.Contains(err.Error(), `clean: load config "`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteCleanDryRunFlagMatchesDefaultPlanningMode(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	historyPath := filepath.Join(home, ".bash_history")
	original := "pwd\nmysql --password hunter2\necho hi\n"
	if err := os.WriteFile(historyPath, []byte(original), 0o600); err != nil {
		t.Fatalf("WriteFile(history) returned error: %v", err)
	}

	var defaultStdout bytes.Buffer
	var defaultStderr bytes.Buffer
	if err := Execute([]string{"clean"}, &defaultStdout, &defaultStderr); err != nil {
		t.Fatalf("Execute(default clean) returned error: %v", err)
	}
	if defaultStderr.Len() != 0 {
		t.Fatalf("expected no stderr output for default clean, got %q", defaultStderr.String())
	}

	var dryRunStdout bytes.Buffer
	var dryRunStderr bytes.Buffer
	if err := Execute([]string{"clean", "--dry-run"}, &dryRunStdout, &dryRunStderr); err != nil {
		t.Fatalf("Execute(clean --dry-run) returned error: %v", err)
	}
	if dryRunStderr.Len() != 0 {
		t.Fatalf("expected no stderr output for clean --dry-run, got %q", dryRunStderr.String())
	}

	if got, want := dryRunStdout.String(), defaultStdout.String(); got != want {
		t.Fatalf("clean --dry-run output = %q, want %q", got, want)
	}

	historyContent, err := os.ReadFile(historyPath)
	if err != nil {
		t.Fatalf("ReadFile(history) returned error: %v", err)
	}
	if string(historyContent) != original {
		t.Fatalf("history content = %q, want %q", string(historyContent), original)
	}
}

func TestExecuteCleanApplyShellFlagFiltersToZshOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	bashHistoryPath := filepath.Join(home, ".bash_history")
	bashOriginal := "pwd\nmysql --password hunter2\n"
	if err := os.WriteFile(bashHistoryPath, []byte(bashOriginal), 0o600); err != nil {
		t.Fatalf("WriteFile(bash history) returned error: %v", err)
	}

	zshHistoryPath := filepath.Join(home, ".zsh_history")
	zshOriginal := ": 1712959000:0;pwd\n: 1712959015:2;curl https://user:secret@example.com/path\n"
	if err := os.WriteFile(zshHistoryPath, []byte(zshOriginal), 0o600); err != nil {
		t.Fatalf("WriteFile(zsh history) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"clean", "--apply", "--shell", "zsh"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, "clean apply: shell=zsh") {
		t.Fatalf("unexpected apply output: %q", got)
	}
	if strings.Contains(stdout.String(), "shell=bash") {
		t.Fatalf("unexpected bash apply output: %q", stdout.String())
	}

	bashHistoryContent, err := os.ReadFile(bashHistoryPath)
	if err != nil {
		t.Fatalf("ReadFile(bash history) returned error: %v", err)
	}
	if string(bashHistoryContent) != bashOriginal {
		t.Fatalf("bash history content = %q, want %q", string(bashHistoryContent), bashOriginal)
	}

	zshHistoryContent, err := os.ReadFile(zshHistoryPath)
	if err != nil {
		t.Fatalf("ReadFile(zsh history) returned error: %v", err)
	}
	if got, want := string(zshHistoryContent), ": 1712959015:2;curl [REDACTED]/path\n"; got != want {
		t.Fatalf("zsh history content = %q, want %q", got, want)
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}

	zshBackupMatches, err := filepath.Glob(filepath.Join(paths.StateDir, "backups", "*", ".zsh_history"))
	if err != nil {
		t.Fatalf("Glob(zsh backups) returned error: %v", err)
	}
	if len(zshBackupMatches) != 1 {
		t.Fatalf("zsh backup match count = %d, want 1", len(zshBackupMatches))
	}

	bashBackupMatches, err := filepath.Glob(filepath.Join(paths.StateDir, "backups", "*", ".bash_history"))
	if err != nil {
		t.Fatalf("Glob(bash backups) returned error: %v", err)
	}
	if len(bashBackupMatches) != 0 {
		t.Fatalf("bash backup match count = %d, want 0", len(bashBackupMatches))
	}

	auditContent, err := os.ReadFile(paths.AuditLog)
	if err != nil {
		t.Fatalf("ReadFile(audit log) returned error: %v", err)
	}
	if !strings.Contains(string(auditContent), "shell=zsh") {
		t.Fatalf("expected audit log to contain shell=zsh, got %q", string(auditContent))
	}
	if strings.Contains(string(auditContent), "shell=bash") {
		t.Fatalf("unexpected bash audit content: %q", string(auditContent))
	}
}

func TestExecuteCleanApplyReturnsErrorButKeepsRewriteWhenAuditAppendFails(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	historyPath := filepath.Join(home, ".bash_history")
	original := "pwd\nmysql --password hunter2\necho hi\n"
	if err := os.WriteFile(historyPath, []byte(original), 0o600); err != nil {
		t.Fatalf("WriteFile(history) returned error: %v", err)
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}
	if err := os.MkdirAll(paths.StateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll(state dir) returned error: %v", err)
	}
	if err := os.Mkdir(paths.AuditLog, 0o700); err != nil {
		t.Fatalf("Mkdir(audit log blocker) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err = Execute([]string{"clean", "--apply"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected audit append failure")
	}
	if !strings.Contains(err.Error(), "append audit log") {
		t.Fatalf("unexpected error: %v", err)
	}

	historyContent, err := os.ReadFile(historyPath)
	if err != nil {
		t.Fatalf("ReadFile(history) returned error: %v", err)
	}
	if got, want := string(historyContent), "mysql --password [REDACTED]\necho hi\n"; got != want {
		t.Fatalf("rewritten history = %q, want %q", got, want)
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
	if string(backupContent) != original {
		t.Fatalf("backup content = %q, want %q", string(backupContent), original)
	}
}
