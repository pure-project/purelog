package purelog

import (
	"os"
	"reflect"
	"time"
	"unsafe"
)

//1949-10-01 07:00:00.000000 pid level file line |
func appendHeader(buf []byte, now time.Time, pid int, file string, line int, level string) []byte {
	buf = appendTimestamp(buf, now)
	buf = append(buf, ' ')
	buf = appendInt(buf, pid)
	buf = append(buf, ' ')
	buf = append(buf, level...)
	buf = append(buf, ' ')
	buf = append(buf, file...)
	buf = append(buf, ':')
	buf = appendInt(buf, line)
	return append(buf, " | "...)
}

//1949-10-01 07:00:00.000000
func appendTimestamp(buf []byte, now time.Time) []byte {
	y, m ,d := now.Date()
	h, min, s := now.Clock()
	us := now.Nanosecond() / 1000

	buf = appendInt0(buf, y, 4)
	buf = append(buf, '-')
	buf = appendInt0(buf, int(m), 2)
	buf = append(buf, '-')
	buf = appendInt0(buf, d, 2)
	buf = append(buf, ' ')
	buf = appendInt0(buf, h, 2)
	buf = append(buf, ':')
	buf = appendInt0(buf, min, 2)
	buf = append(buf, ':')
	buf = appendInt0(buf, s, 2)
	buf = append(buf, '.')
	return appendInt0(buf, us, 6)
}

//1949-10-01_07-10-59_000000000
func appendRotateTime(buf []byte, now time.Time) []byte {
	y, m, d := now.Date()
	h, min, s := now.Clock()
	ns := now.Nanosecond()

	buf = appendInt0(buf, y, 4)
	buf = append(buf, '-')
	buf = appendInt0(buf, int(m), 2)
	buf = append(buf, '-')
	buf = appendInt0(buf, d, 2)
	buf = append(buf, '_')
	buf = appendInt0(buf, h, 2)
	buf = append(buf, '-')
	buf = appendInt0(buf, min, 2)
	buf = append(buf, '-')
	buf = appendInt0(buf, s, 2)
	buf = append(buf, '_')
	return appendInt0(buf, ns, 9)
}

//int to string
func appendInt(buf []byte, num int) []byte {
	if num < 10 {
		return append(buf, '0' + byte(num))
	}

	if num < 20 {
		buf  = append(buf, '0' + byte(num / 10))
		return append(buf, '0' + byte(num % 10))
	}

	//static buffer
	var arr [32]byte
	b := arr[:0]

	//append to string
	for num != 0 {
		b = append(b, '0' + byte(num % 10))
		num /= 10
	}

	//reverse
	for i := len(b) - 1; i >= 0; i-- {
		buf = append(buf, b[i])
	}

	return buf
}

//int to string pad zero
func appendInt0(buf []byte, num, count int) []byte {
	if num < 10 {
		buf  = append(buf, '0')
		return append(buf, '0' + byte(num))
	}

	if num < 20 {
		buf  = append(buf, '0' + byte(num / 10))
		return append(buf, '0' + byte(num % 10))
	}

	//static buffer
	var arr [32]byte
	b := arr[:0]

	//append to string
	for num != 0 {
		b = append(b, '0' + byte(num % 10))
		num /= 10
	}

	//reverse
	for i := 0; i < count - len(b); i++ {
		buf = append(buf, '0')
	}

	//pad zeros
	for i := minInt(len(b) - 1, count); i >= 0; i-- {
		buf = append(buf, b[i])
	}

	return buf
}

func minInt(a, b int) int {
	if a < b { return a }
	return b
}

//returns max time.Duration
func maxDuration(a, b time.Duration) time.Duration {
	if a < b { return b }
	return a
}

func reverseIndex(s string, n int, c byte) int {
	return reverseIndexB(s2b(s), n, c)
}

func reverseIndexB(b []byte, n int, c byte) int {
	count := 0
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] == c {
			count++
			if count == n {
				return i
			}
		}
	}
	return -1
}


//reverse split string
func reverseSplitN(s string, n int, c byte) (string, string) {
	i := reverseIndex(s, n, c)
	if i != -1 {
		return s[:i], s[i + 1:]
	}
	return "", s
}

//fast []byte to string
func b2s(b []byte) string {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: hdr.Data, Len: hdr.Len}))
}

//fast string to []byte
func s2b(s string) []byte {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: hdr.Data, Len: hdr.Len, Cap: hdr.Len}))
}

//get file size
func fileSize(file string) uint64 {
	stat, err := os.Stat(file)
	if err == nil {
		return uint64(stat.Size())
	}
	return 0
}