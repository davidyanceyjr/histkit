package picker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func Select(ctx context.Context, candidates []Candidate) (Candidate, bool, error) {
	if len(candidates) == 0 {
		return Candidate{}, false, nil
	}

	path, err := exec.LookPath("fzf")
	if err != nil {
		return Candidate{}, false, fmt.Errorf("find fzf: %w", err)
	}

	lines := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		lines = append(lines, candidate.Display())
	}

	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n") + "\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

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
