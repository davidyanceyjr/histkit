package sanitize

import (
	"fmt"
	"regexp"
	"strings"
)

const redactedPlaceholder = "[REDACTED]"

func RedactCommand(command string, rule Rule) (string, error) {
	if strings.TrimSpace(command) == "" {
		return "", fmt.Errorf("redact command: command is required")
	}
	if err := rule.Validate(); err != nil {
		return "", fmt.Errorf("redact command: %w", err)
	}

	switch rule.Type {
	case RuleExact:
		if command != rule.Pattern {
			return "", fmt.Errorf("redact command: exact rule %q did not match command", rule.Name)
		}
		return redactedPlaceholder, nil
	case RuleContains:
		if !strings.Contains(command, rule.Pattern) {
			return "", fmt.Errorf("redact command: contains rule %q did not match command", rule.Name)
		}
		return strings.ReplaceAll(command, rule.Pattern, redactedPlaceholder), nil
	case RuleRegex:
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return "", fmt.Errorf("redact command: %w", err)
		}
		if !re.MatchString(command) {
			return "", fmt.Errorf("redact command: regex rule %q did not match command", rule.Name)
		}
		return re.ReplaceAllString(command, redactedPlaceholder), nil
	case RuleKeywordGroup:
		if !keywordGroupMatches(command, rule.Keywords) {
			return "", fmt.Errorf("redact command: keyword-group rule %q did not match command", rule.Name)
		}
		redacted := command
		for _, keyword := range rule.Keywords {
			redacted = strings.ReplaceAll(redacted, keyword, redactedPlaceholder)
		}
		return redacted, nil
	case RuleHeuristic:
		redacted, ok, err := redactHeuristic(command, rule.Detector)
		if err != nil {
			return "", fmt.Errorf("redact command: %w", err)
		}
		if !ok {
			return "", fmt.Errorf("redact command: heuristic rule %q did not match command", rule.Name)
		}
		return redacted, nil
	default:
		return "", fmt.Errorf("redact command: unsupported rule type %q", rule.Type)
	}
}

func redactHeuristic(command, detector string) (string, bool, error) {
	switch detector {
	case "high_entropy_token":
		redacted := command
		matched := false
		for _, token := range tokenize(command) {
			if !isHighEntropyToken(token) {
				continue
			}
			replacement := redactedPlaceholder
			if key, value, ok := strings.Cut(token, "="); ok && isHighEntropyToken(value) {
				replacement = key + "=" + redactedPlaceholder
			}
			redacted = strings.ReplaceAll(redacted, token, replacement)
			matched = true
		}
		return redacted, matched, nil
	case "large_paste_blob":
		if !hasLargePasteBlob(command) {
			return "", false, nil
		}
		return redactedPlaceholder, true, nil
	default:
		return "", false, fmt.Errorf("unknown heuristic detector %q", detector)
	}
}

func keywordGroupMatches(command string, keywords []string) bool {
	for _, keyword := range keywords {
		if !strings.Contains(command, keyword) {
			return false
		}
	}
	return true
}
