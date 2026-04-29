package history

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseZshFixture(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "history", "zsh", "extended.hist")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer f.Close()

	entries, warnings, err := ParseZsh("/home/tester/.zsh_history", f)
	if err != nil {
		t.Fatalf("ParseZsh returned error: %v", err)
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

	expectedTimestamps := []time.Time{
		time.Unix(1712959000, 0).UTC(),
		time.Unix(1712959015, 0).UTC(),
		time.Unix(1712959030, 0).UTC(),
		time.Unix(1712959045, 0).UTC(),
	}

	tests := []struct {
		index   int
		rawLine string
		command string
	}{
		{index: 0, rawLine: ": 1712959000:0;git status", command: "git status"},
		{index: 1, rawLine: ": 1712959015:2;echo 'hello;world'", command: "echo 'hello;world'"},
		{index: 2, rawLine: ": 1712959030:4;  ls -la", command: "  ls -la"},
		{index: 3, rawLine: ": 1712959045:6;printf '%s\\n' foo; printf '%s\\n' bar", command: "printf '%s\\n' foo; printf '%s\\n' bar"},
	}

	for _, tc := range tests {
		entry := entries[tc.index]
		if got, want := entry.Shell, ShellZsh; got != want {
			t.Fatalf("entries[%d].Shell = %q, want %q", tc.index, got, want)
		}
		if got, want := entry.SourceFile, "/home/tester/.zsh_history"; got != want {
			t.Fatalf("entries[%d].SourceFile = %q, want %q", tc.index, got, want)
		}
		if got, want := entry.RawLine, tc.rawLine; got != want {
			t.Fatalf("entries[%d].RawLine = %q, want %q", tc.index, got, want)
		}
		if got, want := entry.Command, tc.command; got != want {
			t.Fatalf("entries[%d].Command = %q, want %q", tc.index, got, want)
		}
		if !entry.HasTimestamp() {
			t.Fatalf("entries[%d].HasTimestamp = false, want true", tc.index)
		}
		if entry.Timestamp == nil || !entry.Timestamp.Equal(expectedTimestamps[tc.index]) {
			t.Fatalf("entries[%d].Timestamp = %v, want %v", tc.index, entry.Timestamp, expectedTimestamps[tc.index])
		}
		if entry.HasExitCode() {
			t.Fatalf("entries[%d].HasExitCode = true, want false", tc.index)
		}
	}
}

func TestParseZshMalformedPrefixWarning(t *testing.T) {
	entries, warnings, err := ParseZsh("/home/tester/.zsh_history", strings.NewReader(": not-a-zsh-prefix"))
	if err != nil {
		t.Fatalf("ParseZsh returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("len(entries) = %d, want 0", len(entries))
	}
	if len(warnings) != 1 {
		t.Fatalf("len(warnings) = %d, want 1", len(warnings))
	}
	if !strings.Contains(warnings[0].Message, "missing command separator") &&
		!strings.Contains(warnings[0].Message, "invalid metadata fields") &&
		!strings.Contains(warnings[0].Message, "invalid timestamp") {
		t.Fatalf("unexpected warning message: %q", warnings[0].Message)
	}
}

func TestParseZshEmptyInput(t *testing.T) {
	entries, warnings, err := ParseZsh("/home/tester/.zsh_history", strings.NewReader(""))
	if err != nil {
		t.Fatalf("ParseZsh returned error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("len(entries) = %d, want 0", len(entries))
	}
	if len(warnings) != 0 {
		t.Fatalf("len(warnings) = %d, want 0", len(warnings))
	}
}

func TestParseZshRequiresSourceFile(t *testing.T) {
	if _, _, err := ParseZsh("", strings.NewReader(": 1712959000:0;pwd")); err == nil {
		t.Fatal("ParseZsh returned nil error for empty source file")
	}
}

func TestParseZshRequiresReader(t *testing.T) {
	if _, _, err := ParseZsh("/home/tester/.zsh_history", nil); err == nil {
		t.Fatal("ParseZsh returned nil error for nil reader")
	}
}
