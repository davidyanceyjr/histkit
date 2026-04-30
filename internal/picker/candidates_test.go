package picker

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"histkit/internal/history"
	"histkit/internal/index"
	"histkit/internal/snippets"
)

func TestLoadCandidatesMergesHistoryAndSnippets(t *testing.T) {
	db := openPickerDB(t)
	defer db.Close()

	timestamp := time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC)
	if _, err := index.WriteHistoryEntries(db, []history.HistoryEntry{
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "git status",
			Command:    "git status",
			Timestamp:  &timestamp,
		},
	}); err != nil {
		t.Fatalf("WriteHistoryEntries returned error: %v", err)
	}

	store := snippets.Store{Path: filepath.Join(t.TempDir(), "snippets.toml")}
	if err := store.Save([]snippets.Snippet{
		{
			ID:          "find-delete-pyc",
			Title:       "Delete Python cache files",
			Command:     "find {{path}} -type f -name '*.pyc' -delete",
			Description: "Delete .pyc files under a path",
			Safety:      snippets.SafetyMedium,
		},
	}); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	candidates, err := LoadCandidates(db, store, true, false, 20)
	if err != nil {
		t.Fatalf("LoadCandidates returned error: %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("len(candidates) = %d, want 2", len(candidates))
	}
	if candidates[0].Label != LabelHistory || candidates[0].Display() != "[history]  git status" {
		t.Fatalf("history candidate = %#v, want labeled history display", candidates[0])
	}
	if candidates[1].Label != LabelSnippet || candidates[1].Display() != "[snippet]  find {{path}} -type f -name '*.pyc' -delete" {
		t.Fatalf("snippet candidate = %#v, want labeled snippet display", candidates[1])
	}
}

func TestLoadCandidatesIncludesMissingBuiltinsWithoutOverwritingUserSnippets(t *testing.T) {
	db := openPickerDB(t)
	defer db.Close()

	store := snippets.Store{Path: filepath.Join(t.TempDir(), "snippets.toml")}
	if err := store.Save([]snippets.Snippet{
		{
			ID:          "find-delete-pyc",
			Title:       "Custom cleanup",
			Command:     "echo custom",
			Description: "User snippet should win over builtin",
			Safety:      snippets.SafetyLow,
		},
	}); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	candidates, err := LoadCandidates(db, store, true, true, 20)
	if err != nil {
		t.Fatalf("LoadCandidates returned error: %v", err)
	}

	var userFound bool
	var builtinCount int
	for _, candidate := range candidates {
		if candidate.Label != LabelSnippet {
			continue
		}
		builtinCount++
		if candidate.SnippetID == "find-delete-pyc" {
			userFound = true
			if candidate.Command != "echo custom" {
				t.Fatalf("candidate command = %q, want user snippet preserved", candidate.Command)
			}
		}
	}
	if !userFound {
		t.Fatal("expected user snippet candidate to be present")
	}
	if builtinCount != len(snippets.Builtins()) {
		t.Fatalf("snippet candidate count = %d, want %d", builtinCount, len(snippets.Builtins()))
	}
}

func TestParseSelectedLine(t *testing.T) {
	historyCandidate, err := ParseSelectedLine("[history]  git status")
	if err != nil {
		t.Fatalf("ParseSelectedLine(history) returned error: %v", err)
	}
	if historyCandidate.Label != LabelHistory || historyCandidate.Command != "git status" {
		t.Fatalf("historyCandidate = %#v, want label history and command git status", historyCandidate)
	}

	snippetCandidate, err := ParseSelectedLine("[snippet]  find {{path}}")
	if err != nil {
		t.Fatalf("ParseSelectedLine(snippet) returned error: %v", err)
	}
	if snippetCandidate.Label != LabelSnippet || snippetCandidate.Command != "find {{path}}" {
		t.Fatalf("snippetCandidate = %#v, want label snippet and command template", snippetCandidate)
	}

	if _, err := ParseSelectedLine("plain text"); err == nil {
		t.Fatal("ParseSelectedLine returned nil error for unsupported format")
	}
}

func TestLoadCandidatesSkipsSnippetsWhenDisabled(t *testing.T) {
	db := openPickerDB(t)
	defer db.Close()

	store := snippets.Store{Path: filepath.Join(t.TempDir(), "snippets.toml")}
	if err := store.Save([]snippets.Snippet{
		{
			ID:          "echo-test",
			Title:       "Echo test",
			Command:     "echo test",
			Description: "Emit a test string",
			Safety:      snippets.SafetyLow,
		},
	}); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	candidates, err := LoadCandidates(db, store, false, true, 20)
	if err != nil {
		t.Fatalf("LoadCandidates returned error: %v", err)
	}
	if len(candidates) != 0 {
		t.Fatalf("len(candidates) = %d, want 0 when snippets disabled and no history", len(candidates))
	}
}

func openPickerDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := index.Open(filepath.Join(t.TempDir(), "histkit.db"))
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	if err := index.InitSchema(db); err != nil {
		db.Close()
		t.Fatalf("InitSchema returned error: %v", err)
	}

	return db
}
