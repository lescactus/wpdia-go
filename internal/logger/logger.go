package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger(level, format string) (*logrus.Logger, error) {
	log := logrus.New()

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Ensure the given loglevel is valid for Logrus
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	log.SetLevel(l)

	// Set the log formatter
	switch format {
	case "text":
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	// Add context logging in debug mode
	if log.Level == logrus.DebugLevel {
		log.SetReportCaller(true)
	}

	return log, nil
}
