package backup

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type Record struct {
	ID         string
	SourceFile string
	BackupPath string
	CreatedAt  time.Time
	Checksum   string
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return fmt.Errorf("backup record id is required")
	}
	if strings.TrimSpace(r.SourceFile) == "" {
		return fmt.Errorf("backup record source file is required")
	}
	if strings.TrimSpace(r.BackupPath) == "" {
		return fmt.Errorf("backup record backup path is required")
	}
	if r.CreatedAt.IsZero() {
		return fmt.Errorf("backup record created time is required")
	}
	if strings.TrimSpace(r.Checksum) == "" {
		return fmt.Errorf("backup record checksum is required")
	}
	return nil
}

func BuildRecord(sourceFile, backupDir, checksum string, createdAt time.Time, sequence int) (Record, error) {
	if strings.TrimSpace(sourceFile) == "" {
		return Record{}, fmt.Errorf("build backup record: source file is required")
	}
	if strings.TrimSpace(backupDir) == "" {
		return Record{}, fmt.Errorf("build backup record: backup directory is required")
	}
	if strings.TrimSpace(checksum) == "" {
		return Record{}, fmt.Errorf("build backup record: checksum is required")
	}
	if createdAt.IsZero() {
		return Record{}, fmt.Errorf("build backup record: created time is required")
	}
	if sequence <= 0 {
		return Record{}, fmt.Errorf("build backup record: sequence must be positive")
	}

	id := BackupID(createdAt, sequence)
	record := Record{
		ID:         id,
		SourceFile: sourceFile,
		BackupPath: filepath.Join(backupDir, id, filepath.Base(sourceFile)),
		CreatedAt:  createdAt.UTC(),
		Checksum:   checksum,
	}
	if err := record.Validate(); err != nil {
		return Record{}, err
	}

	return record, nil
}

func BackupID(createdAt time.Time, sequence int) string {
	return fmt.Sprintf("b_%s_%03d", createdAt.UTC().Format("20060102T150405Z"), sequence)
}
