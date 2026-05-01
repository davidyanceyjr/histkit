package sanitize

import (
	"bytes"
	"strings"
	"testing"

	"histkit/internal/history"
)

func TestApplyToSourceBashRewritesDeleteRedactAndQuarantine(t *testing.T) {
	source := history.Source{
		Shell: history.ShellBash,
		Path:  "/home/tester/.bash_history",
	}
	input := []byte("echo hi\npwd\nmysql --password hunter2\nexport TOKEN=AbcdefGhijklmnopQRST12uvwxYZ34\n   \n")

	report, err := ApplyToSource(source, input)
	if err != nil {
		t.Fatalf("ApplyToSource returned error: %v", err)
	}

	want := "echo hi\nmysql [REDACTED]\n[QUARANTINED]\n   \n"
	if got := string(report.RewrittenContent); got != want {
		t.Fatalf("RewrittenContent = %q, want %q", got, want)
	}
	if report.TotalLines != 5 {
		t.Fatalf("TotalLines = %d, want 5", report.TotalLines)
	}
	if report.ParsedEntries != 4 {
		t.Fatalf("ParsedEntries = %d, want 4", report.ParsedEntries)
	}
	if report.MatchedEntries != 3 {
		t.Fatalf("MatchedEntries = %d, want 3", report.MatchedEntries)
	}
	if report.RewrittenEntries != 2 {
		t.Fatalf("RewrittenEntries = %d, want 2", report.RewrittenEntries)
	}
	if report.DeletedEntries != 1 {
		t.Fatalf("DeletedEntries = %d, want 1", report.DeletedEntries)
	}
	if report.CountsByAction[ActionDelete] != 1 {
		t.Fatalf("delete count = %d, want 1", report.CountsByAction[ActionDelete])
	}
	if report.CountsByAction[ActionRedact] != 1 {
		t.Fatalf("redact count = %d, want 1", report.CountsByAction[ActionRedact])
	}
	if report.CountsByAction[ActionQuarantine] != 1 {
		t.Fatalf("quarantine count = %d, want 1", report.CountsByAction[ActionQuarantine])
	}
	if report.CountsByConfidence[ConfidenceHigh] != 2 {
		t.Fatalf("high confidence count = %d, want 2", report.CountsByConfidence[ConfidenceHigh])
	}
	if report.CountsByConfidence[ConfidenceMedium] != 1 {
		t.Fatalf("medium confidence count = %d, want 1", report.CountsByConfidence[ConfidenceMedium])
	}
	if got := strings.Join(report.RuleNames, ","); got != "high-entropy-token,inline-password-flag,pwd-command" {
		t.Fatalf("RuleNames = %q", got)
	}
}

func TestApplyToSourceZshExtendedHistoryPreservesMetadata(t *testing.T) {
	source := history.Source{
		Shell: history.ShellZsh,
		Path:  "/home/tester/.zsh_history",
	}
	input := []byte(": 1712959000:0;pwd\n: 1712959015:2;curl https://user:secret@example.com/path\n")

	report, err := ApplyToSource(source, input)
	if err != nil {
		t.Fatalf("ApplyToSource returned error: %v", err)
	}

	want := ": 1712959015:2;curl [REDACTED]/path\n"
	if got := string(report.RewrittenContent); got != want {
		t.Fatalf("RewrittenContent = %q, want %q", got, want)
	}
}

func TestApplyToSourcePreservesUnparsedZshLines(t *testing.T) {
	source := history.Source{
		Shell: history.ShellZsh,
		Path:  "/home/tester/.zsh_history",
	}
	input := []byte(": malformed\n: 1712959000:0;pwd\n")

	report, err := ApplyToSource(source, input)
	if err != nil {
		t.Fatalf("ApplyToSource returned error: %v", err)
	}

	if !bytes.Equal(report.RewrittenContent, []byte(": malformed\n")) {
		t.Fatalf("RewrittenContent = %q, want preserved malformed line only", string(report.RewrittenContent))
	}
}

func TestApplyToSourceRejectsUnsupportedShell(t *testing.T) {
	_, err := ApplyToSource(history.Source{Shell: "fish", Path: "/tmp/fish_history"}, []byte("pwd\n"))
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
}
