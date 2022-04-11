package log

type TestingLogger struct{}

func NewTestingLogger() Logger {
	return &TestingLogger{}
}

func (t TestingLogger) Print(level Level, msg string, fields Fields) {
}

func (t TestingLogger) Debug(msg string, fields Fields) {
}

func (t TestingLogger) Info(msg string, fields Fields) {
}

func (t TestingLogger) Warn(msg string, fields Fields) {
}

func (t TestingLogger) Error(msg string, fields Fields) {
}

func (t TestingLogger) Fatal(msg string, fields Fields) {
}

func (t TestingLogger) Panic(msg string, fields Fields) {
}
