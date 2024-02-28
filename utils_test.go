package purelog

import (
	"testing"
	"time"
)

func Test_AppendInt(t *testing.T) {
	var buf [32]byte
	b := appendInt(buf[:0], 666)
	t.Logf("%s", b)
}

func Test_AppendHeader(t *testing.T) {
	var buf [32]byte
	b := appendHeader(buf[:0], time.Now(), 1, "test.go", 123, "ERR")
	t.Logf("%s", b)
}

func Test_ReverseSplitN(t *testing.T) {
	a, b := reverseSplitN("a/b/c", 1, '/')
	t.Logf("a=%s b=%s", a, b)

	a, b = reverseSplitN("a b c", 2, ' ')
	t.Logf("a=%s b=%s", a, b)
}
