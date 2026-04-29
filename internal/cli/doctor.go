package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

func runDoctor(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("doctor: %w", err)
	}
	if *help {
		writeDoctorUsage(stdout)
		return nil
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("doctor: unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	fmt.Fprintln(stdout, "doctor is not implemented yet; basic environment checks will arrive in a later slice.")
	return nil
}

func writeDoctorUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit doctor")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Placeholder command: doctor currently reports only that checks are not implemented yet.")
}
