package snippets

import (
	"path/filepath"
	"testing"
)

func TestBuiltinsValidate(t *testing.T) {
	builtins := Builtins()
	if len(builtins) == 0 {
		t.Fatal("Builtins returned no snippets")
	}

	if err := ValidateCollection(builtins); err != nil {
		t.Fatalf("ValidateCollection returned error: %v", err)
	}
}

func TestImportBuiltinsSeedsMissingSnippets(t *testing.T) {
	store := Store{Path: filepath.Join(t.TempDir(), "snippets.toml")}

	imported, err := ImportBuiltins(store)
	if err != nil {
		t.Fatalf("ImportBuiltins returned error: %v", err)
	}
	if imported != len(Builtins()) {
		t.Fatalf("imported = %d, want %d", imported, len(Builtins()))
	}

	loaded, err := store.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(loaded) != len(Builtins()) {
		t.Fatalf("len(loaded) = %d, want %d", len(loaded), len(Builtins()))
	}
}

func TestImportBuiltinsDoesNotOverwriteExistingSnippet(t *testing.T) {
	store := Store{Path: filepath.Join(t.TempDir(), "snippets.toml")}

	custom := Snippet{
		ID:          "find-delete-pyc",
		Title:       "Custom Python cleanup",
		Command:     "echo custom",
		Description: "A user override-like entry that should remain unchanged",
		Safety:      SafetyLow,
	}
	if err := store.Save([]Snippet{custom}); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	imported, err := ImportBuiltins(store)
	if err != nil {
		t.Fatalf("ImportBuiltins returned error: %v", err)
	}
	if imported != len(Builtins())-1 {
		t.Fatalf("imported = %d, want %d", imported, len(Builtins())-1)
	}

	loaded, err := store.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	var found Snippet
	ok := false
	for _, snippet := range loaded {
		if snippet.ID == custom.ID {
			found = snippet
			ok = true
			break
		}
	}
	if !ok {
		t.Fatalf("expected snippet %q to remain present", custom.ID)
	}
	if found.Command != custom.Command {
		t.Fatalf("Command = %q, want %q", found.Command, custom.Command)
	}
}

func TestImportBuiltinsIdempotent(t *testing.T) {
	store := Store{Path: filepath.Join(t.TempDir(), "snippets.toml")}

	firstImported, err := ImportBuiltins(store)
	if err != nil {
		t.Fatalf("first ImportBuiltins returned error: %v", err)
	}
	if firstImported == 0 {
		t.Fatal("first ImportBuiltins imported 0 snippets, want > 0")
	}

	secondImported, err := ImportBuiltins(store)
	if err != nil {
		t.Fatalf("second ImportBuiltins returned error: %v", err)
	}
	if secondImported != 0 {
		t.Fatalf("second imported = %d, want 0", secondImported)
	}
}
