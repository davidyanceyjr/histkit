package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davidyanceyjr/histkit/internal/config"
	"github.com/davidyanceyjr/histkit/internal/index"
)

func TestExecuteScanHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if err := Execute([]string{"scan", "--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	output := stdout.String()
	assertHelpContains(t, output,
		"Usage:\n  histkit scan [--shell <shell>] [--config <path>]",
		"scan reads supported history files and writes normalized entries into the local index.",
		"It does not rewrite shell history.",
		"--shell <shell>   scan only one supported shell source (bash or zsh)",
		"--config <path>   load a specific histkit config file before scanning",
	)
}

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

func TestExecuteScanConfigPathExpandsTilde(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	historyPath := filepath.Join(home, ".bash_history")
	if err := os.WriteFile(historyPath, []byte("pwd\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(history) returned error: %v", err)
	}

	configPath := filepath.Join(home, "histkit.toml")
	if err := os.WriteFile(configPath, []byte("[general]\ndefault_shell = \"bash\"\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(config) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"scan", "--config", "~/histkit.toml"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "scan complete: 1 source(s), 1 entries parsed, 1 inserted, 0 skipped, 0 warning(s).") {
		t.Fatalf("unexpected scan output: %q", stdout.String())
	}
}

func TestExecuteScanStreamsLargeHistoryInBatches(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var historyContent strings.Builder
	const entryCount = scanWriteBatchSize + 250
	for i := 0; i < entryCount; i++ {
		fmt.Fprintf(&historyContent, "printf 'entry-%d'\n", i)
	}

	historyPath := filepath.Join(home, ".bash_history")
	if err := os.WriteFile(historyPath, []byte(historyContent.String()), 0o600); err != nil {
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

	expectedOutput := fmt.Sprintf(
		"scan complete: 1 source(s), %d entries parsed, %d inserted, 0 skipped, 0 warning(s).",
		entryCount,
		entryCount,
	)
	if !strings.Contains(stdout.String(), expectedOutput) {
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

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM history_entries;`).Scan(&count); err != nil {
		t.Fatalf("QueryRow(count) returned error: %v", err)
	}
	if count != entryCount {
		t.Fatalf("history entry count = %d, want %d", count, entryCount)
	}
}

func TestExecuteScanAcceptsLongHistoryLine(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	longCommand := strings.Repeat("x", 128*1024)
	historyPath := filepath.Join(home, ".bash_history")
	if err := os.WriteFile(historyPath, []byte(longCommand+"\n"), 0o600); err != nil {
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
	if !strings.Contains(stdout.String(), "scan complete: 1 source(s), 1 entries parsed, 1 inserted, 0 skipped, 0 warning(s).") {
		t.Fatalf("unexpected scan output: %q", stdout.String())
	}
}

func TestExecuteScanRejectsMissingConfigPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"scan", "--config", "~/missing.toml"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for missing config path")
	}
	if !strings.Contains(err.Error(), `scan: load config "`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteScanRejectsInvalidConfigTOML(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configPath := filepath.Join(home, "histkit.toml")
	if err := os.WriteFile(configPath, []byte("[general\nbad"), 0o600); err != nil {
		t.Fatalf("WriteFile(config) returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"scan", "--config", configPath}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for invalid config TOML")
	}
	if !strings.Contains(err.Error(), `scan: load config "`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
