package log

import "log/slog"

type Config struct {
	Format Format `mapstructure:"format" validate:"required,oneof=json text"`
	Level  Level  `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Caller bool   `mapstructure:"caller"`
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
