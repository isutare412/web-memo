package log

import "log/slog"

type Config struct {
	Format      Format `koanf:"format" validate:"required,oneof=json text"`
	Level       Level  `koanf:"level" validate:"required,oneof=debug info warn error"`
	Caller      bool   `koanf:"caller"`
	Component   string `koanf:"component"`   // Optional
	Environment string `koanf:"environment"` // Optional
	Version     string `koanf:"version"`     // Optional
}

func (c Config) ConstAttrs() []slog.Attr {
	var attrs []slog.Attr
	if c.Component != "" {
		attrs = append(attrs, slog.String("component", c.Component))
	}
	if c.Environment != "" {
		attrs = append(attrs, slog.String("environment", c.Environment))
	}
	if c.Version != "" {
		attrs = append(attrs, slog.String("version", c.Version))
	}
	return attrs
}

type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

func (l Level) SlogLevel() slog.Level {
	sl := slog.LevelInfo
	switch l {
	case LevelDebug:
		sl = slog.LevelDebug
	case LevelInfo:
		sl = slog.LevelInfo
	case LevelWarn:
		sl = slog.LevelWarn
	case LevelError:
		sl = slog.LevelError
	}
	return sl
}
