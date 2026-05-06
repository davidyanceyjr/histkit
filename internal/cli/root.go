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
	case "clean":
		return runClean(args[1:], stdout)
	case "restore":
		return runRestore(args[1:], stdout)
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
	case "clean":
		writeCleanUsage(stdout)
	case "restore":
		writeRestoreUsage(stdout)
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
	fmt.Fprintln(w, "histkit is a conservative CLI for shell history hygiene, reusable snippets, and fuzzy command recall.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "It keeps raw shell history, the local history index, and snippets separate by design.")
	fmt.Fprintln(w, "The normal workflow is: doctor -> scan -> stats or pick -> clean --dry-run -> clean --apply -> restore.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit <command> [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  doctor  Check config, local paths, detected history sources, fzf, and related environment state")
	fmt.Fprintln(w, "  scan    Parse supported shell history sources and refresh the local SQLite history index")
	fmt.Fprintln(w, "  stats   Show indexed history counts by shell and source")
	fmt.Fprintln(w, "  pick    Select a command from indexed history and snippets through fzf")
	fmt.Fprintln(w, "  clean   Preview or apply cleanup rules to shell history with backups and audit logging")
	fmt.Fprintln(w, "  restore List recorded backups or restore a specific backup")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Use \"histkit help <command>\" or \"histkit <command> --help\" for command-specific help.")
}
