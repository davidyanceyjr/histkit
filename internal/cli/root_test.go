package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if err := Execute([]string{"--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("expected help output, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "scan") {
		t.Fatalf("expected command list in help output, got %q", stdout.String())
	}
}

func TestExecuteScanPlaceholder(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	assertCommandOutput(t, []string{"scan"}, "scan complete:")
}

func TestExecuteStatsPlaceholder(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	assertCommandOutput(t, []string{"stats"}, "Indexed history entries: 0")
}

func TestExecuteDoctorPlaceholder(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	assertCommandOutput(t, []string{"doctor"}, "doctor overall status:")
}

func TestExecuteUnknownCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Execute([]string{"unknown"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(err.Error(), `unknown command "unknown"`) {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Fatalf("expected root usage on stderr, got %q", stderr.String())
	}
}

func assertCommandOutput(t *testing.T, args []string, want string) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if err := Execute(args, &stdout, &stderr); err != nil {
		t.Fatalf("Execute(%v) returned error: %v", args, err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output for %v, got %q", args, stderr.String())
	}
	if !strings.Contains(stdout.String(), want) {
		t.Fatalf("expected stdout for %v to contain %q, got %q", args, want, stdout.String())
	}
}
