package logger

import (
	"errors"
	"log/slog"
	"os"
)

func New(level, format string) (*slog.Logger, error) {
	var logger *slog.Logger
	logLevel, err := toLeveler(level)
	if err != nil {
		//slog.Error(err.Error())
		return nil, err
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,

		// Add context logging in debug mode
		AddSource: func(logLvl slog.Leveler) bool {
			return logLvl == slog.LevelDebug

		}(logLevel),
	}

	switch format {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	return logger, nil
}

func toLeveler(level string) (slog.Leveler, error) {
	switch level {
	case "error":
		return slog.LevelError, nil
	case "warn":
		return slog.LevelWarn, nil
	case "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	default:
		return slog.LevelError, errors.New("could not parse log level")
	}
}
