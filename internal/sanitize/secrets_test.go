package sanitize

import (
	"testing"

	"histkit/internal/history"
)

func TestBuiltinSecretRulesValidate(t *testing.T) {
	rules := BuiltinSecretRules()
	if len(rules) == 0 {
		t.Fatal("BuiltinSecretRules returned no rules")
	}
	if err := ValidateRules(rules); err != nil {
		t.Fatalf("ValidateRules returned error: %v", err)
	}
}

func TestMatchSecretRulesTruePositives(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		ruleName string
		action   ActionType
	}{
		{
			name:     "private key block",
			command:  "cat <<'EOF'\nBEGIN OPENSSH PRIVATE KEY\nabc\nEOF",
			ruleName: "private-key-block",
			action:   ActionQuarantine,
		},
		{
			name:     "bearer token",
			command:  "curl -H 'Authorization: Bearer abcdefghijklmnopqrstuvwx123456' https://example.com",
			ruleName: "bearer-token",
			action:   ActionRedact,
		},
		{
			name:     "inline password flag",
			command:  "mysql --password hunter2",
			ruleName: "inline-password-flag",
			action:   ActionRedact,
		},
		{
			name:     "url embedded credentials",
			command:  "curl https://user:secret@example.com/path",
			ruleName: "url-embedded-credentials",
			action:   ActionRedact,
		},
		{
			name:     "aws access key id",
			command:  "export AWS_ACCESS_KEY_ID=AKIAABCDEFGHIJKLMNOP",
			ruleName: "aws-access-key-id",
			action:   ActionRedact,
		},
		{
			name:     "high entropy token",
			command:  "export TOKEN=AbcdefGhijklmnopQRST12uvwxYZ34",
			ruleName: "high-entropy-token",
			action:   ActionQuarantine,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			entry := history.HistoryEntry{
				Shell:      history.ShellBash,
				SourceFile: "/home/tester/.bash_history",
				RawLine:    tc.command,
				Command:    tc.command,
			}

			matches, err := MatchSecretRules(entry)
			if err != nil {
				t.Fatalf("MatchSecretRules returned error: %v", err)
			}
			if len(matches) == 0 {
				t.Fatal("MatchSecretRules returned no matches")
			}

			var found bool
			for _, match := range matches {
				if match.RuleName == tc.ruleName {
					found = true
					if match.Action != tc.action {
						t.Fatalf("match.Action = %q, want %q", match.Action, tc.action)
					}
				}
			}
			if !found {
				t.Fatalf("expected rule %q in matches %#v", tc.ruleName, matches)
			}
		})
	}
}

func TestMatchSecretRulesFalsePositiveGuards(t *testing.T) {
	tests := []string{
		"ssh prod-box",
		"sudo systemctl status nginx",
		"openssl version",
		"kubectl get pods",
		"curl https://example.com/path",
		"echo password rotation complete",
	}

	for _, command := range tests {
		t.Run(command, func(t *testing.T) {
			entry := history.HistoryEntry{
				Shell:      history.ShellBash,
				SourceFile: "/home/tester/.bash_history",
				RawLine:    command,
				Command:    command,
			}

			matches, err := MatchSecretRules(entry)
			if err != nil {
				t.Fatalf("MatchSecretRules returned error: %v", err)
			}
			if len(matches) != 0 {
				t.Fatalf("len(matches) = %d, want 0 for non-secret command %#v", len(matches), matches)
			}
		})
	}
}

func TestSecretRuleRedactionsProduceMaskedOutput(t *testing.T) {
	entry := history.HistoryEntry{
		Shell:      history.ShellBash,
		SourceFile: "/home/tester/.bash_history",
		RawLine:    "curl -H 'Authorization: Bearer abcdefghijklmnopqrstuvwx123456' https://example.com",
		Command:    "curl -H 'Authorization: Bearer abcdefghijklmnopqrstuvwx123456' https://example.com",
	}

	matches, err := MatchSecretRules(entry)
	if err != nil {
		t.Fatalf("MatchSecretRules returned error: %v", err)
	}

	for _, match := range matches {
		if match.Action != ActionRedact {
			continue
		}
		if match.After == "" || match.After == match.Before {
			t.Fatalf("redact match.After = %q, want transformed masked output", match.After)
		}
	}
}
