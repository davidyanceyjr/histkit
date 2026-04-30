package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	General  General  `toml:"general"`
	Snippets Snippets `toml:"snippets"`
}

type General struct {
	DefaultShell  string `toml:"default_shell"`
	BackupHistory bool   `toml:"backup_history"`
	DryRun        bool   `toml:"dry_run"`
	PreviewDiff   bool   `toml:"preview_diff"`
}

type Snippets struct {
	Enabled  bool   `toml:"enabled"`
	Builtin  bool   `toml:"builtin"`
	UserFile string `toml:"user_file"`
}

type Paths struct {
	ConfigFile   string
	StateDir     string
	CacheDir     string
	HistoryDB    string
	SnippetsFile string
}

func Default() Config {
	return Config{
		General: General{
			DefaultShell:  "bash",
			BackupHistory: true,
			DryRun:        true,
			PreviewDiff:   true,
		},
		Snippets: Snippets{
			Enabled:  true,
			Builtin:  true,
			UserFile: "~/.local/share/histkit/snippets.toml",
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("load config %q: %w", path, err)
	}

	if cfg.General.DefaultShell == "" {
		cfg.General.DefaultShell = Default().General.DefaultShell
	}
	if cfg.Snippets.UserFile == "" {
		cfg.Snippets.UserFile = Default().Snippets.UserFile
	}

	return cfg, nil
}

func DefaultPaths(home string) (Paths, error) {
	if strings.TrimSpace(home) == "" {
		return Paths{}, fmt.Errorf("default paths: home directory is required")
	}

	return Paths{
		ConfigFile:   filepath.Join(home, ".config", "histkit", "config.toml"),
		StateDir:     filepath.Join(home, ".local", "share", "histkit"),
		CacheDir:     filepath.Join(home, ".cache", "histkit"),
		HistoryDB:    filepath.Join(home, ".local", "share", "histkit", "history.db"),
		SnippetsFile: filepath.Join(home, ".local", "share", "histkit", "snippets.toml"),
	}, nil
}

func DetectDefaultPaths() (Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Paths{}, fmt.Errorf("detect default paths: %w", err)
	}

	return DefaultPaths(home)
}

func ExpandUserPath(path, home string) (string, error) {
	if path == "" {
		return "", nil
	}
	if strings.TrimSpace(home) == "" {
		return "", fmt.Errorf("expand user path %q: home directory is required", path)
	}

	switch {
	case path == "~":
		return home, nil
	case strings.HasPrefix(path, "~/"):
		return filepath.Join(home, path[2:]), nil
	default:
		return path, nil
	}
}
