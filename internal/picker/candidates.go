package picker

import (
	"database/sql"
	"fmt"

	"histkit/internal/index"
	"histkit/internal/snippets"
)

const (
	LabelHistory = "[history]"
	LabelSnippet = "[snippet]"
)

type Candidate struct {
	Label       string
	Command     string
	HistoryID   string
	SnippetID   string
	Title       string
	Description string
	Safety      string
	Tags        []string
}

func (c Candidate) Display() string {
	return c.Label + "  " + c.Command
}

func LoadCandidates(db *sql.DB, store snippets.Store, snippetsEnabled bool, includeBuiltins bool, historyLimit int) ([]Candidate, error) {
	historyEntries, err := index.QueryRecentHistoryEntries(db, historyLimit)
	if err != nil {
		return nil, err
	}

	userSnippets := []snippets.Snippet{}
	if snippetsEnabled {
		userSnippets, err = store.List()
		if err != nil {
			return nil, err
		}
	}

	mergedSnippets := mergeSnippets(userSnippets, snippetsEnabled && includeBuiltins)
	candidates := make([]Candidate, 0, len(historyEntries)+len(mergedSnippets))
	for _, entry := range historyEntries {
		candidates = append(candidates, Candidate{
			Label:     LabelHistory,
			Command:   entry.Command,
			HistoryID: entry.ID,
		})
	}
	for _, snippet := range mergedSnippets {
		candidates = append(candidates, Candidate{
			Label:       LabelSnippet,
			Command:     snippet.Command,
			SnippetID:   snippet.ID,
			Title:       snippet.Title,
			Description: snippet.Description,
			Safety:      snippet.Safety,
			Tags:        append([]string(nil), snippet.Tags...),
		})
	}

	return candidates, nil
}

func ParseSelectedLine(line string) (Candidate, error) {
	switch {
	case len(line) >= len(LabelHistory)+2 && line[:len(LabelHistory)] == LabelHistory:
		return Candidate{Label: LabelHistory, Command: line[len(LabelHistory)+2:]}, nil
	case len(line) >= len(LabelSnippet)+2 && line[:len(LabelSnippet)] == LabelSnippet:
		return Candidate{Label: LabelSnippet, Command: line[len(LabelSnippet)+2:]}, nil
	default:
		return Candidate{}, fmt.Errorf("parse selected line: unsupported candidate format")
	}
}

func mergeSnippets(user []snippets.Snippet, includeBuiltins bool) []snippets.Snippet {
	merged := append([]snippets.Snippet{}, user...)
	if !includeBuiltins {
		return merged
	}

	seen := make(map[string]struct{}, len(user))
	for _, snippet := range user {
		seen[snippet.ID] = struct{}{}
	}
	for _, snippet := range snippets.Builtins() {
		if _, ok := seen[snippet.ID]; ok {
			continue
		}
		merged = append(merged, snippet)
		seen[snippet.ID] = struct{}{}
	}

	return merged
}
