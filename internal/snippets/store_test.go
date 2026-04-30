package snippets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoreListMissingFileReturnsEmpty(t *testing.T) {
	store := Store{Path: filepath.Join(t.TempDir(), "snippets.toml")}

	snippets, err := store.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(snippets) != 0 {
		t.Fatalf("len(snippets) = %d, want 0", len(snippets))
	}
}

func TestStoreSaveAndListRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snippets.toml")
	store := Store{Path: path}

	snippets := []Snippet{
		{
			ID:          "find-delete-pyc",
			Title:       "Delete Python cache files",
			Command:     "find {{path}} -type f -name '*.pyc' -delete",
			Description: "Delete .pyc files under a path",
			Tags:        []string{"find", "python", "cleanup"},
			Shells:      []string{"bash", "zsh"},
			Safety:      SafetyMedium,
		},
	}

	if err := store.Save(snippets); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded, err := store.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("len(loaded) = %d, want 1", len(loaded))
	}
	if loaded[0].Command != snippets[0].Command {
		t.Fatalf("Command = %q, want exact template preservation", loaded[0].Command)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if !strings.Contains(string(content), "[[snippets]]") {
		t.Fatalf("snippet file content missing [[snippets]]: %q", string(content))
	}
}

func TestStoreListRejectsDuplicateIDs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snippets.toml")
	content := `
[[snippets]]
id = "dup-id"
title = "One"
command = "echo one"
description = "first"
safety = "low"

[[snippets]]
id = "dup-id"
title = "Two"
command = "echo two"
description = "second"
safety = "medium"
`
	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	store := Store{Path: path}
	if _, err := store.List(); err == nil {
		t.Fatal("List returned nil error for duplicate ids")
	}
}

func TestStoreAddAndRemove(t *testing.T) {
	store := Store{Path: filepath.Join(t.TempDir(), "snippets.toml")}

	first := Snippet{
		ID:          "find-delete-pyc",
		Title:       "Delete Python cache files",
		Command:     "find {{path}} -type f -name '*.pyc' -delete",
		Description: "Delete .pyc files under a path",
		Safety:      SafetyMedium,
	}
	second := Snippet{
		ID:          "git-clean-branches",
		Title:       "Delete merged branches",
		Command:     "git branch --merged | grep -v '\\*\\|main\\|master' | xargs -r git branch -d",
		Description: "Delete local branches already merged",
		Safety:      SafetyHigh,
	}

	if err := store.Add(first); err != nil {
		t.Fatalf("Add(first) returned error: %v", err)
	}
	if err := store.Add(second); err != nil {
		t.Fatalf("Add(second) returned error: %v", err)
	}
	if err := store.Remove(first.ID); err != nil {
		t.Fatalf("Remove returned error: %v", err)
	}

	loaded, err := store.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(loaded) != 1 {
		t.Fatalf("len(loaded) = %d, want 1", len(loaded))
	}
	if loaded[0].ID != second.ID {
		t.Fatalf("remaining snippet id = %q, want %q", loaded[0].ID, second.ID)
	}
}

func TestStoreRemoveMissingIDFails(t *testing.T) {
	store := Store{Path: filepath.Join(t.TempDir(), "snippets.toml")}

	if err := store.Remove("missing-id"); err == nil {
		t.Fatal("Remove returned nil error for missing snippet id")
	}
}
