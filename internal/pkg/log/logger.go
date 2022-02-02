package log

// Level represents a logging priority.
type Level int32

const (
	// DebugLevel TODO: add description
	DebugLevel = iota - 1
	// InfoLevel is the default logging level.
	InfoLevel
	// WarnLevel TODO: add description
	WarnLevel
	// ErrorLevel TODO: add description
	ErrorLevel
	// PanicLevel TODO: add description
	PanicLevel
	// FatalLevel TODO: add description
	FatalLevel
)

// Fields is a map where keys and values represent names and values of log structure fields.
type Fields map[string]interface{}

// Logger is the interface that wraps logging methods.
type Logger interface {
	Print(level Level, msg string, fields Fields)
	Debug(msg string, fields Fields)
	Info(msg string, fields Fields)
	Warn(msg string, fields Fields)
	Error(msg string, fields Fields)
	Fatal(msg string, fields Fields)
	Panic(msg string, fields Fields)
}
