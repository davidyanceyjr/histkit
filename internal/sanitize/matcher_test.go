package sanitize

import (
	"strings"
	"testing"

	"histkit/internal/history"
)

func TestMatchRule(t *testing.T) {
	tests := []struct {
		name    string
		command string
		rule    Rule
		want    bool
	}{
		{
			name:    "exact match",
			command: "clear",
			rule: Rule{
				Name:       "drop-clear",
				Type:       RuleExact,
				Pattern:    "clear",
				Action:     ActionDelete,
				Confidence: ConfidenceHigh,
				Reason:     "Drop trivial clear commands",
			},
			want: true,
		},
		{
			name:    "contains match",
			command: "cat <<'EOF'\nBEGIN OPENSSH PRIVATE KEY\nEOF",
			rule: Rule{
				Name:       "private-key-marker",
				Type:       RuleContains,
				Pattern:    "BEGIN OPENSSH PRIVATE KEY",
				Action:     ActionQuarantine,
				Confidence: ConfidenceHigh,
				Reason:     "Quarantine private key material",
			},
			want: true,
		},
		{
			name:    "regex match",
			command: "kubectl config set-credentials demo --token abc123",
			rule: Rule{
				Name:       "kubectl-token",
				Type:       RuleRegex,
				Pattern:    `kubectl.*--token[ =][^ ]+`,
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
				Reason:     "Redact inline kubectl tokens",
			},
			want: true,
		},
		{
			name:    "keyword group match",
			command: "curl https://user:secret@example.com/path",
			rule: Rule{
				Name:       "url-creds",
				Type:       RuleKeywordGroup,
				Keywords:   []string{"https://", "@", ":"},
				Action:     ActionQuarantine,
				Confidence: ConfidenceMedium,
				Reason:     "Quarantine likely credential-bearing URLs",
			},
			want: true,
		},
		{
			name:    "heuristic match",
			command: "export TOKEN=AbcdefGhijklmnopQRST12uvwxYZ34",
			rule: Rule{
				Name:       "high-entropy-token",
				Type:       RuleHeuristic,
				Detector:   "high_entropy_token",
				Action:     ActionQuarantine,
				Confidence: ConfidenceMedium,
				Reason:     "Quarantine likely secret-like tokens",
			},
			want: true,
		},
		{
			name:    "no broad ssh false positive",
			command: "ssh prod-box",
			rule: Rule{
				Name:       "private-key-marker",
				Type:       RuleContains,
				Pattern:    "BEGIN OPENSSH PRIVATE KEY",
				Action:     ActionQuarantine,
				Confidence: ConfidenceHigh,
				Reason:     "Quarantine private key material",
			},
			want: false,
		},
		{
			name:    "keyword group requires all keywords",
			command: "curl https://example.com/path",
			rule: Rule{
				Name:       "url-creds",
				Type:       RuleKeywordGroup,
				Keywords:   []string{"https://", "@", ":"},
				Action:     ActionQuarantine,
				Confidence: ConfidenceMedium,
				Reason:     "Quarantine likely credential-bearing URLs",
			},
			want: false,
		},
		{
			name:    "exact rule does not trim command",
			command: "clear ",
			rule: Rule{
				Name:       "drop-clear",
				Type:       RuleExact,
				Pattern:    "clear",
				Action:     ActionDelete,
				Confidence: ConfidenceHigh,
				Reason:     "Drop trivial clear commands",
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			match, got, err := MatchRule(tc.command, tc.rule)
			if err != nil {
				t.Fatalf("MatchRule returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("MatchRule matched = %v, want %v", got, tc.want)
			}
			if !got {
				return
			}
			if match.RuleName != tc.rule.Name || match.Action != tc.rule.Action || match.Before != tc.command {
				t.Fatalf("match = %#v, want rule/action/before copied from rule and command", match)
			}
			if tc.rule.Action == ActionRedact && match.After == tc.command {
				t.Fatalf("match.After = %q, want transformed redacted output", match.After)
			}
		})
	}
}

func TestMatchRulePopulatesRedactedAfterValue(t *testing.T) {
	match, ok, err := MatchRule("kubectl config set-credentials demo --token abc123", Rule{
		Name:       "kubectl-token",
		Type:       RuleRegex,
		Pattern:    `kubectl.*--token[ =][^ ]+`,
		Action:     ActionRedact,
		Confidence: ConfidenceHigh,
		Reason:     "Redact inline kubectl tokens",
	})
	if err != nil {
		t.Fatalf("MatchRule returned error: %v", err)
	}
	if !ok {
		t.Fatal("MatchRule returned ok=false, want true")
	}
	if match.After != redactedPlaceholder {
		t.Fatalf("match.After = %q, want %q", match.After, redactedPlaceholder)
	}
}

func TestMatchRuleRejectsUnknownHeuristicDetector(t *testing.T) {
	_, _, err := MatchRule("echo hi", Rule{
		Name:       "unknown",
		Type:       RuleHeuristic,
		Detector:   "does_not_exist",
		Action:     ActionQuarantine,
		Confidence: ConfidenceMedium,
		Reason:     "desc",
	})
	if err == nil {
		t.Fatal("MatchRule returned nil error for unknown detector")
	}
}

func TestMatchEntryReturnsAllMatchingRules(t *testing.T) {
	entry := history.HistoryEntry{
		Shell:      history.ShellBash,
		SourceFile: "/home/tester/.bash_history",
		RawLine:    "kubectl config set-credentials demo --token abc123",
		Command:    "kubectl config set-credentials demo --token abc123",
	}

	matches, err := MatchEntry(entry, []Rule{
		{
			Name:       "contains-kubectl-token-flag",
			Type:       RuleContains,
			Pattern:    "--token",
			Action:     ActionQuarantine,
			Confidence: ConfidenceMedium,
			Reason:     "Flag likely contains a token",
		},
		{
			Name:       "regex-kubectl-token",
			Type:       RuleRegex,
			Pattern:    `kubectl.*--token[ =][^ ]+`,
			Action:     ActionRedact,
			Confidence: ConfidenceHigh,
			Reason:     "Redact inline kubectl tokens",
		},
	})
	if err != nil {
		t.Fatalf("MatchEntry returned error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("len(matches) = %d, want 2", len(matches))
	}
}

func TestMatchEntryRejectsInvalidEntry(t *testing.T) {
	_, err := MatchEntry(history.HistoryEntry{}, []Rule{
		{
			Name:       "drop-clear",
			Type:       RuleExact,
			Pattern:    "clear",
			Action:     ActionDelete,
			Confidence: ConfidenceHigh,
			Reason:     "Drop trivial clear commands",
		},
	})
	if err == nil {
		t.Fatal("MatchEntry returned nil error for invalid entry")
	}
}

func TestLargePasteBlobHeuristic(t *testing.T) {
	command := "cat <<'EOF'\n" + strings.Repeat("x", 520) + "\nEOF"

	match, ok, err := MatchRule(command, Rule{
		Name:       "large-paste-blob",
		Type:       RuleHeuristic,
		Detector:   "large_paste_blob",
		Action:     ActionQuarantine,
		Confidence: ConfidenceMedium,
		Reason:     "Quarantine likely accidental paste blobs",
	})
	if err != nil {
		t.Fatalf("MatchRule returned error: %v", err)
	}
	if !ok {
		t.Fatal("MatchRule returned ok=false, want true")
	}
	if match.RuleName != "large-paste-blob" {
		t.Fatalf("match.RuleName = %q, want large-paste-blob", match.RuleName)
	}
}
