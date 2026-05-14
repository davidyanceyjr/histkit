package picker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectReturnsErrorWhenFZFNotFound(t *testing.T) {
	restore := stubFZFFinder(t, func(string) (string, error) {
		return "", execErr("not found")
	})
	defer restore()

	_, ok, err := Select(context.Background(), []Candidate{{Label: LabelHistory, Command: "git status"}})
	if err == nil {
		t.Fatal("Select returned nil error without fzf in PATH")
	}
	if ok {
		t.Fatal("Select returned ok=true without fzf in PATH")
	}
	if !strings.Contains(err.Error(), "find fzf") {
		t.Fatalf("Select error = %q, want mention of fzf lookup", err)
	}
}

func TestSelectReturnsErrorWhenFZFPathIsRelative(t *testing.T) {
	restore := stubFZFFinder(t, func(string) (string, error) {
		return "fzf", nil
	})
	defer restore()

	_, ok, err := Select(context.Background(), []Candidate{{Label: LabelHistory, Command: "git status"}})
	if err == nil {
		t.Fatal("Select returned nil error for relative fzf path")
	}
	if ok {
		t.Fatal("Select returned ok=true for relative fzf path")
	}
	if !strings.Contains(err.Error(), "fzf path must be absolute") {
		t.Fatalf("Select error = %q, want absolute-path failure", err)
	}
}

func TestSelectReturnsErrorWhenFZFPathHasWrongBasename(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fzf-safe")
	writeFileMode(t, path, "#!/bin/sh\nexit 0\n", 0o755)

	restore := stubFZFFinder(t, func(string) (string, error) {
		return path, nil
	})
	defer restore()

	_, ok, err := Select(context.Background(), []Candidate{{Label: LabelHistory, Command: "git status"}})
	if err == nil {
		t.Fatal("Select returned nil error for wrong fzf basename")
	}
	if ok {
		t.Fatal("Select returned ok=true for wrong fzf basename")
	}
	if !strings.Contains(err.Error(), "fzf path must end with fzf") {
		t.Fatalf("Select error = %q, want basename failure", err)
	}
}

func TestSelectReturnsChosenCandidate(t *testing.T) {
	fzfDir := t.TempDir()
	writeFakeFZF(t, fzfDir, "#!/bin/sh\nIFS= read -r _\nIFS= read -r line\nprintf '%s\\n' \"$line\"\n")
	t.Setenv("PATH", fzfDir)
	ttyPath := useFakeTTY(t)

	selected, ok, err := Select(context.Background(), []Candidate{
		{Label: LabelHistory, Command: "git status"},
		{Label: LabelSnippet, Command: "find {{path}}"},
	})
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if !ok {
		t.Fatal("Select returned ok=false, want true")
	}
	if selected.Label != LabelSnippet || selected.Command != "find {{path}}" {
		t.Fatalf("selected = %#v, want snippet candidate", selected)
	}
	if content, err := os.ReadFile(ttyPath); err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	} else if len(content) != 0 {
		t.Fatalf("fake tty content = %q, want empty output for quiet fzf", string(content))
	}
}

func TestSelectReturnsNoSelectionForAbort(t *testing.T) {
	fzfDir := t.TempDir()
	writeFakeFZF(t, fzfDir, "#!/bin/sh\nexit 130\n")
	t.Setenv("PATH", fzfDir)
	useFakeTTY(t)

	_, ok, err := Select(context.Background(), []Candidate{{Label: LabelHistory, Command: "git status"}})
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if ok {
		t.Fatal("Select returned ok=true, want false for abort")
	}
}

func TestSelectMirrorsFZFStderrToTTYWhenAvailable(t *testing.T) {
	fzfDir := t.TempDir()
	writeFakeFZF(t, fzfDir, "#!/bin/sh\nprintf 'picker-ui\\n' >&2\nIFS= read -r _\nIFS= read -r line\nprintf '%s\\n' \"$line\"\n")
	t.Setenv("PATH", fzfDir)

	ttyPath := useFakeTTY(t)

	selected, ok, err := Select(context.Background(), []Candidate{
		{Label: LabelHistory, Command: "git status"},
		{Label: LabelSnippet, Command: "find {{path}}"},
	})
	if err != nil {
		t.Fatalf("Select returned error: %v", err)
	}
	if !ok {
		t.Fatal("Select returned ok=false, want true")
	}
	if selected.Label != LabelSnippet || selected.Command != "find {{path}}" {
		t.Fatalf("selected = %#v, want snippet candidate", selected)
	}

	content, err := os.ReadFile(ttyPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if got := string(content); got != "picker-ui\n" {
		t.Fatalf("fake tty content = %q, want %q", got, "picker-ui\n")
	}
}

func TestSelectReturnsCapturedErrorWhenTTYUnavailable(t *testing.T) {
	fzfDir := t.TempDir()
	writeFakeFZF(t, fzfDir, "#!/bin/sh\nprintf 'boom\\n' >&2\nexit 2\n")
	t.Setenv("PATH", fzfDir)

	restore := stubTTYOpener(t, nil, fmt.Errorf("no controlling tty"))
	defer restore()

	_, ok, err := Select(context.Background(), []Candidate{{Label: LabelHistory, Command: "git status"}})
	if err == nil {
		t.Fatal("Select returned nil error, want fzf failure")
	}
	if ok {
		t.Fatal("Select returned ok=true, want false on failure")
	}
	if !strings.Contains(err.Error(), "run fzf: boom") {
		t.Fatalf("Select error = %q, want captured stderr message", err)
	}
}

func writeFakeFZF(t *testing.T, dir, content string) {
	t.Helper()

	path := filepath.Join(dir, "fzf")
	writeFileMode(t, path, content, 0o755)
}

func useFakeTTY(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "tty")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o600)
	if err != nil {
		t.Fatalf("OpenFile returned error: %v", err)
	}

	restore := stubTTYOpener(t, file, nil)
	t.Cleanup(func() {
		restore()
	})

	return path
}

func stubTTYOpener(t *testing.T, file *os.File, err error) func() {
	t.Helper()

	original := openTTY
	openTTY = func() (*os.File, error) {
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	return func() {
		openTTY = original
		if file != nil {
			_ = file.Close()
		}
	}
}

func stubFZFFinder(t *testing.T, finder func(string) (string, error)) func() {
	t.Helper()

	original := findFZF
	findFZF = finder
	return func() {
		findFZF = original
	}
}

func writeFileMode(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}

type execError string

func (e execError) Error() string {
	return string(e)
}

func execErr(message string) error {
	return errors.New(message)
}
