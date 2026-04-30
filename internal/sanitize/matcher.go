package sanitize

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"histkit/internal/history"
)

type HeuristicFunc func(command string) bool

var heuristicDetectors = map[string]HeuristicFunc{
	"high_entropy_token": hasHighEntropyToken,
	"large_paste_blob":   hasLargePasteBlob,
}

func MatchEntry(entry history.HistoryEntry, rules []Rule) ([]RuleMatch, error) {
	if err := entry.Validate(); err != nil {
		return nil, fmt.Errorf("match entry: %w", err)
	}
	if err := ValidateRules(rules); err != nil {
		return nil, fmt.Errorf("match entry: %w", err)
	}

	matches := make([]RuleMatch, 0, len(rules))
	for _, rule := range rules {
		match, ok, err := MatchRule(entry.Command, rule)
		if err != nil {
			return nil, fmt.Errorf("match entry: %w", err)
		}
		if !ok {
			continue
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func MatchRule(command string, rule Rule) (RuleMatch, bool, error) {
	if strings.TrimSpace(command) == "" {
		return RuleMatch{}, false, fmt.Errorf("match rule: command is required")
	}
	if err := rule.Validate(); err != nil {
		return RuleMatch{}, false, fmt.Errorf("match rule: %w", err)
	}

	matched, err := matchesRule(command, rule)
	if err != nil {
		return RuleMatch{}, false, fmt.Errorf("match rule %q: %w", rule.Name, err)
	}
	if !matched {
		return RuleMatch{}, false, nil
	}

	match := RuleMatch{
		RuleName:   rule.Name,
		Reason:     rule.Reason,
		Confidence: rule.Confidence,
		Action:     rule.Action,
		Before:     command,
	}
	if rule.Action == ActionRedact {
		after, err := RedactCommand(command, rule)
		if err != nil {
			return RuleMatch{}, false, fmt.Errorf("apply redaction: %w", err)
		}
		match.After = after
	}

	return match, true, nil
}

func matchesRule(command string, rule Rule) (bool, error) {
	switch rule.Type {
	case RuleExact:
		return command == rule.Pattern, nil
	case RuleContains:
		return strings.Contains(command, rule.Pattern), nil
	case RuleRegex:
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return false, err
		}
		return re.MatchString(command), nil
	case RuleKeywordGroup:
		return keywordGroupMatches(command, rule.Keywords), nil
	case RuleHeuristic:
		detector, ok := heuristicDetectors[rule.Detector]
		if !ok {
			return false, fmt.Errorf("unknown heuristic detector %q", rule.Detector)
		}
		return detector(command), nil
	default:
		return false, fmt.Errorf("unsupported rule type %q", rule.Type)
	}
}

func hasHighEntropyToken(command string) bool {
	for _, token := range tokenize(command) {
		if isHighEntropyToken(token) {
			return true
		}
	}

	return false
}

func isHighEntropyToken(token string) bool {
	if len(token) < 24 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, r := range token {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func hasLargePasteBlob(command string) bool {
	return len(command) >= 500
}

func tokenize(command string) []string {
	return strings.FieldsFunc(command, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '/' || r == '+' || r == '=')
	})
}
