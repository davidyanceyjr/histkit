package backup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davidyanceyjr/histkit/internal/fsroot"
)

func RewriteAtomic(targetPath string, contents []byte) error {
	if strings.TrimSpace(targetPath) == "" {
		return fmt.Errorf("rewrite atomic: target path is required")
	}

	root, targetRelativePath, absoluteTargetPath, err := newRewriteTargetRoot(targetPath)
	if err != nil {
		return fmt.Errorf("rewrite atomic: %w", err)
	}

	mode, err := targetFileMode(root, targetRelativePath, absoluteTargetPath)
	if err != nil {
		return fmt.Errorf("rewrite atomic: %w", err)
	}

	tempFile, err := root.CreateTemp(".", "."+filepath.Base(absoluteTargetPath)+".tmp-*")
	if err != nil {
		return fmt.Errorf("rewrite atomic: create temp file: %w", err)
	}

	tempPath := tempFile.Name()
	tempClosed := false
	defer func() {
		if !tempClosed {
			_ = tempFile.Close()
		}
		_ = os.Remove(tempPath)
	}()

	if err := tempFile.Chmod(mode); err != nil {
		return fmt.Errorf("rewrite atomic: chmod temp file %q: %w", tempPath, err)
	}
	if _, err := tempFile.Write(contents); err != nil {
		return fmt.Errorf("rewrite atomic: write temp file %q: %w", tempPath, err)
	}
	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("rewrite atomic: sync temp file %q: %w", tempPath, err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("rewrite atomic: close temp file %q: %w", tempPath, err)
	}
	tempClosed = true

	if err := os.Rename(tempPath, absoluteTargetPath); err != nil {
		return fmt.Errorf("rewrite atomic: rename temp file %q to %q: %w", tempPath, absoluteTargetPath, err)
	}
	if err := syncDir(root); err != nil {
		return fmt.Errorf("rewrite atomic: sync target directory %q: %w", root.Path(), err)
	}

	return nil
}

func targetFileMode(root fsroot.Root, path string, displayPath string) (os.FileMode, error) {
	info, err := root.Stat(path)
	if err == nil {
		if !info.Mode().IsRegular() {
			return 0, fmt.Errorf("target file %q is not a regular file", displayPath)
		}
		return info.Mode().Perm(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return 0o600, nil
	}

	return 0, fmt.Errorf("stat target file %q: %w", displayPath, err)
}

func syncDir(root fsroot.Root) error {
	dir, err := root.Open(".")
	if err != nil {
		return err
	}
	defer dir.Close()

	return dir.Sync()
}

func newRewriteTargetRoot(path string) (fsroot.Root, string, string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fsroot.Root{}, "", "", fmt.Errorf("target file %q: resolve absolute path: %w", path, err)
	}

	root, err := fsroot.New(filepath.Dir(absolutePath))
	if err != nil {
		return fsroot.Root{}, "", "", fmt.Errorf("target file %q: %w", path, err)
	}

	return root, filepath.Base(absolutePath), absolutePath, nil
}
