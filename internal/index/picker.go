package index

import (
	"database/sql"
	"fmt"
	"time"

	"histkit/internal/history"
)

func QueryRecentHistoryEntries(db *sql.DB, limit int) ([]history.HistoryEntry, error) {
	if db == nil {
		return nil, fmt.Errorf("query recent history entries: database is required")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("query recent history entries: limit must be positive")
	}

	rows, err := db.Query(`
		SELECT id, shell, source_file, raw_line, command, timestamp, exit_code, session_id, hash
		FROM history_entries
		ORDER BY COALESCE(timestamp, ingested_at) DESC, ingested_at DESC
		LIMIT ?;
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent history entries: %w", err)
	}
	defer rows.Close()

	var entries []history.HistoryEntry
	for rows.Next() {
		var entry history.HistoryEntry
		var timestamp sql.NullString
		var exitCode sql.NullInt64
		var sessionID sql.NullString
		var hash sql.NullString

		if err := rows.Scan(
			&entry.ID,
			&entry.Shell,
			&entry.SourceFile,
			&entry.RawLine,
			&entry.Command,
			&timestamp,
			&exitCode,
			&sessionID,
			&hash,
		); err != nil {
			return nil, fmt.Errorf("query recent history entries: %w", err)
		}

		if timestamp.Valid {
			parsed, err := time.Parse(time.RFC3339Nano, timestamp.String)
			if err != nil {
				return nil, fmt.Errorf("query recent history entries: parse timestamp %q: %w", timestamp.String, err)
			}
			entry.Timestamp = &parsed
		}
		if exitCode.Valid {
			value := int(exitCode.Int64)
			entry.ExitCode = &value
		}
		if sessionID.Valid {
			entry.SessionID = sessionID.String
		}
		if hash.Valid {
			entry.Hash = hash.String
		}

		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query recent history entries: %w", err)
	}

	return entries, nil
}
