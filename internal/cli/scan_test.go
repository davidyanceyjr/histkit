package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"histkit/internal/config"
	"histkit/internal/index"
)

func TestExecuteScanIndexesBashHistory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	historyPath := filepath.Join(home, ".bash_history")
	if err := os.WriteFile(historyPath, []byte("pwd\n   \ngit status\n"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"scan"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "scan complete: 1 source(s), 2 entries parsed, 2 inserted, 0 skipped, 1 warning(s).") {
		t.Fatalf("unexpected scan output: %q", output)
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}

	db, err := index.Open(paths.HistoryDB)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM history_entries;`).Scan(&count); err != nil {
		t.Fatalf("QueryRow(count) returned error: %v", err)
	}
	if count != 2 {
		t.Fatalf("history entry count = %d, want 2", count)
	}
}

func TestExecuteScanShellFlagFiltersSources(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := os.WriteFile(filepath.Join(home, ".bash_history"), []byte("pwd\n"), 0o600); err != nil {
		t.Fatalf("WriteFile bash returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(": 1712959000:0;echo zsh\n"), 0o600); err != nil {
		t.Fatalf("WriteFile zsh returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"scan", "--shell", "zsh"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "scan complete: 1 source(s), 1 entries parsed, 1 inserted, 0 skipped, 0 warning(s).") {
		t.Fatalf("unexpected scan output: %q", stdout.String())
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}

	db, err := index.Open(paths.HistoryDB)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer db.Close()

	var shell string
	var count int
	if err := db.QueryRow(`SELECT shell, COUNT(*) FROM history_entries GROUP BY shell;`).Scan(&shell, &count); err != nil {
		t.Fatalf("QueryRow(shell count) returned error: %v", err)
	}
	if shell != "zsh" || count != 1 {
		t.Fatalf("stored rows = shell %q count %d, want shell zsh count 1", shell, count)
	}
}

func TestExecuteScanRejectsUnsupportedShell(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"scan", "--shell", "fish"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
	if !strings.Contains(err.Error(), `unsupported shell "fish"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
