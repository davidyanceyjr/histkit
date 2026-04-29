package index

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"histkit/internal/history"
)

type WriteResult struct {
	Attempted int
	Inserted  int
	Skipped   int
}

func WriteHistoryEntries(db *sql.DB, entries []history.HistoryEntry) (WriteResult, error) {
	return writeHistoryEntriesAt(db, entries, time.Now().UTC())
}

func writeHistoryEntriesAt(db *sql.DB, entries []history.HistoryEntry, ingestedAt time.Time) (WriteResult, error) {
	if db == nil {
		return WriteResult{}, fmt.Errorf("write history entries: database is required")
	}
	if ingestedAt.IsZero() {
		return WriteResult{}, fmt.Errorf("write history entries: ingested time is required")
	}
	if len(entries) == 0 {
		return WriteResult{}, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return WriteResult{}, fmt.Errorf("write history entries: %w", err)
	}

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO history_entries
			(id, shell, source_file, raw_line, command, timestamp, exit_code, session_id, hash, ingested_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		_ = tx.Rollback()
		return WriteResult{}, fmt.Errorf("write history entries: %w", err)
	}
	defer stmt.Close()

	result := WriteResult{Attempted: len(entries)}
	ingestedValue := ingestedAt.UTC().Format(time.RFC3339Nano)

	for _, entry := range entries {
		stored, err := prepareEntry(entry)
		if err != nil {
			_ = tx.Rollback()
			return WriteResult{}, fmt.Errorf("write history entries: %w", err)
		}

		execResult, err := stmt.Exec(
			stored.ID,
			stored.Shell,
			stored.SourceFile,
			stored.RawLine,
			stored.Command,
			formatTimestamp(stored.Timestamp),
			nullableInt(stored.ExitCode),
			nullableString(stored.SessionID),
			nullableString(stored.Hash),
			ingestedValue,
		)
		if err != nil {
			_ = tx.Rollback()
			return WriteResult{}, fmt.Errorf("write history entries: insert %q: %w", stored.SourceFile, err)
		}

		rowsAffected, err := execResult.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return WriteResult{}, fmt.Errorf("write history entries: %w", err)
		}

		if rowsAffected == 0 {
			result.Skipped++
			continue
		}
		result.Inserted++
	}

	if err := tx.Commit(); err != nil {
		return WriteResult{}, fmt.Errorf("write history entries: %w", err)
	}

	return result, nil
}

func prepareEntry(entry history.HistoryEntry) (history.HistoryEntry, error) {
	if err := entry.Validate(); err != nil {
		return history.HistoryEntry{}, err
	}

	if strings.TrimSpace(entry.Hash) == "" {
		entry.Hash = hashCommand(entry.Command)
	}
	if strings.TrimSpace(entry.ID) == "" {
		entry.ID = deriveEntryID(entry)
	}

	return entry, nil
}

func hashCommand(command string) string {
	sum := sha256.Sum256([]byte(command))
	return hex.EncodeToString(sum[:])
}

func deriveEntryID(entry history.HistoryEntry) string {
	h := sha256.New()
	writeIDPart(h, entry.Shell)
	writeIDPart(h, entry.SourceFile)
	writeIDPart(h, entry.RawLine)
	writeIDPart(h, entry.Command)
	if entry.Timestamp != nil {
		writeIDPart(h, entry.Timestamp.UTC().Format(time.RFC3339Nano))
	} else {
		writeIDPart(h, "")
	}
	if entry.ExitCode != nil {
		writeIDPart(h, strconv.Itoa(*entry.ExitCode))
	} else {
		writeIDPart(h, "")
	}
	writeIDPart(h, entry.SessionID)
	writeIDPart(h, entry.Hash)

	return "entry-" + hex.EncodeToString(h.Sum(nil))
}

func writeIDPart(h interface{ Write([]byte) (int, error) }, value string) {
	_, _ = h.Write([]byte(value))
	_, _ = h.Write([]byte{0})
}

func formatTimestamp(ts *time.Time) any {
	if ts == nil {
		return nil
	}

	return ts.UTC().Format(time.RFC3339Nano)
}

func nullableInt(value *int) any {
	if value == nil {
		return nil
	}

	return *value
}

func nullableString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	return value
}
