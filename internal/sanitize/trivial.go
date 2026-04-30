package sanitize

import "histkit/internal/history"

func BuiltinTrivialRules() []Rule {
	return []Rule{
		{
			Name:       "clear-command",
			Type:       RuleExact,
			Pattern:    "clear",
			Action:     ActionDelete,
			Confidence: ConfidenceHigh,
			Reason:     "Drop trivial terminal clear commands",
		},
		{
			Name:       "pwd-command",
			Type:       RuleExact,
			Pattern:    "pwd",
			Action:     ActionDelete,
			Confidence: ConfidenceHigh,
			Reason:     "Drop trivial working-directory checks",
		},
		{
			Name:       "ls-command",
			Type:       RuleExact,
			Pattern:    "ls",
			Action:     ActionDelete,
			Confidence: ConfidenceMedium,
			Reason:     "Drop trivial directory listings",
		},
		{
			Name:       "ll-command",
			Type:       RuleExact,
			Pattern:    "ll",
			Action:     ActionDelete,
			Confidence: ConfidenceMedium,
			Reason:     "Drop trivial shell alias directory listings",
		},
		{
			Name:       "large-paste-blob",
			Type:       RuleHeuristic,
			Detector:   "large_paste_blob",
			Action:     ActionQuarantine,
			Confidence: ConfidenceMedium,
			Reason:     "Quarantine likely accidental large paste blobs",
		},
	}
}

func MatchTrivialRules(entry history.HistoryEntry) ([]RuleMatch, error) {
	return MatchEntry(entry, BuiltinTrivialRules())
}
