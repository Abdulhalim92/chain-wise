package logger

import (
	"io"
	"log/slog"
	"os"
)

// Format — формат вывода логов.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Options — настройки логгера (уровень, формат, источник, имя сервиса).
type Options struct {
	Level     string
	Format    Format
	AddSource bool
	Service   string
	Output    io.Writer
}

// Logger — обёртка над slog с поддержкой опций и дочерних логгеров.
type Logger struct {
	*slog.Logger
	opts Options
}

// New создаёт логгер по опциям (уровень, формат json/text, add_source, service).
func New(opts Options) *Logger {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	if opts.Format == "" {
		opts.Format = FormatJSON
	}
	var lvl slog.Level
	switch opts.Level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	handlerOpts := &slog.HandlerOptions{Level: lvl, AddSource: opts.AddSource}
	var h slog.Handler
	if opts.Format == FormatText {
		h = slog.NewTextHandler(opts.Output, handlerOpts)
	} else {
		h = slog.NewJSONHandler(opts.Output, handlerOpts)
	}
	log := slog.New(h)
	if opts.Service != "" {
		log = log.With(slog.String("service", opts.Service))
	}
	return &Logger{Logger: log, opts: opts}
}

// WithService возвращает дочерний логгер с заданным именем сервиса.
func (l *Logger) WithService(service string) *Logger {
	child := l.Logger.With(slog.String("service", service))
	return &Logger{Logger: child, opts: l.opts}
}

// With возвращает дочерний логгер с добавленными атрибутами.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{Logger: l.Logger.With(args...), opts: l.opts}
}
