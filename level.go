package purelog

type Level uint32
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func ParseLevel(level string) Level {
	switch level {
	case "debug", "DBG":
		return LevelDebug
	case "info", "INF":
		return LevelInfo
	case "warn", "WAR":
		return LevelWarn
	case "error", "ERR":
		return LevelError
	default:
		return LevelDebug
	}
}

func (level Level) String() string {
	switch level {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "debug"
	}
}

func (level Level) shortString() string {
	switch level {
	case LevelDebug:
		return "DBG"
	case LevelInfo:
		return "INF"
	case LevelWarn:
		return "WAR"
	case LevelError:
		return "ERR"
	default:
		return "DBG"
	}
}