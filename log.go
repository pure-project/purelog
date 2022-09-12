package purelog

var DefaultConfig = NewConfig().
	SetStderr(true).
	SetStdout(true).
	SetCaller(true)

var DefaultLogger = New(DefaultConfig)

func Debug(args ...interface{}) {
	if DefaultLogger.enabled(LevelDebug) {
		DefaultLogger.log(LevelDebug, 2, "", args...)
	}
}

func Info(args ...interface{}) {
	if DefaultLogger.enabled(LevelInfo) {
		DefaultLogger.log(LevelInfo, 2, "", args...)
	}
}

func Warn(args ...interface{}) {
	if DefaultLogger.enabled(LevelWarn) {
		DefaultLogger.log(LevelWarn, 2, "", args...)
	}
}

func Error(args ...interface{}) {
	if DefaultLogger.enabled(LevelError) {
		DefaultLogger.log(LevelError, 2, "", args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if DefaultLogger.enabled(LevelDebug) {
		DefaultLogger.log(LevelDebug, 2, format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if DefaultLogger.enabled(LevelInfo) {
		DefaultLogger.log(LevelInfo, 2, format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if DefaultLogger.enabled(LevelWarn) {
		DefaultLogger.log(LevelWarn, 2, format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if DefaultLogger.enabled(LevelError) {
		DefaultLogger.log(LevelError, 2, format, args...)
	}
}

func Flush() {
	DefaultLogger.Flush()
}