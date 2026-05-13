package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/davidyanceyjr/histkit/internal/fsroot"
)

const recordFileName = "record.toml"

type storedRecord struct {
	ID         string `toml:"id"`
	SourceFile string `toml:"source_file"`
	BackupPath string `toml:"backup_path"`
	CreatedAt  string `toml:"created_at"`
	Checksum   string `toml:"checksum"`
}

func RecordPath(backupID, backupDir string) string {
	return filepath.Join(backupDir, backupID, recordFileName)
}

func WriteRecord(record Record, backupDir string) error {
	if err := record.Validate(); err != nil {
		return fmt.Errorf("write backup record: %w", err)
	}
	root, err := newBackupRoot(backupDir)
	if err != nil {
		return fmt.Errorf("write backup record: %w", err)
	}

	path := RecordPath(record.ID, root.Path())
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("write backup record: create metadata directory: %w", err)
	}

	payload := storedRecord{
		ID:         record.ID,
		SourceFile: record.SourceFile,
		BackupPath: record.BackupPath,
		CreatedAt:  record.CreatedAt.UTC().Format(timeLayout),
		Checksum:   record.Checksum,
	}

	file, err := root.OpenFile(filepath.Join(record.ID, recordFileName), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("write backup record: open record file %q: %w", path, err)
	}

	if err := toml.NewEncoder(file).Encode(payload); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return fmt.Errorf("write backup record: encode record file %q: %w", path, err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("write backup record: close record file %q: %w", path, err)
	}

	return nil
}

func LoadRecord(path string) (Record, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return Record{}, fmt.Errorf("load backup record %q: resolve absolute path: %w", path, err)
	}
	root, err := fsroot.New(filepath.Dir(absolutePath))
	if err != nil {
		return Record{}, fmt.Errorf("load backup record %q: %w", path, err)
	}

	var payload storedRecord
	data, err := root.ReadFile(filepath.Base(absolutePath))
	if err != nil {
		return Record{}, fmt.Errorf("load backup record %q: %w", absolutePath, err)
	}
	if _, err := toml.Decode(string(data), &payload); err != nil {
		return Record{}, fmt.Errorf("load backup record %q: %w", absolutePath, err)
	}

	createdAt, err := parseStoredTime(payload.CreatedAt)
	if err != nil {
		return Record{}, fmt.Errorf("load backup record %q: %w", absolutePath, err)
	}

	record := Record{
		ID:         payload.ID,
		SourceFile: payload.SourceFile,
		BackupPath: payload.BackupPath,
		CreatedAt:  createdAt,
		Checksum:   payload.Checksum,
	}
	if err := record.Validate(); err != nil {
		return Record{}, fmt.Errorf("load backup record %q: %w", absolutePath, err)
	}

	return record, nil
}

func ListRecords(backupDir string) ([]Record, error) {
	root, err := newBackupRoot(backupDir)
	if err != nil {
		return nil, fmt.Errorf("list backup records: %w", err)
	}

	entries, err := os.ReadDir(root.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list backup records: read backup directory %q: %w", root.Path(), err)
	}

	records := make([]Record, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		record, err := LoadRecord(RecordPath(entry.Name(), root.Path()))
		if err != nil {
			return nil, fmt.Errorf("list backup records: %w", err)
		}
		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.After(records[j].CreatedAt)
	})
	return records, nil
}

func FindRecord(backupDir, backupID string) (Record, error) {
	if strings.TrimSpace(backupID) == "" {
		return Record{}, fmt.Errorf("find backup record: backup id is required")
	}
	root, err := newBackupRoot(backupDir)
	if err != nil {
		return Record{}, fmt.Errorf("find backup record: %w", err)
	}

	record, err := LoadRecord(RecordPath(backupID, root.Path()))
	if err != nil {
		return Record{}, fmt.Errorf("find backup record: %w", err)
	}
	return record, nil
}

func Restore(record Record) error {
	if err := record.Validate(); err != nil {
		return fmt.Errorf("restore backup: %w", err)
	}

	checksum, err := ChecksumFile(record.BackupPath)
	if err != nil {
		return fmt.Errorf("restore backup: %w", err)
	}
	if checksum != record.Checksum {
		return fmt.Errorf("restore backup: checksum mismatch for %q", record.BackupPath)
	}

	contents, err := os.ReadFile(record.BackupPath)
	if err != nil {
		return fmt.Errorf("restore backup: read backup file %q: %w", record.BackupPath, err)
	}
	if err := RewriteAtomic(record.SourceFile, contents); err != nil {
		return fmt.Errorf("restore backup: %w", err)
	}

	return nil
}

func parseStoredTime(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, fmt.Errorf("stored created time is required")
	}
	createdAt, err := time.Parse(timeLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse stored created time %q: %w", value, err)
	}
	return createdAt.UTC(), nil
}

const timeLayout = "2006-01-02T15:04:05Z07:00"

func newBackupRoot(backupDir string) (fsroot.Root, error) {
	if strings.TrimSpace(backupDir) == "" {
		return fsroot.Root{}, fmt.Errorf("backup directory is required")
	}

	absolutePath, err := filepath.Abs(backupDir)
	if err != nil {
		return fsroot.Root{}, fmt.Errorf("backup directory %q: resolve absolute path: %w", backupDir, err)
	}
	root, err := fsroot.New(absolutePath)
	if err != nil {
		return fsroot.Root{}, fmt.Errorf("backup directory %q: %w", backupDir, err)
	}
	return root, nil
}
