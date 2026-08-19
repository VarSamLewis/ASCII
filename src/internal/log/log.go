package log

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init(verbose bool) {
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelDebug
	}

	Log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(Log)
}

func InitDefault() {
	Init(false)
}
