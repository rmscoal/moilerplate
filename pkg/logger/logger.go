package logger

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

type AppLogger struct {
	logger *logrus.Logger
	path   string
}

func NewAppLogger(path string) *AppLogger {
	logger := new(AppLogger)
	logger.path = path

	logger.logger = logrus.New()

	file, err := os.OpenFile(logger.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatalf("FATAL - Unable to initialize app logger: %s", err)
	}

	logger.logger.SetOutput(file)

	return logger
}

// Info -.
func (l *AppLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Warn -.
func (l *AppLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Error -.
func (l *AppLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Getter methods
func (a *AppLogger) Path() string {
	return a.path
}

func (a *AppLogger) Logger() *logrus.Logger {
	return a.logger
}
