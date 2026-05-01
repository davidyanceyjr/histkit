package sanitize

import (
	"fmt"
	"slices"
	"strings"

	"histkit/internal/history"
)

type ApplyReport struct {
	SourceFile         string
	Shell              string
	TotalLines         int
	ParsedEntries      int
	MatchedEntries     int
	RewrittenEntries   int
	DeletedEntries     int
	CountsByAction     map[ActionType]int
	CountsByConfidence map[Confidence]int
	RuleNames          []string
	RewrittenContent   []byte
}

func ApplyToSource(source history.Source, input []byte) (ApplyReport, error) {
	if _, err := history.ParserForShell(source.Shell); err != nil {
		return ApplyReport{}, fmt.Errorf("apply source: %w", err)
	}
	if strings.TrimSpace(source.Path) == "" {
		return ApplyReport{}, fmt.Errorf("apply source: source path is required")
	}

	report := ApplyReport{
		SourceFile:         source.Path,
		Shell:              source.Shell,
		CountsByAction:     make(map[ActionType]int),
		CountsByConfidence: make(map[Confidence]int),
	}

	lines, hasTrailingNewline := splitHistoryLines(string(input))
	outputLines := make([]string, 0, len(lines))
	ruleNames := map[string]struct{}{}

	for _, line := range lines {
		report.TotalLines++

		entry, ok, err := parseHistoryLine(source, line)
		if err != nil {
			return ApplyReport{}, fmt.Errorf("apply source: %w", err)
		}
		if !ok {
			outputLines = append(outputLines, line)
			continue
		}

		report.ParsedEntries++

		matches, err := matchBuiltinRules(entry)
		if err != nil {
			return ApplyReport{}, fmt.Errorf("apply source: %w", err)
		}
		if len(matches) == 0 {
			outputLines = append(outputLines, line)
			continue
		}

		report.MatchedEntries++
		for _, match := range matches {
			report.CountsByAction[match.Action]++
			report.CountsByConfidence[match.Confidence]++
			ruleNames[match.RuleName] = struct{}{}
		}

		switch finalAction(matches) {
		case ActionDelete:
			report.DeletedEntries++
		case ActionQuarantine:
			rewrittenLine, err := rewriteHistoryLine(entry, quarantinedPlaceholder)
			if err != nil {
				return ApplyReport{}, fmt.Errorf("apply source: %w", err)
			}
			outputLines = append(outputLines, rewrittenLine)
			report.RewrittenEntries++
		case ActionRedact:
			redacted, err := redactBuiltinRules(entry.Command)
			if err != nil {
				return ApplyReport{}, fmt.Errorf("apply source: %w", err)
			}
			rewrittenLine, err := rewriteHistoryLine(entry, redacted)
			if err != nil {
				return ApplyReport{}, fmt.Errorf("apply source: %w", err)
			}
			outputLines = append(outputLines, rewrittenLine)
			report.RewrittenEntries++
		default:
			outputLines = append(outputLines, line)
		}
	}

	report.RuleNames = sortedRuleNames(ruleNames)
	report.RewrittenContent = joinHistoryLines(outputLines, hasTrailingNewline)
	return report, nil
}

func matchBuiltinRules(entry history.HistoryEntry) ([]RuleMatch, error) {
	secretMatches, err := MatchSecretRules(entry)
	if err != nil {
		return nil, err
	}
	trivialMatches, err := MatchTrivialRules(entry)
	if err != nil {
		return nil, err
	}

	return append(secretMatches, trivialMatches...), nil
}

func finalAction(matches []RuleMatch) ActionType {
	action := ActionKeep
	for _, match := range matches {
		switch match.Action {
		case ActionQuarantine:
			return ActionQuarantine
		case ActionRedact:
			action = ActionRedact
		case ActionDelete:
			if action == ActionKeep {
				action = ActionDelete
			}
		}
	}
	return action
}

func redactBuiltinRules(command string) (string, error) {
	redacted := command
	for _, rule := range BuiltinSecretRules() {
		if rule.Action != ActionRedact {
			continue
		}

		match, ok, err := MatchRule(redacted, rule)
		if err != nil {
			return "", err
		}
		if ok {
			redacted = match.After
		}
	}
	for _, rule := range BuiltinTrivialRules() {
		if rule.Action != ActionRedact {
			continue
		}

		match, ok, err := MatchRule(redacted, rule)
		if err != nil {
			return "", err
		}
		if ok {
			redacted = match.After
		}
	}

	return redacted, nil
}

func rewriteHistoryLine(entry history.HistoryEntry, command string) (string, error) {
	if err := entry.Validate(); err != nil {
		return "", fmt.Errorf("rewrite history line: %w", err)
	}

	switch entry.Shell {
	case history.ShellBash:
		return command, nil
	case history.ShellZsh:
		if strings.HasPrefix(entry.RawLine, ": ") {
			parts := strings.SplitN(strings.TrimPrefix(entry.RawLine, ": "), ";", 2)
			if len(parts) != 2 {
				return "", fmt.Errorf("rewrite history line: malformed zsh extended history line")
			}
			return ": " + parts[0] + ";" + command, nil
		}
		return command, nil
	default:
		return "", fmt.Errorf("rewrite history line: unsupported shell %q", entry.Shell)
	}
}

func parseHistoryLine(source history.Source, line string) (history.HistoryEntry, bool, error) {
	switch source.Shell {
	case history.ShellBash:
		if line == "" || strings.TrimSpace(line) == "" {
			return history.HistoryEntry{}, false, nil
		}
		return history.HistoryEntry{
			Shell:      history.ShellBash,
			SourceFile: source.Path,
			RawLine:    line,
			Command:    line,
		}, true, nil
	case history.ShellZsh:
		switch {
		case line == "":
			return history.HistoryEntry{}, false, nil
		case strings.TrimSpace(line) == "":
			return history.HistoryEntry{}, false, nil
		case strings.HasPrefix(line, ": "):
			return parseZshHistoryLine(source.Path, line)
		default:
			return history.HistoryEntry{
				Shell:      history.ShellZsh,
				SourceFile: source.Path,
				RawLine:    line,
				Command:    line,
			}, true, nil
		}
	default:
		return history.HistoryEntry{}, false, fmt.Errorf("unsupported shell %q", source.Shell)
	}
}

func parseZshHistoryLine(sourceFile, rawLine string) (history.HistoryEntry, bool, error) {
	metadataAndCommand := strings.TrimPrefix(rawLine, ": ")
	parts := strings.SplitN(metadataAndCommand, ";", 2)
	if len(parts) != 2 {
		return history.HistoryEntry{}, false, nil
	}
	if strings.TrimSpace(parts[1]) == "" {
		return history.HistoryEntry{}, false, nil
	}

	return history.HistoryEntry{
		Shell:      history.ShellZsh,
		SourceFile: sourceFile,
		RawLine:    rawLine,
		Command:    parts[1],
	}, true, nil
}

func splitHistoryLines(content string) ([]string, bool) {
	if content == "" {
		return nil, false
	}

	hasTrailingNewline := strings.HasSuffix(content, "\n")
	trimmed := strings.TrimSuffix(content, "\n")
	if trimmed == "" {
		return []string{""}, hasTrailingNewline
	}

	return strings.Split(trimmed, "\n"), hasTrailingNewline
}

func joinHistoryLines(lines []string, hasTrailingNewline bool) []byte {
	if len(lines) == 0 {
		if hasTrailingNewline {
			return []byte("\n")
		}
		return nil
	}

	content := strings.Join(lines, "\n")
	if hasTrailingNewline {
		content += "\n"
	}
	return []byte(content)
}

func sortedRuleNames(ruleNames map[string]struct{}) []string {
	names := make([]string, 0, len(ruleNames))
	for name := range ruleNames {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
