package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRewriteAtomicReplacesContentsAndPreservesMode(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, ".bash_history")
	if err := os.WriteFile(targetPath, []byte("pwd\n"), 0o640); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if err := RewriteAtomic(targetPath, []byte("git status\n")); err != nil {
		t.Fatalf("RewriteAtomic returned error: %v", err)
	}

	data, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) != "git status\n" {
		t.Fatalf("rewritten contents = %q, want %q", string(data), "git status\n")
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatalf("Stat returned error: %v", err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("file mode = %o, want 640", info.Mode().Perm())
	}

	matches, err := filepath.Glob(filepath.Join(tempDir, ".*.tmp-*"))
	if err != nil {
		t.Fatalf("Glob returned error: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("leftover temp files = %v, want none", matches)
	}
}

func TestRewriteAtomicCreatesMissingFile(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, ".zsh_history")

	if err := RewriteAtomic(targetPath, []byte(": 1712959000:0;echo hello\n")); err != nil {
		t.Fatalf("RewriteAtomic returned error: %v", err)
	}

	data, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) != ": 1712959000:0;echo hello\n" {
		t.Fatalf("created contents = %q, want %q", string(data), ": 1712959000:0;echo hello\n")
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatalf("Stat returned error: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("file mode = %o, want 600", info.Mode().Perm())
	}
}

func TestRewriteAtomicRejectsEmptyPath(t *testing.T) {
	if err := RewriteAtomic("", []byte("pwd\n")); err == nil {
		t.Fatal("RewriteAtomic returned nil error for empty target path")
	}
}

func TestRewriteAtomicRejectsNonRegularTarget(t *testing.T) {
	tempDir := t.TempDir()

	if err := RewriteAtomic(tempDir, []byte("pwd\n")); err == nil {
		t.Fatal("RewriteAtomic returned nil error for directory target")
	}
}

func TestRewriteAtomicRejectsMissingParentDirectory(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "missing", ".bash_history")

	if err := RewriteAtomic(targetPath, []byte("pwd\n")); err == nil {
		t.Fatal("RewriteAtomic returned nil error for missing parent directory")
	}
}
