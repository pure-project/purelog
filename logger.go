package purelog

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

//logger instance
type Logger struct {
	config   *Config
	pid      int
	mtx      sync.Mutex
	wg       sync.WaitGroup
	bp       sync.Pool
	buf      *buffer
	buf2     *buffer
	files    []string
	flushCh  chan bool
	once     sync.Once
}

//new logger instance
func New(config ...*Config) (l *Logger) {
	 l = new(Logger)
	 l.init(config...)
	 return
}

//close logger
func (l *Logger) Close() {
	l.once.Do(func() {
		close(l.flushCh)
		l.wg.Wait()
		l.flush()
	})
}

//notify flush (async)
func (l *Logger) Flush() {
	l.flushCh <- true
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(LevelDebug, 1, "", args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(LevelInfo, 1, "", args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log(LevelWarn, 1, "", args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(LevelError, 1, "", args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(LevelDebug, 1, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(LevelInfo, 1, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(LevelWarn, 1, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(LevelError, 1, format, args...)
}

func (l *Logger) Log(level Level, skip int, args ...interface{}) {
	skip++
	l.log(level, skip, "", args...)
}

func (l *Logger) Logf(level Level, skip int, format string, args ...interface{}) {
	skip++
	l.log(level, skip, format, args...)
}


const (
	flushTimeMin   = 50 * time.Millisecond //minimum flush time: 50ms
	lineBufSize    = 1024                  //line buffer size: 1KB (for pool)
	fileBufSizeMin = 32 * 1024             //file buffer min size: 32KB (for double-buffer)
	fileBufSizeMax = 4 * 1024 * 1024       //file buffer max size:  4MB (for double-buffer)
)

func (l *Logger) init(config ...*Config) {
	if len(config) != 0 {
		l.config = config[0]
	} else {
		l.config = NewConfig().SetStderr(true).SetStdout(true).SetCaller(true)
	}

	l.pid     = os.Getpid()
	l.bp.New  = func() interface{} { return &buffer{ Data: make([]byte, 0, lineBufSize) } }
	l.buf     = &buffer{ Data: make([]byte, 0, fileBufSizeMin) }
	l.buf2    = &buffer{ Data: make([]byte, 0, fileBufSizeMin) }
	l.flushCh = make(chan bool, 1)
	l.wg.Add(1)
	go l.doLog()
}

func (l *Logger) enabled(level Level) bool {
	return l.config.getLevel() <= level && (l.config.getStdout() || len(l.config.getFile()) != 0)
}

func (l *Logger) log(level Level, skip int, format string, args ...interface{}) {
	if !l.enabled(level) {
		return
	}

	skip++
	file, line := l.caller(skip)
	_, file = reverseSplitN(file, 2, '/')

	levelStr := level.ShortString()

	if len(args) == 0 {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		l.buf.Data = appendHeader(l.buf.Data, time.Now(), l.pid, file, line, levelStr)  //for normal time line
		l.buf.Data = append(l.buf.Data, format...)
		l.buf.Data = append(l.buf.Data, '\n')
		return
	}

	if len(format) == 0 && len(args) == 1 {
		str, ok := args[0].(string)
		if ok {
			l.mtx.Lock()
			defer l.mtx.Unlock()
			l.buf.Data = appendHeader(l.buf.Data, time.Now(), l.pid, file, line, levelStr)
			l.buf.Data = append(l.buf.Data, str...)
			l.buf.Data = append(l.buf.Data, '\n')
			return
		}
	}

	buf := l.bp.Get().(*buffer)
	defer l.bp.Put(buf)

	buf.Reset()
	if len(format) == 0 {
		fmt.Fprint(buf, args...)
	} else {
		fmt.Fprintf(buf, format, args...)
	}

	buf.Data = append(buf.Data, '\n')

	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.buf.Data = appendHeader(l.buf.Data, time.Now(), l.pid, file, line, levelStr)
	l.buf.Data = append(l.buf.Data, buf.Data...)
}

func (l *Logger) doLog() {
	defer l.wg.Done()

	timer := time.NewTimer(maxDuration(flushTimeMin, l.config.getFlush()))
	defer timer.Stop()

	for {
		select {
		case _, ok := <-l.flushCh:
			if !ok {
				return
			}
			l.flush()

		case <-timer.C:
			l.flush()
		}

		timer.Reset(maxDuration(flushTimeMin, l.config.getFlush()))
	}
}

//flush log data
func (l *Logger) flush() {
	//swap double buffer
	l.mtx.Lock()
	l.buf, l.buf2 = l.buf2, l.buf
	l.mtx.Unlock()

	bufSize := l.buf2.Len()
	defer l.recycleMemory(bufSize, l.buf2.Cap())
	if bufSize == 0 {
		return
	}

	//output stdout
	if l.config.getStdout() {
		os.Stdout.Write(l.buf2.Data)
	}

	file := l.config.getFile()
	if len(file) == 0 {
		return
	}

	data := l.buf2.Data

	//rotate
	size := l.config.getSize()
	if size != 0 {
		data = l.rotate(size, file, data)
	}

	//sync to disk
	l.sync(file, data)
}

//sync write data to file
func (l *Logger) sync(file string, data []byte) {
	dir, _ := reverseSplitN(file, 1, '/')
	_ = os.MkdirAll(dir, os.ModePerm)

	out, err := os.OpenFile(file, os.O_CREATE | os.O_APPEND | os.O_WRONLY | os.O_SYNC, 0666)
	if err != nil {
		l.internalError("logger.sync: open log file %s err: %v", file, err)
		return
	}
	defer out.Close()

	_, err = out.Write(data)
	if err != nil {
		l.internalError("logger.sync: write log file %s err: %v", file, err)
		return
	}
	_ = out.Sync()
}

func (l *Logger) rotate(size uint64, file string, data []byte) []byte {
	for len(data) != 0 {
		fileSz := fileSize(file)
		if uint64(len(data)) + fileSz >= size {
			//can write size
			sz := size - fileSz
			//adjust line
			idx := reverseIndexB(data[:sz], 1, '\n')
			if idx != -1 {
				l.sync(file, data[:idx])
				data = data[idx+1:]
			} else {
				//adjust fail, direct cut
				l.sync(file, data[:sz])
				data = data[sz:]
			}
			//do rotate
			l.rotateFile(file)
			continue
		}

		break
	}
	return data
}

func (l *Logger) rotateFile(file string) {
	//gen new file name
	name, ext := reverseSplitN(file, 1, '.')
	var arr [128]byte
	buf := append(arr[:0], name...)
	buf = append(buf, '_')
	buf = appendRotateTime(buf, time.Now())
	buf = append(buf, '.')
	buf = append(buf, ext...)
	newFile := b2s(buf)

	//move file
	err := os.Rename(file, newFile)
	if err != nil {
		l.internalError("logger.rotate: move file %s to %s err: %v", file, newFile, err)
		os.Remove(file)
		return
	}

	//add to files
	l.files = append(l.files, newFile)

	//clean older
	l.clean()
}

func (l *Logger) clean() {
	count := l.config.getCount()
	if count != 0 && uint32(len(l.files)) >= count {
		file := l.files[0]
		err := os.Remove(file)
		if err != nil {
			l.internalError("logger.clean: remove log file %s err: %v", file, err)
		}
		l.files = l.files[1:]
	}
}

func (l *Logger) recycleMemory(size, cap int) {
	//recycle memory when most space unused.
	if fileBufSizeMax < cap && size <= fileBufSizeMin {
		//l.Debugf("recycle memory: size=%d cap=%d", size, cap)   //output memory recycle logging if need
		l.buf2.Data = make([]byte, 0, fileBufSizeMin)
		return
	}
	l.buf2.Reset()
}

func (l *Logger) caller(skip int) (string, int) {
	if l.config.getCaller() {
		skip++
		_, file, line, _ := runtime.Caller(skip)
		return file, line
	}
	return "???", 0
}

func (l *Logger) internalError(format string, args ...interface{}) {
	if l.config.getStderr() {
		defer os.Stderr.WriteString("\n")

		if len(args) == 0 {
			os.Stderr.WriteString(format)
			return
		}

		if len(format) == 0 {
			fmt.Fprint(os.Stderr, args...)

			return
		}

		fmt.Fprintf(os.Stderr, format, args...)
	}
}
