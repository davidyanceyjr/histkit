package sanitize

import (
	"fmt"
	"strings"
	"time"

	"histkit/internal/history"
)

const quarantinedPlaceholder = "[QUARANTINED]"

type QuarantineRecord struct {
	ID         string
	Shell      string
	SourceFile string
	RuleName   string
	Reason     string
	Confidence Confidence
	Action     ActionType
	Original   string
	Preview    string
	CreatedAt  time.Time
}

func (r QuarantineRecord) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return fmt.Errorf("quarantine record id is required")
	}
	if strings.TrimSpace(r.Shell) == "" {
		return fmt.Errorf("quarantine record shell is required")
	}
	if strings.TrimSpace(r.SourceFile) == "" {
		return fmt.Errorf("quarantine record source file is required")
	}
	if strings.TrimSpace(r.RuleName) == "" {
		return fmt.Errorf("quarantine record rule name is required")
	}
	if strings.TrimSpace(r.Reason) == "" {
		return fmt.Errorf("quarantine record reason is required")
	}
	if !isValidConfidence(r.Confidence) {
		return fmt.Errorf("quarantine record confidence %q is invalid", r.Confidence)
	}
	if r.Action != ActionQuarantine {
		return fmt.Errorf("quarantine record action must be %q", ActionQuarantine)
	}
	if strings.TrimSpace(r.Original) == "" {
		return fmt.Errorf("quarantine record original value is required")
	}
	if strings.TrimSpace(r.Preview) == "" {
		return fmt.Errorf("quarantine record preview is required")
	}
	if r.CreatedAt.IsZero() {
		return fmt.Errorf("quarantine record created time is required")
	}

	return nil
}

func BuildQuarantineRecord(entry history.HistoryEntry, match RuleMatch, createdAt time.Time, sequence int) (QuarantineRecord, error) {
	if err := entry.Validate(); err != nil {
		return QuarantineRecord{}, fmt.Errorf("build quarantine record: %w", err)
	}
	if err := match.Validate(); err != nil {
		return QuarantineRecord{}, fmt.Errorf("build quarantine record: %w", err)
	}
	if match.Action != ActionQuarantine {
		return QuarantineRecord{}, fmt.Errorf("build quarantine record: match action must be %q", ActionQuarantine)
	}
	if createdAt.IsZero() {
		return QuarantineRecord{}, fmt.Errorf("build quarantine record: created time is required")
	}
	if sequence <= 0 {
		return QuarantineRecord{}, fmt.Errorf("build quarantine record: sequence must be positive")
	}

	record := QuarantineRecord{
		ID:         quarantineID(createdAt, sequence),
		Shell:      entry.Shell,
		SourceFile: entry.SourceFile,
		RuleName:   match.RuleName,
		Reason:     match.Reason,
		Confidence: match.Confidence,
		Action:     match.Action,
		Original:   match.Before,
		Preview:    quarantinedPlaceholder,
		CreatedAt:  createdAt.UTC(),
	}
	if err := record.Validate(); err != nil {
		return QuarantineRecord{}, err
	}

	return record, nil
}

func BuildQuarantineRecords(entries []history.HistoryEntry, createdAt time.Time) ([]QuarantineRecord, error) {
	if createdAt.IsZero() {
		return nil, fmt.Errorf("build quarantine records: created time is required")
	}

	report, err := PreviewEntries(entries)
	if err != nil {
		return nil, fmt.Errorf("build quarantine records: %w", err)
	}

	records := make([]QuarantineRecord, 0)
	sequence := 1
	for _, item := range report.Items {
		for _, match := range item.Matches {
			if match.Action != ActionQuarantine {
				continue
			}
			record, err := BuildQuarantineRecord(item.Entry, match, createdAt, sequence)
			if err != nil {
				return nil, fmt.Errorf("build quarantine records: %w", err)
			}
			records = append(records, record)
			sequence++
		}
	}

	return records, nil
}

func quarantineID(createdAt time.Time, sequence int) string {
	return fmt.Sprintf("q_%s_%03d", createdAt.UTC().Format("20060102T150405Z"), sequence)
}
