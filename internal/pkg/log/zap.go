package log

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger returns a pointer that points to a new instance of zapLogger.
func NewZapLogger(level Level) (Logger, error) {
	zapLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	switch level {
	case DebugLevel:
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case InfoLevel:
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case ErrorLevel:
		zapLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case PanicLevel:
		zapLevel = zap.NewAtomicLevelAt(zap.PanicLevel)
	case FatalLevel:
		zapLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
	}

	cfg := zap.Config{
		Level:       zapLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
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
