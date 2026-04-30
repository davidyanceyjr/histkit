package sanitize

import (
	"strings"
	"testing"

	"histkit/internal/history"
)

func TestBuiltinTrivialRulesValidate(t *testing.T) {
	rules := BuiltinTrivialRules()
	if len(rules) == 0 {
		t.Fatal("BuiltinTrivialRules returned no rules")
	}
	if err := ValidateRules(rules); err != nil {
		t.Fatalf("ValidateRules returned error: %v", err)
	}
}

func TestMatchTrivialRulesTruePositives(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		ruleName string
		action   ActionType
	}{
		{
			name:     "clear command",
			command:  "clear",
			ruleName: "clear-command",
			action:   ActionDelete,
		},
		{
			name:     "pwd command",
			command:  "pwd",
			ruleName: "pwd-command",
			action:   ActionDelete,
		},
		{
			name:     "ls command",
			command:  "ls",
			ruleName: "ls-command",
			action:   ActionDelete,
		},
		{
			name:     "ll command",
			command:  "ll",
			ruleName: "ll-command",
			action:   ActionDelete,
		},
		{
			name:     "large paste blob",
			command:  "cat <<'EOF'\n" + strings.Repeat("x", 520) + "\nEOF",
			ruleName: "large-paste-blob",
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

			matches, err := MatchTrivialRules(entry)
			if err != nil {
				t.Fatalf("MatchTrivialRules returned error: %v", err)
			}
			if len(matches) == 0 {
				t.Fatal("MatchTrivialRules returned no matches")
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

func TestMatchTrivialRulesFalsePositiveGuards(t *testing.T) {
	tests := []string{
		"ls -la",
		"pwd && ls",
		"clear && reset",
		"kubectl get pods",
		"echo ll",
		"printf 'pwd\n'",
	}

	for _, command := range tests {
		t.Run(command, func(t *testing.T) {
			entry := history.HistoryEntry{
				Shell:      history.ShellBash,
				SourceFile: "/home/tester/.bash_history",
				RawLine:    command,
				Command:    command,
			}

			matches, err := MatchTrivialRules(entry)
			if err != nil {
				t.Fatalf("MatchTrivialRules returned error: %v", err)
			}
			if len(matches) != 0 {
				t.Fatalf("len(matches) = %d, want 0 for non-trivial command %#v", len(matches), matches)
			}
		})
	}
}
