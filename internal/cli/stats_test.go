package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteStatsEmptyIndex(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"stats"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Indexed history entries: 0") {
		t.Fatalf("unexpected stats output: %q", output)
	}
	if !strings.Contains(output, "Counts by shell:\n  (none)\n") {
		t.Fatalf("expected empty shell counts, got %q", output)
	}
	if !strings.Contains(output, "Counts by source:\n  (none)\n") {
		t.Fatalf("expected empty source counts, got %q", output)
	}
}

func TestExecuteStatsReportsShellAndSourceCounts(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := os.WriteFile(filepath.Join(home, ".bash_history"), []byte("pwd\ngit status\n"), 0o600); err != nil {
		t.Fatalf("WriteFile bash returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(": 1712959000:0;echo zsh\n"), 0o600); err != nil {
		t.Fatalf("WriteFile zsh returned error: %v", err)
	}

	var scanStdout bytes.Buffer
	var scanStderr bytes.Buffer
	if err := Execute([]string{"scan"}, &scanStdout, &scanStderr); err != nil {
		t.Fatalf("Execute(scan) returned error: %v", err)
	}
	if scanStderr.Len() != 0 {
		t.Fatalf("expected no scan stderr output, got %q", scanStderr.String())
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"stats"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute(stats) returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	output := stdout.String()
	for _, want := range []string{
		"Indexed history entries: 3",
		"Counts by shell:\n  bash: 2\n  zsh: 1\n",
		"Counts by source:\n  " + filepath.Join(home, ".bash_history") + ": 2\n  " + filepath.Join(home, ".zsh_history") + ": 1\n",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected stats output to contain %q, got %q", want, output)
		}
	}
}
