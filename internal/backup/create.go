package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Create(sourceFile, backupDir string, createdAt time.Time, sequence int) (Record, error) {
	if strings.TrimSpace(sourceFile) == "" {
		return Record{}, fmt.Errorf("create backup: source file is required")
	}
	if strings.TrimSpace(backupDir) == "" {
		return Record{}, fmt.Errorf("create backup: backup directory is required")
	}
	if createdAt.IsZero() {
		return Record{}, fmt.Errorf("create backup: created time is required")
	}
	if sequence <= 0 {
		return Record{}, fmt.Errorf("create backup: sequence must be positive")
	}

	checksum, err := ChecksumFile(sourceFile)
	if err != nil {
		return Record{}, fmt.Errorf("create backup: %w", err)
	}

	record, err := BuildRecord(sourceFile, backupDir, checksum, createdAt, sequence)
	if err != nil {
		return Record{}, fmt.Errorf("create backup: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(record.BackupPath), 0o700); err != nil {
		return Record{}, fmt.Errorf("create backup: create backup directory: %w", err)
	}

	if err := copyFile(sourceFile, record.BackupPath); err != nil {
		return Record{}, fmt.Errorf("create backup: %w", err)
	}

	backupChecksum, err := ChecksumFile(record.BackupPath)
	if err != nil {
		_ = os.Remove(record.BackupPath)
		return Record{}, fmt.Errorf("create backup: %w", err)
	}
	if backupChecksum != record.Checksum {
		_ = os.Remove(record.BackupPath)
		return Record{}, fmt.Errorf("create backup: checksum mismatch after copy")
	}

	return record, nil
}

func ChecksumFile(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("checksum file: path is required")
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("checksum file %q: %w", path, err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("checksum file %q: %w", path, err)
	}

	return "sha256:" + hex.EncodeToString(hash.Sum(nil)), nil
}

func copyFile(sourcePath, targetPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source file %q: %w", sourcePath, err)
	}
	defer source.Close()

	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("stat source file %q: %w", sourcePath, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("source file %q is not a regular file", sourcePath)
	}

	target, err := os.OpenFile(targetPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open backup file %q: %w", targetPath, err)
	}

	if _, err := io.Copy(target, source); err != nil {
		_ = target.Close()
		_ = os.Remove(targetPath)
		return fmt.Errorf("copy backup file %q: %w", targetPath, err)
	}

	if err := target.Close(); err != nil {
		_ = os.Remove(targetPath)
		return fmt.Errorf("close backup file %q: %w", targetPath, err)
	}

	return nil
}
