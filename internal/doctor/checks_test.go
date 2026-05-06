package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckSystemdUserUnitsNotConfigured(t *testing.T) {
	home := t.TempDir()

	result := checkSystemdUserUnits(home)
	if result.Status != StatusOK {
		t.Fatalf("expected OK status, got %q", result.Status)
	}
	if !strings.Contains(result.Detail, "systemd automation not configured") {
		t.Fatalf("expected not-configured detail, got %q", result.Detail)
	}
}

func TestCheckSystemdUserUnitsWarnsForPartialInstall(t *testing.T) {
	home := t.TempDir()
	unitDir := filepath.Join(home, ".config", "systemd", "user")
	writeFile(t, filepath.Join(unitDir, "histkit-scan.service"), "[Unit]\nDescription=test\n")

	result := checkSystemdUserUnits(home)
	if result.Status != StatusWarn {
		t.Fatalf("expected WARN status, got %q", result.Status)
	}
	if !strings.Contains(result.Detail, "partial systemd automation install") {
		t.Fatalf("expected partial-install detail, got %q", result.Detail)
	}
	if !strings.Contains(result.Detail, filepath.Join(unitDir, "histkit-scan.timer")) {
		t.Fatalf("expected missing timer path in detail, got %q", result.Detail)
	}
}

func TestCheckSystemdUserUnitsReportsPresentPair(t *testing.T) {
	home := t.TempDir()
	unitDir := filepath.Join(home, ".config", "systemd", "user")
	writeFile(t, filepath.Join(unitDir, "histkit-scan.service"), "[Unit]\nDescription=test\n")
	writeFile(t, filepath.Join(unitDir, "histkit-scan.timer"), "[Unit]\nDescription=test\n")

	result := checkSystemdUserUnits(home)
	if result.Status != StatusOK {
		t.Fatalf("expected OK status, got %q", result.Status)
	}
	if !strings.Contains(result.Detail, "histkit systemd user units present:") {
		t.Fatalf("expected present-pair detail, got %q", result.Detail)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
