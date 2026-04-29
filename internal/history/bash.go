package history

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ParseBash(sourceFile string, r io.Reader) ([]HistoryEntry, []ParseWarning, error) {
	if strings.TrimSpace(sourceFile) == "" {
		return nil, nil, fmt.Errorf("bash parser source file is required")
	}
	if r == nil {
		return nil, nil, fmt.Errorf("bash parser reader is required")
	}

	scanner := bufio.NewScanner(r)
	var entries []HistoryEntry
	var warnings []ParseWarning
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		rawLine := scanner.Text()

		switch {
		case rawLine == "":
			continue
		case strings.TrimSpace(rawLine) == "":
			warnings = append(warnings, ParseWarning{
				Shell:      ShellBash,
				SourceFile: sourceFile,
				LineNumber: lineNumber,
				RawLine:    rawLine,
				Message:    "whitespace-only Bash history line",
			})
			continue
		default:
			entries = append(entries, HistoryEntry{
				Shell:      ShellBash,
				SourceFile: sourceFile,
				RawLine:    rawLine,
				Command:    rawLine,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("read Bash history from %q: %w", sourceFile, err)
	}

	return entries, warnings, nil
}
