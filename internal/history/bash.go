package history

import (
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

	if err := readHistoryLines(r, func(rawLine string, lineNumber int) error {
		switch {
		case rawLine == "":
			return nil
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
			return nil
		default:
			if err := onEntry(HistoryEntry{
				Shell:      ShellBash,
				SourceFile: sourceFile,
				RawLine:    rawLine,
				Command:    rawLine,
			}); err != nil {
				return err
			}
			return nil
		}
	}); err != nil {
		return wrapHistoryReadError("Bash", sourceFile, err)
	}

	return nil
}
