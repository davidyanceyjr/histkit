package history

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ParseBash(sourceFile string, r io.Reader) ([]HistoryEntry, []ParseWarning, error) {
	var entries []HistoryEntry
	var warnings []ParseWarning
	err := StreamBash(
		sourceFile,
		r,
		func(entry HistoryEntry) error {
			entries = append(entries, entry)
			return nil
		},
		func(warning ParseWarning) error {
			warnings = append(warnings, warning)
			return nil
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return entries, warnings, nil
}

func StreamBash(sourceFile string, r io.Reader, onEntry func(HistoryEntry) error, onWarning func(ParseWarning) error) error {
	if strings.TrimSpace(sourceFile) == "" {
		return fmt.Errorf("bash parser source file is required")
	}
	if r == nil {
		return fmt.Errorf("bash parser reader is required")
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		rawLine := scanner.Text()

		switch {
		case rawLine == "":
			continue
		case strings.TrimSpace(rawLine) == "":
			if err := onWarning(ParseWarning{
				Shell:      ShellBash,
				SourceFile: sourceFile,
				LineNumber: lineNumber,
				RawLine:    rawLine,
				Message:    "whitespace-only Bash history line",
			}); err != nil {
				return err
			}
			continue
		default:
			if err := onEntry(HistoryEntry{
				Shell:      ShellBash,
				SourceFile: sourceFile,
				RawLine:    rawLine,
				Command:    rawLine,
			}); err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read Bash history from %q: %w", sourceFile, err)
	}

	return nil
}
