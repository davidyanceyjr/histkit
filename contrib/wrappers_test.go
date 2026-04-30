package contrib_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBashWrapperLoadsSelectionIntoReadlineBuffer(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell wrapper tests require a Unix shell")
	}

	tmpDir := t.TempDir()
	writeFakeHistkit(t, tmpDir, "#!/bin/sh\nprintf 'git status\\n'\n")

	scriptPath := filepath.Join(repoRoot(t), "contrib", "histkit.bash")
	cmd := exec.Command("bash", "-c", `
source "$SCRIPT_PATH"
READLINE_LINE="existing"
READLINE_POINT=0
__histkit_pick_bash
printf '%s\n%s\n' "$READLINE_LINE" "$READLINE_POINT"
`)
	cmd.Env = append(os.Environ(),
		"PATH="+tmpDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"SCRIPT_PATH="+scriptPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bash wrapper command failed: %v\n%s", err, output)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) != 2 {
		t.Fatalf("bash wrapper output lines = %d, want 2: %q", len(lines), output)
	}
	if lines[0] != "git status" {
		t.Fatalf("READLINE_LINE = %q, want git status", lines[0])
	}
	if lines[1] != "10" {
		t.Fatalf("READLINE_POINT = %q, want 10", lines[1])
	}
}

func TestBashWrapperLeavesBufferUntouchedOnEmptySelection(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell wrapper tests require a Unix shell")
	}

	tmpDir := t.TempDir()
	writeFakeHistkit(t, tmpDir, "#!/bin/sh\nexit 0\n")

	scriptPath := filepath.Join(repoRoot(t), "contrib", "histkit.bash")
	cmd := exec.Command("bash", "-c", `
source "$SCRIPT_PATH"
READLINE_LINE="existing"
READLINE_POINT=8
__histkit_pick_bash
printf '%s\n%s\n' "$READLINE_LINE" "$READLINE_POINT"
`)
	cmd.Env = append(os.Environ(),
		"PATH="+tmpDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"SCRIPT_PATH="+scriptPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bash wrapper command failed: %v\n%s", err, output)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) != 2 {
		t.Fatalf("bash wrapper output lines = %d, want 2: %q", len(lines), output)
	}
	if lines[0] != "existing" {
		t.Fatalf("READLINE_LINE = %q, want existing", lines[0])
	}
	if lines[1] != "8" {
		t.Fatalf("READLINE_POINT = %q, want 8", lines[1])
	}
}

func TestZshWrapperScriptContainsBindingHelper(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(repoRoot(t), "contrib", "histkit.zsh"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	text := string(content)
	for _, want := range []string{
		"histkit_pick_zsh()",
		`selected="$(histkit pick "$@")"`,
		"BUFFER=\"$selected\"",
		"CURSOR=${#BUFFER}",
		"zle -N histkit_pick_zsh",
		"bindkey '^R' histkit_pick_zsh",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("zsh wrapper missing %q", want)
		}
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller returned no file")
	}

	return filepath.Dir(filepath.Dir(filename))
}

func writeFakeHistkit(t *testing.T, dir, content string) {
	t.Helper()

	path := filepath.Join(dir, "histkit")
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
