package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.General.DefaultShell != "bash" {
		t.Fatalf("DefaultShell = %q, want bash", cfg.General.DefaultShell)
	}
	if !cfg.General.BackupHistory {
		t.Fatal("BackupHistory = false, want true")
	}
	if !cfg.General.DryRun {
		t.Fatal("DryRun = false, want true")
	}
	if !cfg.General.PreviewDiff {
		t.Fatal("PreviewDiff = false, want true")
	}
	if !cfg.Snippets.Enabled {
		t.Fatal("Snippets.Enabled = false, want true")
	}
	if !cfg.Snippets.Builtin {
		t.Fatal("Snippets.Builtin = false, want true")
	}
	if cfg.Snippets.UserFile != "~/.local/share/histkit/snippets.toml" {
		t.Fatalf("Snippets.UserFile = %q, want default path", cfg.Snippets.UserFile)
	}
}

func TestLoadDefaultsWhenPathEmpty(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg != Default() {
		t.Fatalf("Load(\"\") = %#v, want %#v", cfg, Default())
	}
}

func TestLoadFromPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	content := `
[general]
default_shell = "zsh"
backup_history = false
dry_run = true
preview_diff = false

[snippets]
enabled = false
builtin = false
user_file = "~/custom-snippets.toml"
`

	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.General.DefaultShell != "zsh" {
		t.Fatalf("DefaultShell = %q, want zsh", cfg.General.DefaultShell)
	}
	if cfg.General.BackupHistory {
		t.Fatal("BackupHistory = true, want false")
	}
	if cfg.General.DryRun != true {
		t.Fatal("DryRun = false, want true")
	}
	if cfg.General.PreviewDiff {
		t.Fatal("PreviewDiff = true, want false")
	}
	if cfg.Snippets.Enabled {
		t.Fatal("Snippets.Enabled = true, want false")
	}
	if cfg.Snippets.Builtin {
		t.Fatal("Snippets.Builtin = true, want false")
	}
	if cfg.Snippets.UserFile != "~/custom-snippets.toml" {
		t.Fatalf("Snippets.UserFile = %q, want custom path", cfg.Snippets.UserFile)
	}
}

func TestLoadExampleConfig(t *testing.T) {
	path := filepath.Join("..", "..", "configs", "config.example.toml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg != Default() {
		t.Fatalf("example config = %#v, want %#v", cfg, Default())
	}
}

func TestLoadRejectsInvalidTOML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	if err := os.WriteFile(path, []byte("[general\nbad"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("Load returned nil error for invalid TOML")
	}
}

func TestDefaultPaths(t *testing.T) {
	paths, err := DefaultPaths("/tmp/histkit-home")
	if err != nil {
		t.Fatalf("DefaultPaths returned error: %v", err)
	}

	if got, want := paths.ConfigFile, "/tmp/histkit-home/.config/histkit/config.toml"; got != want {
		t.Fatalf("ConfigFile = %q, want %q", got, want)
	}
	if got, want := paths.StateDir, "/tmp/histkit-home/.local/share/histkit"; got != want {
		t.Fatalf("StateDir = %q, want %q", got, want)
	}
	if got, want := paths.CacheDir, "/tmp/histkit-home/.cache/histkit"; got != want {
		t.Fatalf("CacheDir = %q, want %q", got, want)
	}
	if got, want := paths.HistoryDB, "/tmp/histkit-home/.local/share/histkit/history.db"; got != want {
		t.Fatalf("HistoryDB = %q, want %q", got, want)
	}
	if got, want := paths.SnippetsFile, "/tmp/histkit-home/.local/share/histkit/snippets.toml"; got != want {
		t.Fatalf("SnippetsFile = %q, want %q", got, want)
	}
}

func TestExpandUserPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "home only", in: "~", want: "/home/tester"},
		{name: "home child", in: "~/.config/histkit/config.toml", want: "/home/tester/.config/histkit/config.toml"},
		{name: "absolute untouched", in: "/var/tmp/config.toml", want: "/var/tmp/config.toml"},
		{name: "relative untouched", in: "configs/config.example.toml", want: "configs/config.example.toml"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExpandUserPath(tc.in, "/home/tester")
			if err != nil {
				t.Fatalf("ExpandUserPath returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("ExpandUserPath(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestExpandUserPathRequiresHome(t *testing.T) {
	if _, err := ExpandUserPath("~/.config/histkit/config.toml", ""); err == nil {
		t.Fatal("ExpandUserPath returned nil error for empty home directory")
	}
}
