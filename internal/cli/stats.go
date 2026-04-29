package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

func runStats(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("stats", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("stats: %w", err)
	}
	if *help {
		writeStatsUsage(stdout)
		return nil
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("stats: unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	fmt.Fprintln(stdout, "stats are not available yet; the history index has not been implemented.")
	return nil
}

func writeStatsUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit stats")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Placeholder command: stats are not available until indexing exists.")
}
