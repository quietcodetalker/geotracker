package log

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"log"
)

type fluentdLogger struct {
	logger *fluent.Fluent
	tag    string
}

// NewFluentdLogger returns a pointer that points to a new instance of fluentdLogger.
func NewFluentdLogger(tag string, cfg fluent.Config) (Logger, error) {
	logger, err := fluent.New(cfg)
	if err != nil {
		return nil, err
	}

	return &fluentdLogger{
		logger: logger,
		tag:    tag,
	}, nil
}

// Print TODO: add description
func (z *fluentdLogger) Print(level Level, msg string, fields Fields) {
	if fields == nil {
		fields = Fields{}
	}
	fields["level"] = level.String()
	if err := z.logger.Post(msg, fields); err != nil {
		log.Printf("failed to send log to fluentd: %v", err)
	}
}

// Debug ...
func (z *fluentdLogger) Debug(msg string, fields Fields) {
	z.Print(DebugLevel, msg, fields)
}

// Info ...
func (z *fluentdLogger) Info(msg string, fields Fields) {
	z.Print(InfoLevel, msg, fields)
}

// Warn ...
func (z *fluentdLogger) Warn(msg string, fields Fields) {
	z.Print(WarnLevel, msg, fields)
}

// Error ...
func (z *fluentdLogger) Error(msg string, fields Fields) {
	z.Print(ErrorLevel, msg, fields)
}

// Fatal ...
func (z *fluentdLogger) Fatal(msg string, fields Fields) {
	z.Print(FatalLevel, msg, fields)
}

// Panic ...
func (z *fluentdLogger) Panic(msg string, fields Fields) {
	z.Print(PanicLevel, msg, fields)
}
