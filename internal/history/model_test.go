package history

import (
	"testing"
	"time"
)

func TestHistoryEntryValidate(t *testing.T) {
	now := time.Unix(1712959000, 0).UTC()
	exitCode := 0

	entry := HistoryEntry{
		ID:         "entry-001",
		Shell:      ShellBash,
		SourceFile: "/home/tester/.bash_history",
		RawLine:    "git status",
		Command:    "git status",
		Timestamp:  &now,
		ExitCode:   &exitCode,
		SessionID:  "session-001",
		Hash:       "hash-001",
	}

	if err := entry.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestHistoryEntryValidateRequiresFields(t *testing.T) {
	tests := []struct {
		name  string
		entry HistoryEntry
	}{
		{
			name: "missing shell",
			entry: HistoryEntry{
				SourceFile: "/tmp/.bash_history",
				RawLine:    "pwd",
				Command:    "pwd",
			},
		},
		{
			name: "missing source file",
			entry: HistoryEntry{
				Shell:   ShellBash,
				RawLine: "pwd",
				Command: "pwd",
			},
		},
		{
			name: "missing raw line",
			entry: HistoryEntry{
				Shell:      ShellBash,
				SourceFile: "/tmp/.bash_history",
				Command:    "pwd",
			},
		},
		{
			name: "missing command",
			entry: HistoryEntry{
				Shell:      ShellBash,
				SourceFile: "/tmp/.bash_history",
				RawLine:    "pwd",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.entry.Validate(); err == nil {
				t.Fatal("Validate returned nil error")
			}
		})
	}
}

func TestHistoryEntryOptionalMetadata(t *testing.T) {
	entry := HistoryEntry{
		Shell:      ShellZsh,
		SourceFile: "/home/tester/.zsh_history",
		RawLine:    ": 1712959000:0;echo hi",
		Command:    "echo hi",
	}

	if entry.HasTimestamp() {
		t.Fatal("HasTimestamp = true, want false")
	}
	if entry.HasExitCode() {
		t.Fatal("HasExitCode = true, want false")
	}

	now := time.Unix(1712959000, 0).UTC()
	exitCode := 42
	entry.Timestamp = &now
	entry.ExitCode = &exitCode

	if !entry.HasTimestamp() {
		t.Fatal("HasTimestamp = false, want true")
	}
	if !entry.HasExitCode() {
		t.Fatal("HasExitCode = false, want true")
	}
}

func TestParseWarningValidate(t *testing.T) {
	warning := ParseWarning{
		Shell:      ShellZsh,
		SourceFile: "/home/tester/.zsh_history",
		LineNumber: 3,
		RawLine:    ": not-a-zsh-prefix",
		Message:    "malformed extended history prefix",
	}

	if err := warning.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestParseWarningValidateRequiresFields(t *testing.T) {
	tests := []struct {
		name    string
		warning ParseWarning
	}{
		{
			name: "missing shell",
			warning: ParseWarning{
				SourceFile: "/tmp/.zsh_history",
				LineNumber: 1,
				RawLine:    "bad",
				Message:    "bad line",
			},
		},
		{
			name: "missing source file",
			warning: ParseWarning{
				Shell:      ShellZsh,
				LineNumber: 1,
				RawLine:    "bad",
				Message:    "bad line",
			},
		},
		{
			name: "invalid line number",
			warning: ParseWarning{
				Shell:      ShellZsh,
				SourceFile: "/tmp/.zsh_history",
				LineNumber: 0,
				RawLine:    "bad",
				Message:    "bad line",
			},
		},
		{
			name: "missing raw line",
			warning: ParseWarning{
				Shell:      ShellZsh,
				SourceFile: "/tmp/.zsh_history",
				LineNumber: 1,
				Message:    "bad line",
			},
		},
		{
			name: "missing message",
			warning: ParseWarning{
				Shell:      ShellZsh,
				SourceFile: "/tmp/.zsh_history",
				LineNumber: 1,
				RawLine:    "bad",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.warning.Validate(); err == nil {
				t.Fatal("Validate returned nil error")
			}
		})
	}
}
