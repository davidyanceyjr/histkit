package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"histkit/internal/config"
	"histkit/internal/history"
	"histkit/internal/index"
	"histkit/internal/snippets"
)

func TestExecutePickPrintsSelectedCommand(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	seedPickFixtures(t, home)

	fzfDir := filepath.Join(home, "bin")
	if err := os.MkdirAll(fzfDir, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(fzfDir, "fzf"), []byte("#!/bin/sh\nIFS= read -r _\nIFS= read -r line\nprintf '%s\\n' \"$line\"\n"), 0o755); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	t.Setenv("PATH", fzfDir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"pick"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if got, want := stdout.String(), "echo custom\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestExecutePickFailsWhenFZFIsMissing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", "")

	seedPickFixtures(t, home)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"pick"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Execute returned nil error without fzf in PATH")
	}
	if !strings.Contains(err.Error(), "pick: find fzf") {
		t.Fatalf("Execute error = %q, want fzf lookup failure", err)
	}
}

func seedPickFixtures(t *testing.T, home string) {
	t.Helper()

	paths, err := config.DefaultPaths(home)
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}

	db, err := index.Open(paths.HistoryDB)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer db.Close()

	if err := index.InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	timestamp := time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC)
	if _, err := index.WriteHistoryEntries(db, []history.HistoryEntry{
		{
			Shell:      history.ShellBash,
			SourceFile: filepath.Join(home, ".bash_history"),
			RawLine:    "git status",
			Command:    "git status",
			Timestamp:  &timestamp,
		},
	}); err != nil {
		t.Fatalf("WriteHistoryEntries returned error: %v", err)
	}

	store := snippets.Store{Path: paths.SnippetsFile}
	if err := store.Save([]snippets.Snippet{
		{
			ID:          "echo-custom",
			Title:       "Echo custom",
			Command:     "echo custom",
			Description: "Print a custom value",
			Safety:      snippets.SafetyLow,
		},
	}); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
}
