package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/onsi/ginkgo/v2"
	slogmulti "github.com/samber/slog-multi"
)

func init() {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:       slog.LevelDebug,
		TimeFormat:  time.RFC3339,
		NoColor:     !isatty.IsTerminal(os.Stdout.Fd()),
		ReplaceAttr: replaceAttrTint,
	})

	logger := slog.New(
		slogmulti.
			Pipe(
				newAttrContextMiddleware(),
				newAttrTraceMiddleware(),
			).
			Handler(handler),
	)

	slog.SetDefault(logger)
}

func Init(cfg Config) {
	var (
		writer    = os.Stdout
		level     = cfg.Level.SlogLevel()
		addSource = cfg.Caller
	)

	var handler slog.Handler
	switch cfg.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level:       level,
			AddSource:   addSource,
			ReplaceAttr: replaceAttrJSON,
		})
	case FormatText:
		handler = tint.NewHandler(writer, &tint.Options{
			Level:       level,
			TimeFormat:  time.RFC3339,
			NoColor:     !isatty.IsTerminal(writer.Fd()),
			AddSource:   addSource,
			ReplaceAttr: replaceAttrTint,
		})
	}

	middlewares := []slogmulti.Middleware{
		newAttrContextMiddleware(),
		newAttrTraceMiddleware(),
	}
	if attrs := cfg.ConstAttrs(); len(attrs) > 0 {
		middlewares = append(middlewares, newAttrConstantMiddleware(attrs...))
	}

	handler = slogmulti.Pipe(middlewares...).Handler(handler)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func AdaptGinkgo() {
	logger := slog.New(tint.NewHandler(ginkgo.GinkgoWriter, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.RFC3339,
		NoColor:    !isatty.IsTerminal(os.Stdout.Fd()),
	}))
	slog.SetDefault(logger)
}
