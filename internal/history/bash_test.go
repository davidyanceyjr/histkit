package history

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseBashFixture(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "history", "bash", "plain.hist")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer f.Close()

	entries, warnings, err := ParseBash("/home/tester/.bash_history", f)
	if err != nil {
		t.Fatalf("ParseBash returned error: %v", err)
	}

	if len(warnings) != 1 {
		t.Fatalf("len(warnings) = %d, want 1", len(warnings))
	}
	if got, want := warnings[0].LineNumber, 4; got != want {
		t.Fatalf("warning line number = %d, want %d", got, want)
	}

	if len(entries) != 4 {
		t.Fatalf("len(entries) = %d, want 4", len(entries))
	}

	tests := []struct {
		index   int
		rawLine string
		command string
	}{
		{index: 0, rawLine: "git status", command: "git status"},
		{index: 1, rawLine: "echo 'hello world'", command: "echo 'hello world'"},
		{index: 2, rawLine: "  ls -la", command: "  ls -la"},
		{index: 3, rawLine: "printf '%s\\n' foo; printf '%s\\n' bar", command: "printf '%s\\n' foo; printf '%s\\n' bar"},
	}

	for _, tc := range tests {
		entry := entries[tc.index]
		if got, want := entry.Shell, ShellBash; got != want {
			t.Fatalf("entries[%d].Shell = %q, want %q", tc.index, got, want)
		}
		if got, want := entry.SourceFile, "/home/tester/.bash_history"; got != want {
			t.Fatalf("entries[%d].SourceFile = %q, want %q", tc.index, got, want)
		}
		if got, want := entry.RawLine, tc.rawLine; got != want {
			t.Fatalf("entries[%d].RawLine = %q, want %q", tc.index, got, want)
		}
		if got, want := entry.Command, tc.command; got != want {
			t.Fatalf("entries[%d].Command = %q, want %q", tc.index, got, want)
		}
		if entry.HasTimestamp() {
			t.Fatalf("entries[%d].HasTimestamp = true, want false", tc.index)
		}
		if entry.HasExitCode() {
			t.Fatalf("entries[%d].HasExitCode = true, want false", tc.index)
		}
	}
}

func TestParseBashEmptyInput(t *testing.T) {
	entries, warnings, err := ParseBash("/home/tester/.bash_history", strings.NewReader(""))
	if err != nil {
		t.Fatalf("ParseBash returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("len(entries) = %d, want 0", len(entries))
	}
	if len(warnings) != 0 {
		t.Fatalf("len(warnings) = %d, want 0", len(warnings))
	}
}

func TestParseBashRequiresSourceFile(t *testing.T) {
	if _, _, err := ParseBash("", strings.NewReader("pwd")); err == nil {
		t.Fatal("ParseBash returned nil error for empty source file")
	}
}

func TestParseBashRequiresReader(t *testing.T) {
	if _, _, err := ParseBash("/home/tester/.bash_history", nil); err == nil {
		t.Fatal("ParseBash returned nil error for nil reader")
	}
}
