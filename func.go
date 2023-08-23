package purelog

//custom format function
type FormatFunc func(buf []byte, level Level, file string, line int, format string, args ...interface{}) []byte

//custom flush function
type FlushFunc func(data []byte)