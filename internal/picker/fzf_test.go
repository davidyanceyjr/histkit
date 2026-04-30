package picker

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectReturnsErrorWhenFZFNotFound(t *testing.T) {
	t.Setenv("PATH", "")

	_, ok, err := Select(context.Background(), []Candidate{{Label: LabelHistory, Command: "git status"}})
	if err == nil {
		t.Fatal("Select returned nil error without fzf in PATH")
	}
	if ok {
		t.Fatal("Select returned ok=true without fzf in PATH")
	}
	if !strings.Contains(err.Error(), "find fzf") {
		t.Fatalf("Select error = %q, want mention of fzf lookup", err)
	}
}

func TestSelectReturnsChosenCandidate(t *testing.T) {
	fzfDir := t.TempDir()
	writeFakeFZF(t, fzfDir, "#!/bin/sh\nIFS= read -r _\nIFS= read -r line\nprintf '%s\\n' \"$line\"\n")
	t.Setenv("PATH", fzfDir)

	selected, ok, err := Select(context.Background(), []Candidate{
		{Label: LabelHistory, Command: "git status"},
		{Label: LabelSnippet, Command: "find {{path}}"},
	})
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if !ok {
		t.Fatal("Select returned ok=false, want true")
	}
	if selected.Label != LabelSnippet || selected.Command != "find {{path}}" {
		t.Fatalf("selected = %#v, want snippet candidate", selected)
	}
}

func TestSelectReturnsNoSelectionForAbort(t *testing.T) {
	fzfDir := t.TempDir()
	writeFakeFZF(t, fzfDir, "#!/bin/sh\nexit 130\n")
	t.Setenv("PATH", fzfDir)

	_, ok, err := Select(context.Background(), []Candidate{{Label: LabelHistory, Command: "git status"}})
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if ok {
		t.Fatal("Select returned ok=true, want false for abort")
	}
}

func writeFakeFZF(t *testing.T, dir, content string) {
	t.Helper()

	path := filepath.Join(dir, "fzf")
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
