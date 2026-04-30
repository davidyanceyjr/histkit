package sanitize

import "histkit/internal/history"

func BuiltinSecretRules() []Rule {
	return []Rule{
		{
			Name:       "private-key-block",
			Type:       RuleContains,
			Pattern:    "BEGIN OPENSSH PRIVATE KEY",
			Action:     ActionQuarantine,
			Confidence: ConfidenceHigh,
			Reason:     "Quarantine pasted private key material",
		},
		{
			Name:       "bearer-token",
			Type:       RuleRegex,
			Pattern:    `(?i)\bbearer\s+[A-Za-z0-9._\-+/=]{12,}`,
			Action:     ActionRedact,
			Confidence: ConfidenceHigh,
			Reason:     "Redact bearer tokens",
		},
		{
			Name:       "inline-password-flag",
			Type:       RuleRegex,
			Pattern:    `(?i)(--password|--passwd)(=| )[^\s]+|\bpassword=[^\s]+|\s-p\s*[^\s]+`,
			Action:     ActionRedact,
			Confidence: ConfidenceHigh,
			Reason:     "Redact inline password values",
		},
		{
			Name:       "url-embedded-credentials",
			Type:       RuleRegex,
			Pattern:    `(?i)\bhttps?://[^/\s:@]+:[^/\s@]+@[^/\s]+`,
			Action:     ActionRedact,
			Confidence: ConfidenceHigh,
			Reason:     "Redact URL-embedded credentials",
		},
		{
			Name:       "aws-access-key-id",
			Type:       RuleRegex,
			Pattern:    `\b(AKIA|ASIA)[A-Z0-9]{16}\b`,
			Action:     ActionRedact,
			Confidence: ConfidenceHigh,
			Reason:     "Redact cloud access key identifiers",
		},
		{
			Name:       "high-entropy-token",
			Type:       RuleHeuristic,
			Detector:   "high_entropy_token",
			Action:     ActionQuarantine,
			Confidence: ConfidenceMedium,
			Reason:     "Quarantine likely secret-like high-entropy tokens",
		},
	}
}

func MatchSecretRules(entry history.HistoryEntry) ([]RuleMatch, error) {
	return MatchEntry(entry, BuiltinSecretRules())
}
