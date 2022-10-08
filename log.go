package purelog

var DefaultConfig = NewConfig().
	SetStderr(true).
	SetStdout(true).
	SetCaller(true)

var DefaultLogger = New(DefaultConfig)

func Debug(args ...interface{}) {
	DefaultLogger.Log(LevelDebug, 1, args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Log(LevelInfo, 1, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Log(LevelWarn, 1, args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Log(LevelError, 1, args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Logf(LevelDebug, 1, format, args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Logf(LevelInfo, 1, format, args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Logf(LevelWarn, 1, format, args...)
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger.Logf(LevelError, 1, format, args...)
}

func Flush() {
	DefaultLogger.Flush()
}