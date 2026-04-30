package sanitize

import (
	"strings"
	"testing"
)

func TestRedactCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		rule    Rule
		want    string
	}{
		{
			name:    "exact redaction",
			command: "clear",
			rule: Rule{
				Name:       "drop-clear",
				Type:       RuleExact,
				Pattern:    "clear",
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
				Reason:     "Mask exact command",
			},
			want: redactedPlaceholder,
		},
		{
			name:    "contains redaction",
			command: "cat <<'EOF'\nBEGIN OPENSSH PRIVATE KEY\nEOF",
			rule: Rule{
				Name:       "private-key-marker",
				Type:       RuleContains,
				Pattern:    "BEGIN OPENSSH PRIVATE KEY",
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
				Reason:     "Mask private key marker",
			},
			want: "cat <<'EOF'\n[REDACTED]\nEOF",
		},
		{
			name:    "regex redaction",
			command: "kubectl config set-credentials demo --token abc123",
			rule: Rule{
				Name:       "kubectl-token",
				Type:       RuleRegex,
				Pattern:    `kubectl.*--token[ =][^ ]+`,
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
				Reason:     "Redact inline kubectl tokens",
			},
			want: redactedPlaceholder,
		},
		{
			name:    "keyword-group redaction",
			command: "curl https://user:secret@example.com/path",
			rule: Rule{
				Name:       "url-creds",
				Type:       RuleKeywordGroup,
				Keywords:   []string{"https://", "@", ":"},
				Action:     ActionRedact,
				Confidence: ConfidenceMedium,
				Reason:     "Mask likely credential-bearing URL separators",
			},
			want: "curl [REDACTED]user[REDACTED]secret[REDACTED]example.com/path",
		},
		{
			name:    "high entropy token redaction",
			command: "export TOKEN=AbcdefGhijklmnopQRST12uvwxYZ34",
			rule: Rule{
				Name:       "high-entropy-token",
				Type:       RuleHeuristic,
				Detector:   "high_entropy_token",
				Action:     ActionRedact,
				Confidence: ConfidenceMedium,
				Reason:     "Mask likely secret-like tokens",
			},
			want: "export TOKEN=[REDACTED]",
		},
		{
			name:    "large paste blob redaction",
			command: "cat <<'EOF'\n" + strings.Repeat("x", 520) + "\nEOF",
			rule: Rule{
				Name:       "large-paste-blob",
				Type:       RuleHeuristic,
				Detector:   "large_paste_blob",
				Action:     ActionRedact,
				Confidence: ConfidenceMedium,
				Reason:     "Mask large accidental paste blobs",
			},
			want: redactedPlaceholder,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RedactCommand(tc.command, tc.rule)
			if err != nil {
				t.Fatalf("RedactCommand returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("RedactCommand = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRedactCommandRejectsNonMatchingInput(t *testing.T) {
	_, err := RedactCommand("ssh prod-box", Rule{
		Name:       "private-key-marker",
		Type:       RuleContains,
		Pattern:    "BEGIN OPENSSH PRIVATE KEY",
		Action:     ActionRedact,
		Confidence: ConfidenceHigh,
		Reason:     "Mask private key marker",
	})
	if err == nil {
		t.Fatal("RedactCommand returned nil error for non-matching input")
	}
}

func TestRedactCommandRejectsUnknownHeuristicDetector(t *testing.T) {
	_, err := RedactCommand("echo hi", Rule{
		Name:       "unknown",
		Type:       RuleHeuristic,
		Detector:   "does_not_exist",
		Action:     ActionRedact,
		Confidence: ConfidenceMedium,
		Reason:     "desc",
	})
	if err == nil {
		t.Fatal("RedactCommand returned nil error for unknown detector")
	}
}
