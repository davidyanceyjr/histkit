package sanitize

import "testing"

func TestRuleValidateAcceptsSupportedRuleTypes(t *testing.T) {
	tests := []struct {
		name string
		rule Rule
	}{
		{
			name: "exact",
			rule: Rule{
				Name:       "drop-clear",
				Type:       RuleExact,
				Pattern:    "clear",
				Action:     ActionDelete,
				Confidence: ConfidenceHigh,
				Reason:     "Drop trivial clear commands",
			},
		},
		{
			name: "contains",
			rule: Rule{
				Name:       "private-key-marker",
				Type:       RuleContains,
				Pattern:    "BEGIN OPENSSH PRIVATE KEY",
				Action:     ActionQuarantine,
				Confidence: ConfidenceHigh,
				Reason:     "Quarantine private key material",
			},
		},
		{
			name: "regex",
			rule: Rule{
				Name:       "kubectl-token",
				Type:       RuleRegex,
				Pattern:    `kubectl.*--token[ =][^ ]+`,
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
				Reason:     "Redact inline kubectl tokens",
			},
		},
		{
			name: "keyword group",
			rule: Rule{
				Name:       "url-creds",
				Type:       RuleKeywordGroup,
				Keywords:   []string{"https://", "@", ":"},
				Action:     ActionQuarantine,
				Confidence: ConfidenceMedium,
				Reason:     "Quarantine likely credential-bearing URLs",
			},
		},
		{
			name: "heuristic",
			rule: Rule{
				Name:       "high-entropy-token",
				Type:       RuleHeuristic,
				Detector:   "high_entropy_token",
				Action:     ActionQuarantine,
				Confidence: ConfidenceMedium,
				Reason:     "Quarantine likely secret-like tokens",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.rule.Validate(); err != nil {
				t.Fatalf("Validate returned error: %v", err)
			}
		})
	}
}

func TestRuleValidateRejectsInvalidRules(t *testing.T) {
	tests := []struct {
		name string
		rule Rule
	}{
		{
			name: "missing name",
			rule: Rule{
				Type:       RuleContains,
				Pattern:    "token",
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
				Reason:     "desc",
			},
		},
		{
			name: "invalid type",
			rule: Rule{
				Name:       "bad-type",
				Type:       "substring",
				Pattern:    "token",
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
				Reason:     "desc",
			},
		},
		{
			name: "invalid action",
			rule: Rule{
				Name:       "bad-action",
				Type:       RuleContains,
				Pattern:    "token",
				Action:     "mask",
				Confidence: ConfidenceHigh,
				Reason:     "desc",
			},
		},
		{
			name: "invalid confidence",
			rule: Rule{
				Name:       "bad-confidence",
				Type:       RuleContains,
				Pattern:    "token",
				Action:     ActionRedact,
				Confidence: "critical",
				Reason:     "desc",
			},
		},
		{
			name: "missing reason",
			rule: Rule{
				Name:       "missing-reason",
				Type:       RuleContains,
				Pattern:    "token",
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
			},
		},
		{
			name: "missing pattern for regex",
			rule: Rule{
				Name:       "missing-pattern",
				Type:       RuleRegex,
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
				Reason:     "desc",
			},
		},
		{
			name: "invalid regex",
			rule: Rule{
				Name:       "invalid-regex",
				Type:       RuleRegex,
				Pattern:    "(",
				Action:     ActionRedact,
				Confidence: ConfidenceHigh,
				Reason:     "desc",
			},
		},
		{
			name: "missing keywords for keyword group",
			rule: Rule{
				Name:       "missing-keywords",
				Type:       RuleKeywordGroup,
				Action:     ActionQuarantine,
				Confidence: ConfidenceMedium,
				Reason:     "desc",
			},
		},
		{
			name: "blank keyword",
			rule: Rule{
				Name:       "blank-keyword",
				Type:       RuleKeywordGroup,
				Keywords:   []string{"token", " "},
				Action:     ActionQuarantine,
				Confidence: ConfidenceMedium,
				Reason:     "desc",
			},
		},
		{
			name: "missing detector for heuristic",
			rule: Rule{
				Name:       "missing-detector",
				Type:       RuleHeuristic,
				Action:     ActionQuarantine,
				Confidence: ConfidenceMedium,
				Reason:     "desc",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.rule.Validate(); err == nil {
				t.Fatal("Validate returned nil error")
			}
		})
	}
}

func TestValidateRulesRejectsDuplicateNames(t *testing.T) {
	rules := []Rule{
		{
			Name:       "dup",
			Type:       RuleContains,
			Pattern:    "token",
			Action:     ActionRedact,
			Confidence: ConfidenceHigh,
			Reason:     "first",
		},
		{
			Name:       "dup",
			Type:       RuleContains,
			Pattern:    "password",
			Action:     ActionQuarantine,
			Confidence: ConfidenceMedium,
			Reason:     "second",
		},
	}

	if err := ValidateRules(rules); err == nil {
		t.Fatal("ValidateRules returned nil error for duplicate names")
	}
}

func TestValidateRulesAcceptsDistinctRules(t *testing.T) {
	rules := []Rule{
		{
			Name:       "private-key-marker",
			Type:       RuleContains,
			Pattern:    "BEGIN OPENSSH PRIVATE KEY",
			Action:     ActionQuarantine,
			Confidence: ConfidenceHigh,
			Reason:     "Quarantine private key material",
		},
		{
			Name:       "redact-kubectl-token",
			Type:       RuleRegex,
			Pattern:    `kubectl.*--token[ =][^ ]+`,
			Action:     ActionRedact,
			Confidence: ConfidenceHigh,
			Reason:     "Redact inline kubectl tokens",
		},
	}

	if err := ValidateRules(rules); err != nil {
		t.Fatalf("ValidateRules returned error: %v", err)
	}
}

func TestRuleMatchValidate(t *testing.T) {
	match := RuleMatch{
		RuleName:   "redact-kubectl-token",
		Reason:     "Redact inline kubectl tokens",
		Confidence: ConfidenceHigh,
		Action:     ActionRedact,
		Before:     "kubectl --token abc",
		After:      "kubectl --token [REDACTED]",
	}

	if err := match.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestRuleMatchValidateRejectsInvalidMatches(t *testing.T) {
	tests := []struct {
		name  string
		match RuleMatch
	}{
		{
			name: "missing name",
			match: RuleMatch{
				Reason:     "desc",
				Confidence: ConfidenceHigh,
				Action:     ActionDelete,
				Before:     "echo hi",
			},
		},
		{
			name: "missing reason",
			match: RuleMatch{
				RuleName:   "name",
				Confidence: ConfidenceHigh,
				Action:     ActionDelete,
				Before:     "echo hi",
			},
		},
		{
			name: "invalid confidence",
			match: RuleMatch{
				RuleName:   "name",
				Reason:     "desc",
				Confidence: "critical",
				Action:     ActionDelete,
				Before:     "echo hi",
			},
		},
		{
			name: "invalid action",
			match: RuleMatch{
				RuleName:   "name",
				Reason:     "desc",
				Confidence: ConfidenceHigh,
				Action:     "mask",
				Before:     "echo hi",
			},
		},
		{
			name: "missing before",
			match: RuleMatch{
				RuleName:   "name",
				Reason:     "desc",
				Confidence: ConfidenceHigh,
				Action:     ActionDelete,
			},
		},
		{
			name: "redact missing after",
			match: RuleMatch{
				RuleName:   "name",
				Reason:     "desc",
				Confidence: ConfidenceHigh,
				Action:     ActionRedact,
				Before:     "secret",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.match.Validate(); err == nil {
				t.Fatal("Validate returned nil error")
			}
		})
	}
}
