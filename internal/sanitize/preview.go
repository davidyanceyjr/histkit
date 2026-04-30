package sanitize

import (
	"fmt"
	"strings"

	"histkit/internal/history"
)

type PreviewItem struct {
	Entry   history.HistoryEntry
	Matches []RuleMatch
}

type PreviewReport struct {
	Items          []PreviewItem
	TotalEntries   int
	MatchedEntries int
	TotalMatches   int
	CountsByAction map[ActionType]int
}

func PreviewEntry(entry history.HistoryEntry) (PreviewItem, error) {
	if err := entry.Validate(); err != nil {
		return PreviewItem{}, fmt.Errorf("preview entry: %w", err)
	}

	secretMatches, err := MatchSecretRules(entry)
	if err != nil {
		return PreviewItem{}, fmt.Errorf("preview entry: %w", err)
	}
	trivialMatches, err := MatchTrivialRules(entry)
	if err != nil {
		return PreviewItem{}, fmt.Errorf("preview entry: %w", err)
	}

	matches := append(secretMatches, trivialMatches...)
	return PreviewItem{
		Entry:   entry,
		Matches: matches,
	}, nil
}

func PreviewEntries(entries []history.HistoryEntry) (PreviewReport, error) {
	report := PreviewReport{
		Items:          []PreviewItem{},
		TotalEntries:   len(entries),
		CountsByAction: make(map[ActionType]int),
	}

	for _, entry := range entries {
		item, err := PreviewEntry(entry)
		if err != nil {
			return PreviewReport{}, fmt.Errorf("preview entries: %w", err)
		}
		if len(item.Matches) == 0 {
			continue
		}

		report.Items = append(report.Items, item)
		report.MatchedEntries++
		report.TotalMatches += len(item.Matches)
		for _, match := range item.Matches {
			report.CountsByAction[match.Action]++
		}
	}

	return report, nil
}

func RenderPreviewText(report PreviewReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "dry-run preview: %d matched entrie(s), %d total match(es), %d scanned entrie(s)\n",
		report.MatchedEntries,
		report.TotalMatches,
		report.TotalEntries,
	)

	b.WriteString("counts by action:\n")
	for _, action := range orderedActions() {
		fmt.Fprintf(&b, "  %s: %d\n", action, report.CountsByAction[action])
	}

	if len(report.Items) == 0 {
		b.WriteString("no matches\n")
		return b.String()
	}

	for _, item := range report.Items {
		fmt.Fprintf(&b, "entry: shell=%s source=%s\n", item.Entry.Shell, item.Entry.SourceFile)
		fmt.Fprintf(&b, "  original: %s\n", item.Entry.Command)
		for _, match := range item.Matches {
			fmt.Fprintf(&b, "  match: rule=%s confidence=%s action=%s\n", match.RuleName, match.Confidence, match.Action)
			if match.Action == ActionRedact {
				fmt.Fprintf(&b, "  transformed: %s\n", match.After)
			}
			fmt.Fprintf(&b, "  reason: %s\n", match.Reason)
		}
	}

	return b.String()
}

func orderedActions() []ActionType {
	return []ActionType{
		ActionKeep,
		ActionDelete,
		ActionRedact,
		ActionQuarantine,
	}
}
