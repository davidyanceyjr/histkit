package sanitize

import (
	"strings"
	"testing"

	"histkit/internal/history"
)

func TestPreviewEntryReturnsMatchesWithoutMutatingEntry(t *testing.T) {
	entry := history.HistoryEntry{
		Shell:      history.ShellBash,
		SourceFile: "/home/tester/.bash_history",
		RawLine:    "mysql --password hunter2",
		Command:    "mysql --password hunter2",
	}

	item, err := PreviewEntry(entry)
	if err != nil {
		t.Fatalf("PreviewEntry returned error: %v", err)
	}
	if len(item.Matches) == 0 {
		t.Fatal("PreviewEntry returned no matches")
	}
	if entry.Command != "mysql --password hunter2" {
		t.Fatalf("entry.Command mutated to %q", entry.Command)
	}
	if item.Entry.Command != entry.Command {
		t.Fatalf("item.Entry.Command = %q, want %q", item.Entry.Command, entry.Command)
	}
}

func TestPreviewEntriesCountsActionsAcrossSecretAndTrivialRules(t *testing.T) {
	entries := []history.HistoryEntry{
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "clear",
			Command:    "clear",
		},
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "curl -H 'Authorization: Bearer abcdefghijklmnopqrstuvwx123456' https://example.com",
			Command:    "curl -H 'Authorization: Bearer abcdefghijklmnopqrstuvwx123456' https://example.com",
		},
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "echo hi",
			Command:    "echo hi",
		},
	}

	report, err := PreviewEntries(entries)
	if err != nil {
		t.Fatalf("PreviewEntries returned error: %v", err)
	}
	if report.TotalEntries != 3 {
		t.Fatalf("report.TotalEntries = %d, want 3", report.TotalEntries)
	}
	if report.MatchedEntries != 2 {
		t.Fatalf("report.MatchedEntries = %d, want 2", report.MatchedEntries)
	}
	if report.TotalMatches != 2 {
		t.Fatalf("report.TotalMatches = %d, want 2", report.TotalMatches)
	}
	if report.CountsByAction[ActionDelete] != 1 {
		t.Fatalf("delete count = %d, want 1", report.CountsByAction[ActionDelete])
	}
	if report.CountsByAction[ActionRedact] != 1 {
		t.Fatalf("redact count = %d, want 1", report.CountsByAction[ActionRedact])
	}
}

func TestRenderPreviewTextIncludesRequiredFields(t *testing.T) {
	report, err := PreviewEntries([]history.HistoryEntry{
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "mysql --password hunter2",
			Command:    "mysql --password hunter2",
		},
	})
	if err != nil {
		t.Fatalf("PreviewEntries returned error: %v", err)
	}

	text := RenderPreviewText(report)
	for _, want := range []string{
		"dry-run preview:",
		"counts by action:",
		"entry: shell=bash source=/home/tester/.bash_history",
		"original: mysql --password hunter2",
		"match: rule=inline-password-flag confidence=high action=redact",
		"transformed:",
		"reason: Redact inline password values",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("preview text missing %q in %q", want, text)
		}
	}
}

func TestRenderPreviewTextHandlesNoMatches(t *testing.T) {
	report, err := PreviewEntries([]history.HistoryEntry{
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "echo hi",
			Command:    "echo hi",
		},
	})
	if err != nil {
		t.Fatalf("PreviewEntries returned error: %v", err)
	}

	text := RenderPreviewText(report)
	if !strings.Contains(text, "no matches") {
		t.Fatalf("preview text = %q, want no matches message", text)
	}
}
