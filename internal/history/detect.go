package history

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Parser func(sourceFile string, r io.Reader) ([]HistoryEntry, []ParseWarning, error)

type Source struct {
	Shell string
	Path  string
}

func CandidateSources(home string) ([]Source, error) {
	if strings.TrimSpace(home) == "" {
		return nil, fmt.Errorf("candidate sources: home directory is required")
	}

	return []Source{
		{Shell: ShellBash, Path: filepath.Join(home, ".bash_history")},
		{Shell: ShellZsh, Path: filepath.Join(home, ".zsh_history")},
	}, nil
}

func DetectSources(home, shell string) ([]Source, error) {
	candidates, err := CandidateSources(home)
	if err != nil {
		return nil, err
	}

	filtered := candidates
	if strings.TrimSpace(shell) != "" {
		if _, err := ParserForShell(shell); err != nil {
			return nil, err
		}

		filtered = nil
		for _, candidate := range candidates {
			if candidate.Shell == shell {
				filtered = append(filtered, candidate)
			}
		}
	}

	var detected []Source
	for _, candidate := range filtered {
		info, err := os.Stat(candidate.Path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("stat history source %q: %w", candidate.Path, err)
		}
		if info.IsDir() {
			continue
		}
		detected = append(detected, candidate)
	}

	return detected, nil
}

func ParserForShell(shell string) (Parser, error) {
	switch shell {
	case ShellBash:
		return ParseBash, nil
	case ShellZsh:
		return ParseZsh, nil
	default:
		return nil, fmt.Errorf("unsupported shell %q", shell)
	}
}
