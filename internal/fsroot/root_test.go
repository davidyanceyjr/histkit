package fsroot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRequiresAbsoluteRoot(t *testing.T) {
	t.Parallel()

	if _, err := New(""); err == nil {
		t.Fatal("expected error for empty root")
	}
	if _, err := New("relative"); err == nil {
		t.Fatal("expected error for relative root")
	}
}

func TestResolveWithinRoot(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()
	root, err := New(rootPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	resolved, err := root.Resolve(filepath.Join("nested", "..", "history.db"))
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	want := filepath.Join(rootPath, "history.db")
	if resolved != want {
		t.Fatalf("Resolve() = %q, want %q", resolved, want)
	}
}

func TestResolveAllowsRootAndDescendants(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()
	root, err := New(rootPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	resolvedRoot, err := root.Resolve(rootPath)
	if err != nil {
		t.Fatalf("Resolve(root) error = %v", err)
	}
	if resolvedRoot != rootPath {
		t.Fatalf("Resolve(root) = %q, want %q", resolvedRoot, rootPath)
	}

	absoluteChild := filepath.Join(rootPath, "logs", "audit.log")
	resolvedChild, err := root.Resolve(absoluteChild)
	if err != nil {
		t.Fatalf("Resolve(child) error = %v", err)
	}
	if resolvedChild != absoluteChild {
		t.Fatalf("Resolve(child) = %q, want %q", resolvedChild, absoluteChild)
	}
}

func TestResolveRejectsEscape(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()
	root, err := New(rootPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if _, err := root.Resolve(filepath.Join("..", "outside")); err == nil {
		t.Fatal("expected relative escape to fail")
	}
	if _, err := root.Resolve(filepath.Join(filepath.Dir(rootPath), "outside")); err == nil {
		t.Fatal("expected absolute escape to fail")
	}
}

func TestOpenFileAndReadFileWithinRoot(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()
	root, err := New(rootPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	file, err := root.OpenFile("audit.log", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}
	if _, err := file.WriteString("created\n"); err != nil {
		_ = file.Close()
		t.Fatalf("WriteString() error = %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	data, err := root.ReadFile(filepath.Join(rootPath, "audit.log"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "created\n" {
		t.Fatalf("ReadFile() = %q, want %q", string(data), "created\n")
	}

	opened, err := root.Open("audit.log")
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := opened.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	info, err := root.Stat("audit.log")
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.Name() != "audit.log" {
		t.Fatalf("Stat().Name() = %q, want %q", info.Name(), "audit.log")
	}
}

func TestCreateTempWithinRoot(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()
	root, err := New(rootPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	file, err := root.CreateTemp(".", ".audit.tmp-*")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	name := file.Name()
	if err := file.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	relative, err := filepath.Rel(rootPath, name)
	if err != nil {
		t.Fatalf("Rel() error = %v", err)
	}
	if relative == ".." || filepath.IsAbs(relative) {
		t.Fatalf("CreateTemp() created file outside root: %q", name)
	}
}

func TestOpenFileRejectsEscape(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()
	root, err := New(rootPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if _, err := root.OpenFile(filepath.Join("..", "audit.log"), os.O_CREATE|os.O_WRONLY, 0o600); err == nil {
		t.Fatal("expected OpenFile() escape to fail")
	}
}

func TestCreateTempRejectsEscape(t *testing.T) {
	t.Parallel()

	rootPath := t.TempDir()
	root, err := New(rootPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if _, err := root.CreateTemp(filepath.Join("..", "outside"), ".audit.tmp-*"); err == nil {
		t.Fatal("expected CreateTemp() escape to fail")
	}
}
