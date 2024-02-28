package purelog

import "github.com/pure-project/purebuf"

//custom format function
//
//format maybe empty(etc. Info)
type FormatFunc func(buf *purebuf.Buffer, level Level, file string, line int, format string, args ...interface{})

//custom flush function
//
//return value means continue default-flush or not.
type FlushFunc func(data []byte) bool