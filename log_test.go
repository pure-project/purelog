package purelog

import (
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)


//logger:

func TestLog(t *testing.T) {
	Debug(1, 2, 3, 4)
	Debugf("format %d %d %d %d", 1, 2, 3, 4)
	Info(1, 2, 3, 4)
	Infof("format %d %d %d %d", 1, 2, 3, 4)
	Warn(1, 2, 3, 4)
	Warnf("format %d %d %d %d", 1, 2, 3, 4)
	Error(1, 2, 3, 4)
	Errorf("format %d %d %d %d", 1, 2, 3, 4)
	time.Sleep(time.Second)
}

func TestLog2(t *testing.T) {
	logger := New(NewConfig().
		SetStdout(true).
		SetFile("test.log").
		SetCaller(true))
	defer logger.Close()

	for i := 0; i < 100; i++ {
		logger.Infof("test %d", i)
	}

	time.Sleep(time.Second)
}

func TestLogFile(t *testing.T) {
	logger := New(NewConfig().
		SetCaller(true).
		SetStderr(true).
		SetFile("test.log").
		SetFlush(100 * time.Millisecond))
	defer logger.Close()

	logger.Debug(1, 2, 3, 4)
	logger.Debugf("format %d %d %d %d", 1, 2, 3, 4)
	logger.Info(1, 2, 3, 4)
	logger.Infof("format %d %d %d %d", 1, 2, 3, 4)
	logger.Warn(1, 2, 3, 4)
	logger.Warnf("format %d %d %d %d", 1, 2, 3, 4)
	logger.Error(1, 2, 3, 4)
	logger.Errorf("format %d %d %d %d", 1, 2, 3, 4)

	time.Sleep(time.Second)
}

func TestMultiLog(t *testing.T) {
	logger := New(NewConfig().
		SetCaller(true).
		SetStderr(true).
		SetStdout(true).
		SetFlush(100 * time.Millisecond))
	defer logger.Close()

	const (
		goroutineCount = 1000
		singleRunCount = 1000
	)

	var wg sync.WaitGroup
	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func(i int) {
			defer wg.Done()
			beg := time.Now()
			for j := 0; j < singleRunCount; j++ {
				logger.Infof("test i=%d j=%d", i, j)
			}
			logger.Debugf("%d cost %v", i, time.Since(beg))
		}(i)
	}
	wg.Wait()
}

func TestMultiLogFile(t *testing.T) {
	logger := New(NewConfig().
		SetCaller(true).
		SetStderr(true).
		SetFile("test.log").
		SetFlush(100 * time.Millisecond))
	defer logger.Close()

	const (
		goroutineCount = 1000
		singleRunCount = 1000
	)

	var wg sync.WaitGroup
	wg.Add(goroutineCount)
	var count uint64
	for i := 0; i < goroutineCount; i++ {
		go func(i int) {
			defer wg.Done()
			beg := time.Now()
			for j := 0; j < singleRunCount; j++ {
				logger.Info("simple count=", atomic.AddUint64(&count, 1))
			}
			logger.Debugf("%d cost %v", i, time.Since(beg))
		}(i)
	}
	wg.Wait()

	time.Sleep(time.Second)

}

func TestMultiLogFileRecycle(t *testing.T) {
	logger := New(NewConfig().
		SetCaller(true).
		SetStderr(true).
		SetFile("test.log").
		SetFlush(100 * time.Millisecond))
	defer logger.Close()

	const (
		goroutineCount = 100
		singleRunCount = 100
	)

	var wg sync.WaitGroup
	wg.Add(goroutineCount)
	var count uint64
	for i := 0; i < goroutineCount; i++ {
		go func(i int) {
			defer wg.Done()
			beg := time.Now()
			for j := 0; j < singleRunCount; j++ {
				logger.Info("simple count=", atomic.AddUint64(&count, 1))
			}
			logger.Debugf("%d cost %v", i, time.Since(beg))
		}(i)
	}
	wg.Wait()

	//wait recycle
	time.Sleep(5 * time.Second)

	//run again
	wg.Add(goroutineCount)
	count = 0
	for i := 0; i < goroutineCount; i++ {
		go func(i int) {
			defer wg.Done()
			beg := time.Now()
			for j := 0; j < singleRunCount; j++ {
				logger.Info("simple count=", atomic.AddUint64(&count, 1))
			}
			logger.Debugf("%d cost %v", i, time.Since(beg))
		}(i)
	}
	wg.Wait()

}

func TestRotateBasic(t *testing.T) {
	logger := New(NewConfig().
		SetFile("test.log").
		SetSize(1024).
		SetCount(5).
		SetCaller(true))

	defer logger.Close()

	var count uint64
	for i := 0; i < 100; i++ {
		logger.Info(atomic.AddUint64(&count, 1))
	}

	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		logger.Info(atomic.AddUint64(&count, 1))
	}
}

func TestRotateSmaller(t *testing.T) {
	logger := New(NewConfig().
		SetFile("test.log").
		SetSize(20).
		SetCount(5).
		SetCaller(true))

	defer logger.Close()

	s300 := func() string {
		buf := [300]byte{}
		for i := range buf {
			buf[i] = ' '
		}
		return string(buf[:])
	}()

	for i := 0; i < 10; i++ {
		logger.Info(s300)
	}
}

func TestChangeConfig(t *testing.T) {
	config := NewConfig().
		SetLevel(LevelInfo).
		SetStdout(true)
	logger := New(config)
	defer logger.Close()

	go func() {
		for {
			logger.Infof("test")
			time.Sleep(100 * time.Millisecond)
		}
	}()

	sleep := func() {
		time.Sleep(2 * time.Second)
	}

	sleep()
	config.SetLevel(LevelWarn)
	sleep()
	config.SetLevel(LevelInfo)
	sleep()
	config.SetCaller(true)
	sleep()
	config.SetFile("test.log")
	sleep()
}

func TestError(t *testing.T) {
	config := NewConfig().
		SetStderr(true)
	logger := New(config)
	defer logger.Close()

	config.SetFile("./test/")
	logger.Infof("test fail")
	time.Sleep(time.Second)
}

func TestFlush(t *testing.T) {
	logger := New(NewConfig().
		SetStdout(true).
		SetFlush(1000 * time.Second))
	defer logger.Close()

	logger.Info("test1")
	logger.Warn("test2")
	logger.Error("test3")

	t.Log("begin flush")
	logger.Flush()
	t.Log("end flush")

	time.Sleep(time.Second)
}


//benchmarks:


func BenchmarkStdLog(b *testing.B) {
	file, _ := os.OpenFile("test.log", os.O_CREATE | os.O_APPEND | os.O_WRONLY | os.O_SYNC, 0666)
	logger := log.New(file, "", log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//logger.Printf("test %d", i)
		logger.Print("simple")
	}
}

func BenchmarkLog(b *testing.B) {
	logger := New(NewConfig().
		SetFile("test.log").
		SetCaller(true))
	//defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("simple")
	}
}

func BenchmarkNormal(b *testing.B) {
	logger := New(NewConfig().
		SetFile("test.log").
		SetCaller(true))
	//defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("test %d", i)
	}
}

func BenchmarkFast(b *testing.B) {
	logger := New(NewConfig().
		SetFile("test.log"))
	//defer logger.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("simple")
	}
}

func BenchmarkSlow(b *testing.B) {
	logger := New(NewConfig().
		SetFile("test.log").
		SetCaller(true))
	//defer logger.Close()

	st := struct {
		A int
		B string
		C float64
		D map[string]interface{}
	} {
		A: 1,
		B: "b",
		C: 3.14,
		D: map[string]interface{} {
			"E": 4,
			"F": 666,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("test %+v", st)
	}
}


//utils:

func TestReverseSplitN(t *testing.T) {
	t.Log(reverseSplitN("",           2, '/'))
	t.Log(reverseSplitN("a.go",       2, '/'))
	t.Log(reverseSplitN("/a.go",      2, '/'))
	t.Log(reverseSplitN("a/b.go",     2, '/'))
	t.Log(reverseSplitN("a/b/c.go",   2, '/'))
	t.Log(reverseSplitN("a/b/c/d.go", 2, '/'))

	t.Log(reverseSplitN("a.1",     1, '.'))
	t.Log(reverseSplitN("a.1.2",   1, '.'))
	t.Log(reverseSplitN("a.1.2.3", 1, '.'))
}

func TestAppendHeader(t *testing.T) {
	t.Logf("%s\n", appendHeader(nil, time.Now(), 99, "test.go", 666, "ERR"))
}

func TestAppendInt(t *testing.T) {
	t.Logf("%s\n", appendInt(nil, 1234))
	t.Logf("%s\n", appendInt0(nil, 1234, 9))
}

func BenchmarkStdAppendInt(b *testing.B) {
	var arr [32]byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.AppendInt(arr[:0], 9,     10)
		strconv.AppendInt(arr[:0], 99,    10)
		strconv.AppendInt(arr[:0], 999,   10)
		strconv.AppendInt(arr[:0], 9999,  10)
		strconv.AppendInt(arr[:0], 99999, 10)
	}
}

func BenchmarkAppendInt(b *testing.B) {
	var arr [32]byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		appendInt(arr[:0], 9    )
		appendInt(arr[:0], 99   )
		appendInt(arr[:0], 999  )
		appendInt(arr[:0], 9999 )
		appendInt(arr[:0], 99999)
	}
}