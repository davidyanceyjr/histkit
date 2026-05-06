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
	output := stdout.String()
	assertHelpContains(t, output,
		"histkit is a conservative CLI for shell history hygiene, reusable snippets, and fuzzy command recall.",
		"It keeps raw shell history, the local history index, and snippets separate by design.",
		"The normal workflow is: doctor -> scan -> stats or pick -> clean --dry-run -> clean --apply -> restore.",
		"doctor  Check config, local paths, detected history sources, fzf, and related environment state",
		"scan    Parse supported shell history sources and refresh the local SQLite history index",
		"stats   Show indexed history counts by shell and source",
		"pick    Select a command from indexed history and snippets through fzf",
		"clean   Preview or apply cleanup rules to shell history with backups and audit logging",
		"restore List recorded backups or restore a specific backup",
		"Use \"histkit help <command>\" or \"histkit <command> --help\" for command-specific help.",
	)
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

func assertHelpContains(t *testing.T, output string, wants ...string) {
	t.Helper()

	for _, want := range wants {
		if !strings.Contains(output, want) {
			t.Fatalf("expected help output to contain %q, got %q", want, output)
		}
	}
}
