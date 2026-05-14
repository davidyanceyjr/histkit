package picker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var openTTY = func() (*os.File, error) {
	return os.OpenFile("/dev/tty", os.O_RDWR, 0)
}

var findFZF = exec.LookPath

func Select(ctx context.Context, candidates []Candidate) (Candidate, bool, error) {
	if len(candidates) == 0 {
		return Candidate{}, false, nil
	}

	path, err := findFZF("fzf")
	if err != nil {
		return Candidate{}, false, fmt.Errorf("find fzf: %w", err)
	}
	if err := validateFZFPath(path); err != nil {
		return Candidate{}, false, fmt.Errorf("find fzf: %w", err)
	}

	lines := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		lines = append(lines, candidate.Display())
	}

	// #nosec G204 -- path comes from exec.LookPath("fzf") and must pass local absolute-path,
	// basename, regular-file, and executable-bit validation before launch.
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n") + "\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	tty, err := openTTY()
	if err == nil {
		defer tty.Close()
		cmd.Stderr = io.MultiWriter(&stderr, tty)
	}

	if err := cmd.Run(); err != nil {
		if isNoSelection(err) {
			return Candidate{}, false, nil
		}
		if stderr.Len() > 0 {
			return Candidate{}, false, fmt.Errorf("run fzf: %s", strings.TrimSpace(stderr.String()))
		}
		return Candidate{}, false, fmt.Errorf("run fzf: %w", err)
	}

	line := strings.TrimRight(stdout.String(), "\r\n")
	if line == "" {
		return Candidate{}, false, nil
	}

	selected, err := ParseSelectedLine(line)
	if err != nil {
		return Candidate{}, false, err
	}

	return selected, true, nil
}

func isNoSelection(err error) bool {
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}

	code := exitErr.ExitCode()
	return code == 1 || code == 130
}

func validateFZFPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("fzf path is required")
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("fzf path must be absolute: %s", path)
	}
	if filepath.Base(path) != "fzf" {
		return fmt.Errorf("fzf path must end with fzf: %s", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat fzf path %s: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("fzf path is a directory: %s", path)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("fzf path is not a regular file: %s", path)
	}
	if info.Mode().Perm()&0o111 == 0 {
		return fmt.Errorf("fzf path is not executable: %s", path)
	}

	return nil
}
