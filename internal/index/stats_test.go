package index

import "testing"

func TestQueryHistoryStatsEmptyDatabase(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	stats, err := QueryHistoryStats(db)
	if err != nil {
		t.Fatalf("QueryHistoryStats returned error: %v", err)
	}

	if stats.TotalEntries != 0 {
		t.Fatalf("TotalEntries = %d, want 0", stats.TotalEntries)
	}
	if len(stats.ByShell) != 0 {
		t.Fatalf("ByShell length = %d, want 0", len(stats.ByShell))
	}
	if len(stats.BySource) != 0 {
		t.Fatalf("BySource length = %d, want 0", len(stats.BySource))
	}
}

func TestQueryHistoryStatsGroupedCounts(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	if err := InitSchema(db); err != nil {
		t.Fatalf("InitSchema returned error: %v", err)
	}

	insert := `INSERT INTO history_entries
		(id, shell, source_file, raw_line, command, timestamp, exit_code, session_id, hash, ingested_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	fixtures := [][]any{
		{"entry-001", "bash", "/home/tester/.bash_history", "pwd", "pwd", nil, nil, nil, "hash-001", "2026-04-29T18:00:00Z"},
		{"entry-002", "bash", "/home/tester/.bash_history", "git status", "git status", nil, nil, nil, "hash-002", "2026-04-29T18:00:01Z"},
		{"entry-003", "zsh", "/home/tester/.zsh_history", ": 1712959000:0;echo hi", "echo hi", "2024-04-13T00:36:40Z", nil, nil, "hash-003", "2026-04-29T18:00:02Z"},
	}

	for _, fixture := range fixtures {
		if _, err := db.Exec(insert, fixture...); err != nil {
			t.Fatalf("insert returned error: %v", err)
		}
	}

	stats, err := QueryHistoryStats(db)
	if err != nil {
		t.Fatalf("QueryHistoryStats returned error: %v", err)
	}

	if stats.TotalEntries != 3 {
		t.Fatalf("TotalEntries = %d, want 3", stats.TotalEntries)
	}
	if len(stats.ByShell) != 2 {
		t.Fatalf("ByShell length = %d, want 2", len(stats.ByShell))
	}
	if stats.ByShell[0] != (GroupCount{Name: "bash", Count: 2}) {
		t.Fatalf("ByShell[0] = %#v, want bash=2", stats.ByShell[0])
	}
	if stats.ByShell[1] != (GroupCount{Name: "zsh", Count: 1}) {
		t.Fatalf("ByShell[1] = %#v, want zsh=1", stats.ByShell[1])
	}
	if len(stats.BySource) != 2 {
		t.Fatalf("BySource length = %d, want 2", len(stats.BySource))
	}
	if stats.BySource[0] != (GroupCount{Name: "/home/tester/.bash_history", Count: 2}) {
		t.Fatalf("BySource[0] = %#v, want bash history=2", stats.BySource[0])
	}
	if stats.BySource[1] != (GroupCount{Name: "/home/tester/.zsh_history", Count: 1}) {
		t.Fatalf("BySource[1] = %#v, want zsh history=1", stats.BySource[1])
	}
}

func TestQueryHistoryStatsRequiresDB(t *testing.T) {
	if _, err := QueryHistoryStats(nil); err == nil {
		t.Fatal("QueryHistoryStats returned nil error for nil database")
	}
}
