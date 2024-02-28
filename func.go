package purelog

//custom format function
type FormatFunc func(buf []byte, level Level, file string, line int, format string, args ...interface{}) []byte

//custom flush function
//
//return value means continue default-flush or not.
type FlushFunc func(data []byte) bool