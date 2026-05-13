package fsroot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Root struct {
	path string
}

func New(path string) (Root, error) {
	if strings.TrimSpace(path) == "" {
		return Root{}, fmt.Errorf("new rooted path helper: root path is required")
	}
	if !filepath.IsAbs(path) {
		return Root{}, fmt.Errorf("new rooted path helper: root path %q must be absolute", path)
	}

	return Root{path: filepath.Clean(path)}, nil
}

func (r Root) Path() string {
	return r.path
}

func (r Root) Resolve(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("resolve rooted path: target path is required")
	}
	if strings.TrimSpace(r.path) == "" {
		return "", fmt.Errorf("resolve rooted path: root path is required")
	}

	candidate := path
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(r.path, candidate)
	}
	candidate = filepath.Clean(candidate)

	relative, err := filepath.Rel(r.path, candidate)
	if err != nil {
		return "", fmt.Errorf("resolve rooted path %q against %q: %w", path, r.path, err)
	}
	if relative == "." {
		return candidate, nil
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return "", fmt.Errorf("resolve rooted path %q: %q escapes root %q", path, candidate, r.path)
	}

	return candidate, nil
}

func (r Root) Open(path string) (*os.File, error) {
	resolved, err := r.Resolve(path)
	if err != nil {
		return nil, fmt.Errorf("open rooted path: %w", err)
	}

	// #nosec G304 -- resolved is validated to stay within an absolute trusted root.
	file, err := os.Open(resolved)
	if err != nil {
		return nil, fmt.Errorf("open rooted path %q: %w", resolved, err)
	}
	return file, nil
}

func (r Root) OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	resolved, err := r.Resolve(path)
	if err != nil {
		return nil, fmt.Errorf("open rooted path: %w", err)
	}

	// #nosec G304 -- resolved is validated to stay within an absolute trusted root.
	file, err := os.OpenFile(resolved, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("open rooted path %q: %w", resolved, err)
	}
	return file, nil
}

func (r Root) ReadFile(path string) ([]byte, error) {
	resolved, err := r.Resolve(path)
	if err != nil {
		return nil, fmt.Errorf("read rooted path: %w", err)
	}

	// #nosec G304 -- resolved is validated to stay within an absolute trusted root.
	data, err := os.ReadFile(resolved)
	if err != nil {
		return nil, fmt.Errorf("read rooted path %q: %w", resolved, err)
	}
	return data, nil
}
