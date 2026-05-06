package cli

import (
	"fmt"
	"io"
)

type helpBlock struct {
	Title string
	Lines []string
}

func writeHelpBlocks(w io.Writer, blocks ...helpBlock) {
	for i, block := range blocks {
		if i > 0 {
			fmt.Fprintln(w)
		}
		if block.Title != "" {
			fmt.Fprintf(w, "%s:\n", block.Title)
		}
		for _, line := range block.Lines {
			fmt.Fprintln(w, line)
		}
	}
}
