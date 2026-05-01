package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"histkit/internal/audit"
	"histkit/internal/backup"
	"histkit/internal/config"
	"histkit/internal/history"
	"histkit/internal/sanitize"
)

func runRestore(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("restore", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")
	configPath := fs.String("config", "", "path to a histkit config file")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("restore: %w", err)
	}
	if *help {
		writeRestoreUsage(stdout)
		return nil
	}
	if fs.NArg() > 1 {
		return fmt.Errorf("restore: unexpected arguments: %s", strings.Join(fs.Args()[1:], " "))
	}

	backupID := ""
	if fs.NArg() == 1 {
		backupID = fs.Arg(0)
	}

	return executeRestore(stdout, restoreOptions{
		ConfigPath: *configPath,
		BackupID:   backupID,
	})
}

func writeRestoreUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit restore [backup-id] [--config <path>]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "List available backup records or restore a specific backup by identifier.")
}

type restoreOptions struct {
	ConfigPath string
	BackupID   string
}

func executeRestore(stdout io.Writer, opts restoreOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("restore: detect home directory: %w", err)
	}

	if opts.ConfigPath != "" {
		expandedConfigPath, err := config.ExpandUserPath(opts.ConfigPath, home)
		if err != nil {
			return fmt.Errorf("restore: %w", err)
		}
		if _, err := config.Load(expandedConfigPath); err != nil {
			return fmt.Errorf("restore: %w", err)
		}
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	backupDir := filepath.Join(paths.StateDir, "backups")
	if opts.BackupID == "" {
		return listBackups(stdout, backupDir)
	}

	record, err := backup.FindRecord(backupDir, opts.BackupID)
	if err != nil {
		return fmt.Errorf("restore: %w", err)
	}
	if err := backup.Restore(record); err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	startedAt := time.Now().UTC()
	auditRecord := audit.Record{
		RunID:              "restore_" + startedAt.Format("20060102T150405.000000000Z"),
		StartedAt:          startedAt,
		CompletedAt:        startedAt,
		Shell:              inferShell(record.SourceFile),
		RuleNames:          []string{"restore"},
		CountsByAction:     map[sanitize.ActionType]int{},
		CountsByConfidence: map[sanitize.Confidence]int{},
		BackupID:           record.ID,
		Apply:              false,
	}
	if err := audit.Append(paths.AuditLog, auditRecord); err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	fmt.Fprintf(stdout, "restore complete: backup=%s source=%s\n", record.ID, record.SourceFile)
	return nil
}

func listBackups(stdout io.Writer, backupDir string) error {
	records, err := backup.ListRecords(backupDir)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		fmt.Fprintln(stdout, "restore: no backups available")
		return nil
	}

	for _, record := range records {
		fmt.Fprintf(
			stdout,
			"backup=%s created_at=%s source=%s checksum=%s\n",
			record.ID,
			record.CreatedAt.UTC().Format(time.RFC3339),
			record.SourceFile,
			record.Checksum,
		)
	}
	return nil
}

func inferShell(sourceFile string) string {
	switch filepath.Base(sourceFile) {
	case ".bash_history":
		return history.ShellBash
	case ".zsh_history":
		return history.ShellZsh
	default:
		return filepath.Base(sourceFile)
	}
}
