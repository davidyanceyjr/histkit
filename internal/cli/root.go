package cli

import (
	"fmt"
	"io"
)

func Execute(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		writeRootUsage(stdout)
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		return runHelp(args[1:], stdout)
	case "scan":
		return runScan(args[1:], stdout)
	case "pick":
		return runPick(args[1:], stdout)
	case "stats":
		return runStats(args[1:], stdout)
	case "doctor":
		return runDoctor(args[1:], stdout)
	default:
		writeRootUsage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runHelp(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		writeRootUsage(stdout)
		return nil
	}

	switch args[0] {
	case "scan":
		writeScanUsage(stdout)
	case "pick":
		writePickUsage(stdout)
	case "stats":
		writeStatsUsage(stdout)
	case "doctor":
		writeDoctorUsage(stdout)
	default:
		return fmt.Errorf("unknown help topic %q", args[0])
	}

	return nil
}

func writeRootUsage(w io.Writer) {
	fmt.Fprintln(w, "histkit is a conservative CLI for shell history hygiene.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit <command> [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  scan    Parse history sources and update the local index")
	fmt.Fprintln(w, "  pick    Select a command from indexed history and snippets")
	fmt.Fprintln(w, "  stats   Show indexed history statistics")
	fmt.Fprintln(w, "  doctor  Check the local histkit environment")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Use \"histkit help <command>\" for command-specific help.")
}
