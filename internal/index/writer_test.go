package index

import (
	"testing"
	"time"

	"histkit/internal/history"
)

func TestWriteHistoryEntriesDerivesFieldsAndStoresMetadata(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	timestamp := time.Date(2026, 4, 29, 18, 30, 0, 0, time.UTC)
	exitCode := 17
	ingestedAt := time.Date(2026, 4, 29, 18, 31, 0, 0, time.UTC)
	entry := history.HistoryEntry{
		Shell:      history.ShellZsh,
		SourceFile: "/home/tester/.zsh_history",
		RawLine:    ": 1714415400:0;kubectl get pods",
		Command:    "kubectl get pods",
		Timestamp:  &timestamp,
		ExitCode:   &exitCode,
		SessionID:  "session-123",
	}

	result, err := writeHistoryEntriesAt(db, []history.HistoryEntry{entry}, ingestedAt)
	if err != nil {
		t.Fatalf("writeHistoryEntriesAt returned error: %v", err)
	}
	if result.Attempted != 1 || result.Inserted != 1 || result.Skipped != 0 {
		t.Fatalf("result = %+v, want attempted=1 inserted=1 skipped=0", result)
	}

	expectedHash := hashCommand(entry.Command)
	expectedID := deriveEntryID(history.HistoryEntry{
		ID:         "",
		Shell:      entry.Shell,
		SourceFile: entry.SourceFile,
		RawLine:    entry.RawLine,
		Command:    entry.Command,
		Timestamp:  entry.Timestamp,
		ExitCode:   entry.ExitCode,
		SessionID:  entry.SessionID,
		Hash:       expectedHash,
	})

	var gotID string
	var gotHash string
	var gotTimestamp string
	var gotExitCode int
	var gotSessionID string
	var gotIngestedAt string
	if err := db.QueryRow(`
		SELECT id, hash, timestamp, exit_code, session_id, ingested_at
		FROM history_entries
		WHERE source_file = ? AND command = ?;
	`, entry.SourceFile, entry.Command).Scan(
		&gotID, &gotHash, &gotTimestamp, &gotExitCode, &gotSessionID, &gotIngestedAt,
	); err != nil {
		t.Fatalf("QueryRow(history_entries) returned error: %v", err)
	}

	if gotID != expectedID {
		t.Fatalf("stored id = %q, want %q", gotID, expectedID)
	}
	if gotHash != expectedHash {
		t.Fatalf("stored hash = %q, want %q", gotHash, expectedHash)
	}
	if gotTimestamp != timestamp.Format(time.RFC3339Nano) {
		t.Fatalf("stored timestamp = %q, want %q", gotTimestamp, timestamp.Format(time.RFC3339Nano))
	}
	if gotExitCode != exitCode {
		t.Fatalf("stored exit code = %d, want %d", gotExitCode, exitCode)
	}
	if gotSessionID != entry.SessionID {
		t.Fatalf("stored session id = %q, want %q", gotSessionID, entry.SessionID)
	}
	if gotIngestedAt != ingestedAt.Format(time.RFC3339Nano) {
		t.Fatalf("stored ingested_at = %q, want %q", gotIngestedAt, ingestedAt.Format(time.RFC3339Nano))
	}
}

func TestWriteHistoryEntriesSkipsDuplicateSourceAndHash(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	ingestedAt := time.Date(2026, 4, 29, 18, 45, 0, 0, time.UTC)
	entries := []history.HistoryEntry{
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "git status",
			Command:    "git status",
		},
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "git status",
			Command:    "git status",
		},
	}

	result, err := writeHistoryEntriesAt(db, entries, ingestedAt)
	if err != nil {
		t.Fatalf("writeHistoryEntriesAt returned error: %v", err)
	}
	if result.Attempted != 2 || result.Inserted != 1 || result.Skipped != 1 {
		t.Fatalf("result = %+v, want attempted=2 inserted=1 skipped=1", result)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM history_entries;`).Scan(&count); err != nil {
		t.Fatalf("QueryRow(count) returned error: %v", err)
	}
	if count != 1 {
		t.Fatalf("row count = %d, want 1", count)
	}
}

func TestWriteHistoryEntriesRollsBackOnInvalidEntry(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	ingestedAt := time.Date(2026, 4, 29, 19, 0, 0, 0, time.UTC)
	entries := []history.HistoryEntry{
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "pwd",
			Command:    "pwd",
		},
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "   ",
			Command:    "",
		},
	}

	if _, err := writeHistoryEntriesAt(db, entries, ingestedAt); err == nil {
		t.Fatal("writeHistoryEntriesAt returned nil error for invalid entry")
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM history_entries;`).Scan(&count); err != nil {
		t.Fatalf("QueryRow(count) returned error: %v", err)
	}
	if count != 0 {
		t.Fatalf("row count after rollback = %d, want 0", count)
	}
}

func TestWriteHistoryEntriesRequiresDB(t *testing.T) {
	entry := history.HistoryEntry{
		Shell:      history.ShellBash,
		SourceFile: "/home/tester/.bash_history",
		RawLine:    "ls",
		Command:    "ls",
	}

	if _, err := writeHistoryEntriesAt(nil, []history.HistoryEntry{entry}, time.Now().UTC()); err == nil {
		t.Fatal("writeHistoryEntriesAt returned nil error for nil database")
	}
}
