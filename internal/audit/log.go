package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Append(path string, record Record) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("append audit log: path is required")
	}
	if err := record.Validate(); err != nil {
		return fmt.Errorf("append audit log: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("append audit log: create parent directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("append audit log: open file %q: %w", path, err)
	}
	defer file.Close()

	if _, err := file.WriteString(RenderLine(record) + "\n"); err != nil {
		return fmt.Errorf("append audit log: write file %q: %w", path, err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("append audit log: sync file %q: %w", path, err)
	}

	return nil
}
