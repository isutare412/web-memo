package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.RFC3339,
		NoColor:    !isatty.IsTerminal(os.Stdout.Fd()),
	}))
	slog.SetDefault(logger)
}

func Init(cfg Config) {
	var (
		writer    = os.Stdout
		level     = cfg.Level.SlogLevel()
		addSource = cfg.Caller
	)

	var logger *slog.Logger
	switch cfg.Format {
	case FormatJSON:
		logger = slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level:     level,
			AddSource: addSource,
		}))
	case FormatText:
		logger = slog.New(tint.NewHandler(writer, &tint.Options{
			Level:      level,
			TimeFormat: time.RFC3339,
			NoColor:    !isatty.IsTerminal(writer.Fd()),
			AddSource:  addSource,
		}))
	}

	slog.SetDefault(logger)
}
