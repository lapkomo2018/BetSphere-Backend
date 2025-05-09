package log

import "github.com/sirupsen/logrus"

type (
	Logger      = logrus.Logger
	FieldLogger = logrus.FieldLogger
	Entry       = logrus.Entry
	Fields      = logrus.Fields
	Level       = logrus.Level
	Hook        = logrus.Hook

	Formatter     = logrus.Formatter
	TextFormatter = logrus.TextFormatter
	JsonFormatter = logrus.JSONFormatter
)

var (
	log       = logrus.New()
	AllLevels = logrus.AllLevels
)

func New() *Logger {
	return logrus.New()
}

func WithField(key string, value any) *Entry {
	return log.WithField(key, value)
}

func WithFields(fields Fields) *Entry {
	return log.WithFields(fields)
}

func WithError(err error) *Entry {
	return log.WithError(err)
}

func Debug(args ...any) {
	log.Debug(args...)
}

func Info(args ...any) {
	log.Info(args...)
}

func Error(args ...any) {
	log.Error(args...)
}

func Warn(args ...any) {
	log.Warn(args...)
}

func Fatal(args ...any) {
	log.Fatal(args...)
}

func SetFormatter(formatter Formatter) {
	log.SetFormatter(formatter)
}

func SetLevel(level Level) {
	log.SetLevel(level)
}

func AddHook(hook Hook) {
	log.AddHook(hook)
}
