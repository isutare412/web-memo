package log

import (
	"log/slog"

	"github.com/lmittmann/tint"
)

var SlogLevelAccess = slog.Level(1)

func replaceAttrJSON(groups []string, a slog.Attr) slog.Attr {
	// Handle access log level
	if a.Key == slog.LevelKey && len(groups) == 0 {
		if level, ok := a.Value.Any().(slog.Level); ok {
			if level == SlogLevelAccess {
				a.Value = slog.StringValue("ACCESS")
			}
		}
	}

	return a
}

func replaceAttrTint(groups []string, a slog.Attr) slog.Attr {
	// Handle access, debug log level
	if a.Key == slog.LevelKey && len(groups) == 0 {
		if level, ok := a.Value.Any().(slog.Level); ok {
			switch level {
			case SlogLevelAccess:
				// Color access log level as cyan
				a = tint.Attr(6, slog.String(a.Key, "ACC"))
			case slog.LevelDebug:
				// Color debug log level as magenta
				a = tint.Attr(4, slog.String(a.Key, "DBG"))
			}
		}
	}

	// Handle error values
	if a.Value.Kind() == slog.KindAny {
		if _, ok := a.Value.Any().(error); ok {
			// Color error attribute as red
			a = tint.Attr(1, a)
		}
	}

	return a
}
