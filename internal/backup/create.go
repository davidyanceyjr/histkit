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

	"github.com/davidyanceyjr/histkit/internal/fsroot"
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

	sourceRoot, sourceRelativePath, sourceAbsolutePath, err := newSourceFileRoot(sourceFile)
	if err != nil {
		return Record{}, fmt.Errorf("create backup: %w", err)
	}
	backupRoot, err := newBackupRoot(backupDir)
	if err != nil {
		return Record{}, fmt.Errorf("create backup: %w", err)
	}

	checksum, err := checksumFileWithRoot(sourceRoot, sourceRelativePath, sourceFile)
	if err != nil {
		return Record{}, fmt.Errorf("create backup: %w", err)
	}

	record, err := BuildRecord(sourceFile, backupRoot.Path(), checksum, createdAt, sequence)
	if err != nil {
		return Record{}, fmt.Errorf("create backup: %w", err)
	}

	backupDirPath, err := backupRoot.Resolve(record.ID)
	if err != nil {
		return Record{}, fmt.Errorf("create backup: resolve backup directory: %w", err)
	}
	if err := os.MkdirAll(backupDirPath, 0o700); err != nil {
		return Record{}, fmt.Errorf("create backup: create backup directory: %w", err)
	}

	targetRelativePath := filepath.Join(record.ID, filepath.Base(sourceAbsolutePath))
	if err := copyFile(sourceRoot, sourceRelativePath, backupRoot, targetRelativePath); err != nil {
		return Record{}, fmt.Errorf("create backup: %w", err)
	}

	backupChecksum, err := checksumFileWithRoot(backupRoot, targetRelativePath, record.BackupPath)
	if err != nil {
		_ = os.Remove(record.BackupPath)
		return Record{}, fmt.Errorf("create backup: %w", err)
	}
	if backupChecksum != record.Checksum {
		_ = os.Remove(record.BackupPath)
		return Record{}, fmt.Errorf("create backup: checksum mismatch after copy")
	}
	if err := WriteRecord(record, backupDir); err != nil {
		_ = os.Remove(record.BackupPath)
		return Record{}, fmt.Errorf("create backup: %w", err)
	}

	return record, nil
}

func ChecksumFile(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("checksum file: path is required")
	}

	root, relativePath, _, err := newSourceFileRoot(path)
	if err != nil {
		return "", err
	}
	return checksumFileWithRoot(root, relativePath, path)
}

func checksumFileWithRoot(root fsroot.Root, path, displayPath string) (string, error) {
	file, err := root.Open(path)
	if err != nil {
		return "", fmt.Errorf("checksum file %q: %w", displayPath, err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("checksum file %q: %w", displayPath, err)
	}

	return "sha256:" + hex.EncodeToString(hash.Sum(nil)), nil
}

func copyFile(sourceRoot fsroot.Root, sourcePath string, targetRoot fsroot.Root, targetPath string) error {
	source, err := sourceRoot.Open(sourcePath)
	if err != nil {
		resolvedSourcePath, resolveErr := sourceRoot.Resolve(sourcePath)
		if resolveErr != nil {
			return fmt.Errorf("open source file: %w", err)
		}
		return fmt.Errorf("open source file %q: %w", resolvedSourcePath, err)
	}
	defer source.Close()

	info, err := source.Stat()
	if err != nil {
		resolvedSourcePath, resolveErr := sourceRoot.Resolve(sourcePath)
		if resolveErr != nil {
			return fmt.Errorf("stat source file: %w", err)
		}
		return fmt.Errorf("stat source file %q: %w", resolvedSourcePath, err)
	}
	if !info.Mode().IsRegular() {
		resolvedSourcePath, resolveErr := sourceRoot.Resolve(sourcePath)
		if resolveErr != nil {
			return fmt.Errorf("source file is not a regular file")
		}
		return fmt.Errorf("source file %q is not a regular file", resolvedSourcePath)
	}

	resolvedTargetPath, err := targetRoot.Resolve(targetPath)
	if err != nil {
		return fmt.Errorf("open backup file: %w", err)
	}

	target, err := targetRoot.OpenFile(targetPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open backup file %q: %w", resolvedTargetPath, err)
	}

	if _, err := io.Copy(target, source); err != nil {
		_ = target.Close()
		_ = os.Remove(resolvedTargetPath)
		return fmt.Errorf("copy backup file %q: %w", resolvedTargetPath, err)
	}

	if err := target.Close(); err != nil {
		_ = os.Remove(resolvedTargetPath)
		return fmt.Errorf("close backup file %q: %w", resolvedTargetPath, err)
	}

	return nil
}

func newSourceFileRoot(path string) (fsroot.Root, string, string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fsroot.Root{}, "", "", fmt.Errorf("source file %q: resolve absolute path: %w", path, err)
	}

	root, err := fsroot.New(filepath.Dir(absolutePath))
	if err != nil {
		return fsroot.Root{}, "", "", fmt.Errorf("source file %q: %w", path, err)
	}

	return root, filepath.Base(absolutePath), absolutePath, nil
}
