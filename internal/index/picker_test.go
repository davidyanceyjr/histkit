package index

import (
	"testing"
	"time"

	"histkit/internal/history"
)

func TestQueryRecentHistoryEntriesOrdersNewestFirst(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	oldTimestamp := time.Date(2026, 4, 30, 10, 0, 0, 0, time.UTC)
	newTimestamp := time.Date(2026, 4, 30, 11, 0, 0, 0, time.UTC)

	if _, err := writeHistoryEntriesAt(db, []history.HistoryEntry{
		{
			ID:         "entry-old",
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "pwd",
			Command:    "pwd",
			Timestamp:  &oldTimestamp,
			Hash:       "hash-old",
		},
		{
			ID:         "entry-new",
			Shell:      history.ShellZsh,
			SourceFile: "/home/tester/.zsh_history",
			RawLine:    ": 1714474800:0;git status",
			Command:    "git status",
			Timestamp:  &newTimestamp,
			SessionID:  "session-002",
			Hash:       "hash-new",
		},
	}, time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("writeHistoryEntriesAt returned error: %v", err)
	}

	entries, err := QueryRecentHistoryEntries(db, 10)
	if err != nil {
		t.Fatalf("QueryRecentHistoryEntries returned error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
	}
	if entries[0].ID != "entry-new" || entries[1].ID != "entry-old" {
		t.Fatalf("entry order = [%q %q], want [entry-new entry-old]", entries[0].ID, entries[1].ID)
	}
	if entries[0].SessionID != "session-002" {
		t.Fatalf("entries[0].SessionID = %q, want session-002", entries[0].SessionID)
	}
}

func TestQueryRecentHistoryEntriesRequiresDBAndLimit(t *testing.T) {
	if _, err := QueryRecentHistoryEntries(nil, 10); err == nil {
		t.Fatal("QueryRecentHistoryEntries returned nil error for nil database")
	}

	db := openTestDB(t)
	defer db.Close()

	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	if _, err := QueryRecentHistoryEntries(db, 0); err == nil {
		t.Fatal("QueryRecentHistoryEntries returned nil error for non-positive limit")
	}
}
