package log_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	asciiLog "ascii/ascii/src/internal/log"
)

func TestInitDefault(t *testing.T) {
	asciiLog.InitDefault()
	if asciiLog.Log == nil {
		t.Fatal("Log should not be nil after InitDefault")
	}
}

func TestInitVerbose(t *testing.T) {
	asciiLog.Init(true)
	if asciiLog.Log == nil {
		t.Fatal("Log should not be nil after Init(true)")
	}

	var buf bytes.Buffer
	asciiLog.Log = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	asciiLog.Log.Debug("test debug message")
	if !strings.Contains(buf.String(), "test debug message") {
		t.Fatalf("expected debug message in output, got: %s", buf.String())
	}
}

func TestInitNonVerbose(t *testing.T) {
	asciiLog.Init(false)
	if asciiLog.Log == nil {
		t.Fatal("Log should not be nil after Init(false)")
	}

	var buf bytes.Buffer
	asciiLog.Log = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	asciiLog.Log.Debug("should not appear")
	if buf.Len() != 0 {
		t.Fatalf("expected empty buffer, got: %s", buf.String())
	}
}

func TestInitWarnLevelVisible(t *testing.T) {
	asciiLog.Init(false)
	var buf bytes.Buffer
	asciiLog.Log = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	asciiLog.Log.Warn("warning message")
	if !strings.Contains(buf.String(), "warning message") {
		t.Fatalf("expected warn message in output, got: %s", buf.String())
	}
}
