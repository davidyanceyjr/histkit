package index

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestOpenRequiresPath(t *testing.T) {
	if _, err := Open(""); err == nil {
		t.Fatal("Open returned nil error for empty path")
	}
}

func TestInitSchemaCreatesTables(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	assertTableExists(t, db, "history_entries")
	assertTableExists(t, db, "runs")

	var version int
	if err := db.QueryRow(`PRAGMA user_version;`).Scan(&version); err != nil {
		t.Fatalf("QueryRow(PRAGMA user_version) returned error: %v", err)
	}
	if version != SchemaVersion {
		t.Fatalf("user_version = %d, want %d", version, SchemaVersion)
	}
}

func TestHistoryEntriesDedupeBySourceAndHash(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	insert := `INSERT INTO history_entries
		(id, shell, source_file, raw_line, command, timestamp, exit_code, session_id, hash, ingested_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	if _, err := db.Exec(insert,
		"entry-001", "bash", "/home/tester/.bash_history", "git status", "git status",
		nil, nil, "", "hash-001", "2026-04-29T18:00:00Z",
	); err != nil {
		t.Fatalf("first insert returned error: %v", err)
	}

	if _, err := db.Exec(insert,
		"entry-002", "bash", "/home/tester/.bash_history", "git status", "git status",
		nil, nil, "", "hash-001", "2026-04-29T18:01:00Z",
	); err == nil {
		t.Fatal("expected duplicate hash/source insert to fail")
	}

	if _, err := db.Exec(insert,
		"entry-003", "bash", "/home/tester/.zsh_history", "git status", "git status",
		nil, nil, "", "hash-001", "2026-04-29T18:02:00Z",
	); err != nil {
		t.Fatalf("different source insert returned error: %v", err)
	}
}

func TestRunsInsertAndReopenDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "histkit.db")

	db, err := Open(path)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}

	if err := InitSchema(db); err != nil {
		_ = db.Close()
		t.Fatalf("InitSchema returned error: %v", err)
	}

	if _, err := db.Exec(
		`INSERT INTO runs (id, command, started_at, finished_at, status, notes)
		 VALUES (?, ?, ?, ?, ?, ?);`,
		"run-001", "scan", "2026-04-29T18:00:00Z", "2026-04-29T18:00:02Z", "ok", "fixture run",
	); err != nil {
		_ = db.Close()
		t.Fatalf("insert run returned error: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen returned error: %v", err)
	}
	defer reopened.Close()

	var count int
	if err := reopened.QueryRow(`SELECT COUNT(*) FROM runs;`).Scan(&count); err != nil {
		t.Fatalf("QueryRow(count) returned error: %v", err)
	}
	if count != 1 {
		t.Fatalf("run count = %d, want 1", count)
	}
}

func TestInitSchemaRequiresDB(t *testing.T) {
	if err := InitSchema(nil); err == nil {
		t.Fatal("InitSchema returned nil error for nil database")
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	path := filepath.Join(t.TempDir(), "histkit.db")
	db, err := Open(path)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}

	return db
}

func assertTableExists(t *testing.T, db *sql.DB, name string) {
	t.Helper()

	var found string
	if err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?;`, name).Scan(&found); err != nil {
		t.Fatalf("table lookup for %q returned error: %v", name, err)
	}
	if found != name {
		t.Fatalf("found table = %q, want %q", found, name)
	}
}
