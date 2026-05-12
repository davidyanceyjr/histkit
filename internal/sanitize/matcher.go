package sanitize

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/davidyanceyjr/histkit/internal/history"
)

type HeuristicFunc func(command string) bool

var heuristicDetectors = map[string]HeuristicFunc{
	"high_entropy_token":   hasHighEntropyToken,
	"inline_password_flag": hasInlinePasswordFlag,
	"large_paste_blob":     hasLargePasteBlob,
}

var (
	inlinePasswordEqualsPattern = regexp.MustCompile(`(?i)^--passw(?:ord|d)=(\S+)$`)
	inlinePasswordValuePattern  = regexp.MustCompile(`(?i)^--passw(?:ord|d)$`)
	mysqlShortPasswordPattern   = regexp.MustCompile(`^-p(\S+)$`)
)

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
	_, matched := redactHighEntropySensitiveValues(command)
	return matched
}

func hasInlinePasswordFlag(command string) bool {
	_, matched := redactInlinePasswordFlags(command)
	return matched
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

func redactInlinePasswordFlags(command string) (string, bool) {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return "", false
	}

	redacted := append([]string(nil), fields...)
	commandName := effectiveCommandName(fields)
	matched := false

	for i, field := range fields {
		switch {
		case inlinePasswordEqualsPattern.MatchString(field):
			key, _, _ := strings.Cut(field, "=")
			redacted[i] = key + "=" + redactedPlaceholder
			matched = true
		case inlinePasswordValuePattern.MatchString(field):
			if i+1 >= len(fields) {
				continue
			}
			redacted[i+1] = redactedPlaceholder
			matched = true
		case isMySQLShortPasswordField(field, commandName):
			redacted[i] = "-p" + redactedPlaceholder
			matched = true
		}
	}

	return strings.Join(redacted, " "), matched
}

func redactHighEntropySensitiveValues(command string) (string, bool) {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return "", false
	}

	redacted := append([]string(nil), fields...)
	matched := false

	for i := 0; i < len(fields); i++ {
		targetIndex, replacement, ok := highEntropyReplacement(fields, i)
		if !ok {
			continue
		}
		redacted[targetIndex] = replacement
		matched = true
	}

	return strings.Join(redacted, " "), matched
}

func highEntropyReplacement(fields []string, index int) (int, string, bool) {
	field := fields[index]

	key, value, ok := strings.Cut(field, "=")
	if ok && isSensitiveKeyName(key) && isHighEntropyCandidate(value) {
		return index, key + "=" + redactedPlaceholder, true
	}

	if !strings.HasPrefix(field, "--") {
		return 0, "", false
	}
	flag := strings.TrimPrefix(field, "--")
	if flag == "" {
		return 0, "", false
	}

	if key, value, ok := strings.Cut(flag, "="); ok && isSensitiveKeyName(key) && isHighEntropyCandidate(value) {
		return index, "--" + key + "=" + redactedPlaceholder, true
	}

	if !isSensitiveKeyName(flag) || index+1 >= len(fields) || !isHighEntropyCandidate(fields[index+1]) {
		return 0, "", false
	}

	return index + 1, redactedPlaceholder, true
}

func isHighEntropyCandidate(value string) bool {
	if !isHighEntropyToken(value) {
		return false
	}
	if strings.ContainsAny(value, `/\`) {
		return false
	}
	if strings.HasPrefix(value, "~") || strings.HasPrefix(value, ".") {
		return false
	}
	return true
}

func isSensitiveKeyName(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))
	key = strings.TrimLeft(key, "-")
	if key == "" {
		return false
	}

	for _, fragment := range []string{"password", "passwd", "pwd", "token", "secret", "credential", "auth", "session"} {
		if strings.Contains(key, fragment) {
			return true
		}
	}

	if key == "key" || strings.Contains(key, "api_key") || strings.Contains(key, "access_key") ||
		strings.Contains(key, "secret_key") || strings.HasSuffix(key, "_key") || strings.HasSuffix(key, "-key") {
		return true
	}

	return false
}

func effectiveCommandName(fields []string) string {
	for _, field := range fields {
		if looksLikeEnvAssignment(field) {
			continue
		}
		return filepath.Base(field)
	}
	return ""
}

func looksLikeEnvAssignment(field string) bool {
	if strings.HasPrefix(field, "-") {
		return false
	}
	key, _, ok := strings.Cut(field, "=")
	if !ok || key == "" {
		return false
	}
	for i, r := range key {
		if i == 0 {
			if !(unicode.IsLetter(r) || r == '_') {
				return false
			}
			continue
		}
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return false
		}
	}
	return true
}

func isMySQLShortPasswordField(field, commandName string) bool {
	if !mysqlShortPasswordPattern.MatchString(field) {
		return false
	}

	switch strings.ToLower(commandName) {
	case "mysql", "mysqldump", "mysqladmin", "mariadb":
		return true
	default:
		return false
	}
}

func tokenize(command string) []string {
	return strings.FieldsFunc(command, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '/' || r == '+' || r == '=')
	})
}
