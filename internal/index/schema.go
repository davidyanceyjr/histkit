package index

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const SchemaVersion = 1

func Open(path string) (*sql.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("open sqlite database: path is required")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("open sqlite database %q: %w", path, err)
	}

	file, err := os.OpenFile(path, os.O_CREATE, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database %q: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close sqlite database bootstrap file %q: %w", path, err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database %q: %w", path, err)
	}

	return db, nil
}

func InitSchema(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("init schema: database is required")
	}

	stmts := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS history_entries (
			id TEXT PRIMARY KEY,
			shell TEXT NOT NULL,
			source_file TEXT NOT NULL,
			raw_line TEXT NOT NULL,
			command TEXT NOT NULL,
			timestamp TEXT,
			exit_code INTEGER,
			session_id TEXT,
			hash TEXT,
			ingested_at TEXT NOT NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_history_entries_source_hash
			ON history_entries(source_file, hash)
			WHERE hash IS NOT NULL AND hash <> '';`,
		`CREATE INDEX IF NOT EXISTS idx_history_entries_shell
			ON history_entries(shell);`,
		`CREATE INDEX IF NOT EXISTS idx_history_entries_timestamp
			ON history_entries(timestamp);`,
		`CREATE TABLE IF NOT EXISTS runs (
			id TEXT PRIMARY KEY,
			command TEXT NOT NULL,
			started_at TEXT NOT NULL,
			finished_at TEXT,
			status TEXT NOT NULL,
			notes TEXT
		);`,
		fmt.Sprintf(`PRAGMA user_version = %d;`, SchemaVersion),
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("init schema: %w", err)
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("init schema: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("init schema: %w", err)
	}

	return nil
}
