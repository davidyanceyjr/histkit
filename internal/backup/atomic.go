package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RewriteAtomic(targetPath string, contents []byte) error {
	if strings.TrimSpace(targetPath) == "" {
		return fmt.Errorf("rewrite atomic: target path is required")
	}

	mode, err := targetFileMode(targetPath)
	if err != nil {
		return fmt.Errorf("rewrite atomic: %w", err)
	}

	dir := filepath.Dir(targetPath)
	tempFile, err := os.CreateTemp(dir, "."+filepath.Base(targetPath)+".tmp-*")
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

	if err := os.Rename(tempPath, targetPath); err != nil {
		return fmt.Errorf("rewrite atomic: rename temp file %q to %q: %w", tempPath, targetPath, err)
	}
	if err := syncDir(dir); err != nil {
		return fmt.Errorf("rewrite atomic: sync target directory %q: %w", dir, err)
	}

	return nil
}

func targetFileMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err == nil {
		if !info.Mode().IsRegular() {
			return 0, fmt.Errorf("target file %q is not a regular file", path)
		}
		return info.Mode().Perm(), nil
	}
	if os.IsNotExist(err) {
		return 0o600, nil
	}

	return 0, fmt.Errorf("stat target file %q: %w", path, err)
}

func syncDir(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	return dir.Sync()
}
