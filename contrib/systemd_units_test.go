package contrib_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSystemdScanServiceTemplate(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(repoRoot(t), "contrib", "histkit-scan.service"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	want := strings.TrimSpace(`[Unit]
Description=Scan and index shell history with histkit

[Service]
Type=oneshot
ExecStart=%h/.local/bin/histkit scan --config %h/.config/histkit/config.toml`)

	if got := strings.TrimSpace(string(content)); got != want {
		t.Fatalf("service template mismatch:\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}

func TestSystemdScanTimerTemplate(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(repoRoot(t), "contrib", "histkit-scan.timer"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	want := strings.TrimSpace(`[Unit]
Description=Run histkit scan periodically

[Timer]
OnBootSec=5m
OnUnitActiveSec=12h
Persistent=true

[Install]
WantedBy=timers.target`)

	if got := strings.TrimSpace(string(content)); got != want {
		t.Fatalf("timer template mismatch:\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}
