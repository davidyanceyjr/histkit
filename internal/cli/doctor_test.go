package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteDoctorReportsWarningsForFreshHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", "")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"doctor"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	output := stdout.String()
	for _, want := range []string{
		"doctor overall status: warn",
		"config: OK - default config not present; using built-in defaults",
		"state_dir: WARN - state directory does not exist yet; writable parent detected:",
		"history_sources: WARN - no supported history files detected",
		"history_db: WARN - history database does not exist yet; writable parent detected:",
		"fzf: WARN - fzf not found in PATH",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected doctor output to contain %q, got %q", want, output)
		}
	}
}

func TestExecuteDoctorReportsHealthyChecks(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	stateDir := filepath.Join(home, ".local", "share", "histkit")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatalf("MkdirAll stateDir returned error: %v", err)
	}
	configDir := filepath.Join(home, ".config", "histkit")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll configDir returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.toml"), []byte("[general]\ndefault_shell = \"bash\"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile config returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(home, ".bash_history"), []byte("pwd\n"), 0o600); err != nil {
		t.Fatalf("WriteFile history returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, "history.db"), []byte{}, 0o600); err != nil {
		t.Fatalf("WriteFile history.db returned error: %v", err)
	}

	fzfDir := filepath.Join(home, "bin")
	if err := os.MkdirAll(fzfDir, 0o755); err != nil {
		t.Fatalf("MkdirAll fzfDir returned error: %v", err)
	}
	fzfPath := filepath.Join(fzfDir, "fzf")
	if err := os.WriteFile(fzfPath, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("WriteFile fzf returned error: %v", err)
	}
	t.Setenv("PATH", fzfDir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := Execute([]string{"doctor"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}

	output := stdout.String()
	for _, want := range []string{
		"doctor overall status: ok",
		"config: OK - default config loaded:",
		"state_dir: OK - state directory is writable:",
		"history_sources: OK - readable history sources: bash (",
		"history_db: OK - history database is accessible:",
		"fzf: OK - fzf available at",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected doctor output to contain %q, got %q", want, output)
		}
	}
}

func TestExecuteDoctorRejectsMissingExplicitConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Execute([]string{"doctor", "--config", "~/missing.toml"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute returned unexpected error: %v", err)
	}
	if !strings.Contains(stdout.String(), "config: FAIL - requested config file not found:") {
		t.Fatalf("expected explicit config failure in doctor output, got %q", stdout.String())
	}
}
