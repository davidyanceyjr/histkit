package history

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func readHistoryLines(r io.Reader, onLine func(line string, lineNumber int) error) error {
	reader := bufio.NewReader(r)
	lineNumber := 0

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) == 0 && err == io.EOF {
			return nil
		}

		lineNumber++
		line = bytes.TrimSuffix(line, []byte{'\n'})
		line = bytes.TrimSuffix(line, []byte{'\r'})

		if cbErr := onLine(string(line), lineNumber); cbErr != nil {
			return cbErr
		}
		if err == io.EOF {
			return nil
		}
	}
}

func wrapHistoryReadError(shell, sourceFile string, err error) error {
	return fmt.Errorf("read %s history from %q: %w", shell, sourceFile, err)
}
