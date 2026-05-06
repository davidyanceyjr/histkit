package history

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func ParseZsh(sourceFile string, r io.Reader) ([]HistoryEntry, []ParseWarning, error) {
	var entries []HistoryEntry
	var warnings []ParseWarning
	err := StreamZsh(
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

func StreamZsh(sourceFile string, r io.Reader, onEntry func(HistoryEntry) error, onWarning func(ParseWarning) error) error {
	if strings.TrimSpace(sourceFile) == "" {
		return fmt.Errorf("zsh parser source file is required")
	}
	if r == nil {
		return fmt.Errorf("zsh parser reader is required")
	}

	if err := readHistoryLines(r, func(rawLine string, lineNumber int) error {
		switch {
		case rawLine == "":
			return nil
		case strings.TrimSpace(rawLine) == "":
			if err := onWarning(ParseWarning{
				Shell:      ShellZsh,
				SourceFile: sourceFile,
				LineNumber: lineNumber,
				RawLine:    rawLine,
				Message:    "whitespace-only Zsh history line",
			}); err != nil {
				return err
			}
			return nil
		case strings.HasPrefix(rawLine, ": "):
			entry, warning := parseZshExtendedLine(sourceFile, lineNumber, rawLine)
			if warning != nil {
				if err := onWarning(*warning); err != nil {
					return err
				}
				return nil
			}
			if err := onEntry(entry); err != nil {
				return err
			}
			return nil
		default:
			if err := onEntry(HistoryEntry{
				Shell:      ShellZsh,
				SourceFile: sourceFile,
				RawLine:    rawLine,
				Command:    rawLine,
			}); err != nil {
				return err
			}
			return nil
		}
	}); err != nil {
		return wrapHistoryReadError("Zsh", sourceFile, err)
	}

	return nil
}

func parseZshExtendedLine(sourceFile string, lineNumber int, rawLine string) (HistoryEntry, *ParseWarning) {
	metadataAndCommand := strings.TrimPrefix(rawLine, ": ")
	parts := strings.SplitN(metadataAndCommand, ";", 2)
	if len(parts) != 2 {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: missing command separator",
		}
	}

	command := parts[1]
	if strings.TrimSpace(command) == "" {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: empty command",
		}
	}

	metadata := strings.SplitN(parts[0], ":", 2)
	if len(metadata) != 2 {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: invalid metadata fields",
		}
	}

	unixSeconds, err := strconv.ParseInt(metadata[0], 10, 64)
	if err != nil {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: invalid timestamp",
		}
	}

	if _, err := strconv.Atoi(metadata[1]); err != nil {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: invalid duration",
		}
	}

	timestamp := time.Unix(unixSeconds, 0).UTC()

	return HistoryEntry{
		Shell:      ShellZsh,
		SourceFile: sourceFile,
		RawLine:    rawLine,
		Command:    command,
		Timestamp:  &timestamp,
	}, nil
}
