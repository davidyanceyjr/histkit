package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

func runScan(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	if *help {
		writeScanUsage(stdout)
		return nil
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("scan: unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	fmt.Fprintln(stdout, "scan is not implemented yet; no history files were modified.")
	return nil
}

func writeScanUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit scan")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Placeholder command: scan is not implemented yet and performs no mutation.")
}
