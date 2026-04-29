package history

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCandidateSources(t *testing.T) {
	sources, err := CandidateSources("/home/tester")
	if err != nil {
		t.Fatalf("CandidateSources returned error: %v", err)
	}

	if len(sources) != 2 {
		t.Fatalf("len(sources) = %d, want 2", len(sources))
	}
	if got, want := sources[0], (Source{Shell: ShellBash, Path: "/home/tester/.bash_history"}); got != want {
		t.Fatalf("sources[0] = %#v, want %#v", got, want)
	}
	if got, want := sources[1], (Source{Shell: ShellZsh, Path: "/home/tester/.zsh_history"}); got != want {
		t.Fatalf("sources[1] = %#v, want %#v", got, want)
	}
}

func TestCandidateSourcesRequiresHome(t *testing.T) {
	if _, err := CandidateSources(""); err == nil {
		t.Fatal("CandidateSources returned nil error for empty home")
	}
}

func TestDetectSourcesFindsExistingFiles(t *testing.T) {
	home := t.TempDir()
	bashPath := filepath.Join(home, ".bash_history")
	zshPath := filepath.Join(home, ".zsh_history")

	if err := os.WriteFile(bashPath, []byte("git status\n"), 0o644); err != nil {
		t.Fatalf("WriteFile bash returned error: %v", err)
	}
	if err := os.WriteFile(zshPath, []byte(": 1712959000:0;git status\n"), 0o644); err != nil {
		t.Fatalf("WriteFile zsh returned error: %v", err)
	}

	sources, err := DetectSources(home, "")
	if err != nil {
		t.Fatalf("DetectSources returned error: %v", err)
	}

	if len(sources) != 2 {
		t.Fatalf("len(sources) = %d, want 2", len(sources))
	}
	if got, want := sources[0], (Source{Shell: ShellBash, Path: bashPath}); got != want {
		t.Fatalf("sources[0] = %#v, want %#v", got, want)
	}
	if got, want := sources[1], (Source{Shell: ShellZsh, Path: zshPath}); got != want {
		t.Fatalf("sources[1] = %#v, want %#v", got, want)
	}
}

func TestDetectSourcesFiltersByShell(t *testing.T) {
	home := t.TempDir()
	zshPath := filepath.Join(home, ".zsh_history")

	if err := os.WriteFile(zshPath, []byte(": 1712959000:0;git status\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	sources, err := DetectSources(home, ShellZsh)
	if err != nil {
		t.Fatalf("DetectSources returned error: %v", err)
	}

	if len(sources) != 1 {
		t.Fatalf("len(sources) = %d, want 1", len(sources))
	}
	if got, want := sources[0], (Source{Shell: ShellZsh, Path: zshPath}); got != want {
		t.Fatalf("sources[0] = %#v, want %#v", got, want)
	}
}

func TestDetectSourcesIgnoresMissingFiles(t *testing.T) {
	home := t.TempDir()

	sources, err := DetectSources(home, "")
	if err != nil {
		t.Fatalf("DetectSources returned error: %v", err)
	}
	if len(sources) != 0 {
		t.Fatalf("len(sources) = %d, want 0", len(sources))
	}
}

func TestDetectSourcesRejectsUnsupportedShell(t *testing.T) {
	if _, err := DetectSources(t.TempDir(), "fish"); err == nil {
		t.Fatal("DetectSources returned nil error for unsupported shell")
	}
}

func TestParserForShell(t *testing.T) {
	tests := []struct {
		shell string
	}{
		{shell: ShellBash},
		{shell: ShellZsh},
	}

	for _, tc := range tests {
		t.Run(tc.shell, func(t *testing.T) {
			parser, err := ParserForShell(tc.shell)
			if err != nil {
				t.Fatalf("ParserForShell returned error: %v", err)
			}
			if parser == nil {
				t.Fatal("ParserForShell returned nil parser")
			}
		})
	}
}

func TestParserForShellRejectsUnsupportedShell(t *testing.T) {
	if _, err := ParserForShell("fish"); err == nil {
		t.Fatal("ParserForShell returned nil error for unsupported shell")
	}
}
