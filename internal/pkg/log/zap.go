package log

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger returns a pointer that points to a new instance of zapLogger.
func NewZapLogger(development bool) (Logger, error) {
	var logger *zap.Logger
	var err error

	options := []zap.Option{
		zap.AddCallerSkip(2),
	}

	if development {
		logger, err = zap.NewDevelopment(options...)
	} else {
		logger, err = zap.NewProduction(options...)
	}
	if err != nil {
		return nil, err
	}

	return &zapLogger{
		logger: logger,
	}, nil
}

// Print TODO: add description
func (z *zapLogger) Print(level Level, msg string, fields Fields) {
	zFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zFields = append(zFields, zap.Any(k, v))
	}

	switch level {
	case DebugLevel:
		z.logger.Debug(msg, zFields...)
	case WarnLevel:
		z.logger.Warn(msg, zFields...)
	case ErrorLevel:
		z.logger.Error(msg, zFields...)
	case PanicLevel:
		z.logger.Panic(msg, zFields...)
	case FatalLevel:
		z.logger.Fatal(msg, zFields...)

	case InfoLevel:
		fallthrough
	default:
		z.logger.Info(msg, zFields...)
	}
}

// Debug ...
func (z *zapLogger) Debug(msg string, fields Fields) {
	z.Print(DebugLevel, msg, fields)
}

// Info ...
func (z *zapLogger) Info(msg string, fields Fields) {
	z.Print(InfoLevel, msg, fields)
}

// Warn ...
func (z *zapLogger) Warn(msg string, fields Fields) {
	z.Print(WarnLevel, msg, fields)
}

// Error ...
func (z *zapLogger) Error(msg string, fields Fields) {
	z.Print(ErrorLevel, msg, fields)
}

// Fatal ...
func (z *zapLogger) Fatal(msg string, fields Fields) {
	z.Print(FatalLevel, msg, fields)
}

// Panic ...
func (z *zapLogger) Panic(msg string, fields Fields) {
	z.Print(PanicLevel, msg, fields)
}
