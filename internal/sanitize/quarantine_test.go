package sanitize

import (
	"testing"
	"time"

	"histkit/internal/history"
)

func TestQuarantineRecordValidate(t *testing.T) {
	record := QuarantineRecord{
		ID:         "q_20260430T120000Z_001",
		Shell:      history.ShellBash,
		SourceFile: "/home/tester/.bash_history",
		RuleName:   "private-key-block",
		Reason:     "Quarantine pasted private key material",
		Confidence: ConfidenceHigh,
		Action:     ActionQuarantine,
		Original:   "BEGIN OPENSSH PRIVATE KEY",
		Preview:    quarantinedPlaceholder,
		CreatedAt:  time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC),
	}

	if err := record.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestBuildQuarantineRecord(t *testing.T) {
	entry := history.HistoryEntry{
		Shell:      history.ShellBash,
		SourceFile: "/home/tester/.bash_history",
		RawLine:    "cat <<'EOF'\nBEGIN OPENSSH PRIVATE KEY\nabc\nEOF",
		Command:    "cat <<'EOF'\nBEGIN OPENSSH PRIVATE KEY\nabc\nEOF",
	}
	match := RuleMatch{
		RuleName:   "private-key-block",
		Reason:     "Quarantine pasted private key material",
		Confidence: ConfidenceHigh,
		Action:     ActionQuarantine,
		Before:     entry.Command,
	}

	record, err := BuildQuarantineRecord(entry, match, time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC), 1)
	if err != nil {
		t.Fatalf("BuildQuarantineRecord returned error: %v", err)
	}
	if record.ID != "q_20260430T120000Z_001" {
		t.Fatalf("record.ID = %q, want q_20260430T120000Z_001", record.ID)
	}
	if record.Preview != quarantinedPlaceholder {
		t.Fatalf("record.Preview = %q, want %q", record.Preview, quarantinedPlaceholder)
	}
	if record.Original != entry.Command {
		t.Fatalf("record.Original = %q, want original command", record.Original)
	}
}

func TestBuildQuarantineRecordsFromEntries(t *testing.T) {
	entries := []history.HistoryEntry{
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "cat <<'EOF'\nBEGIN OPENSSH PRIVATE KEY\nabc\nEOF",
			Command:    "cat <<'EOF'\nBEGIN OPENSSH PRIVATE KEY\nabc\nEOF",
		},
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "export TOKEN=AbcdefGhijklmnopQRST12uvwxYZ34",
			Command:    "export TOKEN=AbcdefGhijklmnopQRST12uvwxYZ34",
		},
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "clear",
			Command:    "clear",
		},
	}

	records, err := BuildQuarantineRecords(entries, time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("BuildQuarantineRecords returned error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("len(records) = %d, want 2", len(records))
	}
	if records[0].RuleName != "private-key-block" {
		t.Fatalf("records[0].RuleName = %q, want private-key-block", records[0].RuleName)
	}
	if records[1].RuleName != "high-entropy-token" {
		t.Fatalf("records[1].RuleName = %q, want high-entropy-token", records[1].RuleName)
	}
}

func TestBuildQuarantineRecordRejectsNonQuarantineMatch(t *testing.T) {
	entry := history.HistoryEntry{
		Shell:      history.ShellBash,
		SourceFile: "/home/tester/.bash_history",
		RawLine:    "mysql --password hunter2",
		Command:    "mysql --password hunter2",
	}
	match := RuleMatch{
		RuleName:   "inline-password-flag",
		Reason:     "Redact inline password values",
		Confidence: ConfidenceHigh,
		Action:     ActionRedact,
		Before:     entry.Command,
		After:      "mysql --password [REDACTED]",
	}

	if _, err := BuildQuarantineRecord(entry, match, time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC), 1); err == nil {
		t.Fatal("BuildQuarantineRecord returned nil error for non-quarantine match")
	}
}

func TestBuildQuarantineRecordsReturnsNoneWhenNoQuarantineMatches(t *testing.T) {
	records, err := BuildQuarantineRecords([]history.HistoryEntry{
		{
			Shell:      history.ShellBash,
			SourceFile: "/home/tester/.bash_history",
			RawLine:    "echo hi",
			Command:    "echo hi",
		},
	}, time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("BuildQuarantineRecords returned error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("len(records) = %d, want 0", len(records))
	}
}
