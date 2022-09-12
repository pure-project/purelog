package purelog

import (
	"sync/atomic"
	"time"
)

type Config struct {
	file   atomic.Value
	size   uint64
	flush  uint64
	level  uint32
	stderr uint32
	stdout uint32
	count  uint32
	caller uint32
}

func NewConfig() *Config {
	c := &Config{}
	c.file.Store("")
	return c
}

func (c *Config) SetStderr(enb bool) *Config {
	atomic.StoreUint32(&c.stderr, bool2uint32(enb))
	return c
}

func (c *Config) SetStdout(enb bool) *Config {
	atomic.StoreUint32(&c.stdout, bool2uint32(enb))
	return c
}

func (c *Config) SetFile(file string) *Config {
	c.file.Store(file)
	return c
}

func (c *Config) SetLevel(level Level) *Config {
	atomic.StoreUint32(&c.level, uint32(level))
	return c
}

func (c *Config) SetSize(size uint) *Config {
	atomic.StoreUint64(&c.size, uint64(size))
	return c
}

func (c *Config) SetCount(count uint) *Config {
	atomic.StoreUint32(&c.count, uint32(count))
	return c
}

func (c *Config) SetCaller(enb bool) *Config {
	atomic.StoreUint32(&c.caller, bool2uint32(enb))
	return c
}

func (c *Config) SetFlush(flush time.Duration) *Config {
	atomic.StoreUint64(&c.flush, uint64(flush))
	return c
}



func (c *Config) getFile() string {
	iFile := c.file.Load()
	if iFile != nil {
		file, _ := iFile.(string)
		return file
	}
	return ""
}

func (c *Config) getStderr() bool {
	return atomic.LoadUint32(&c.stderr) != 0
}

func (c *Config) getStdout() bool {
	return atomic.LoadUint32(&c.stdout) != 0
}

func (c *Config) getLevel() Level {
	return Level(atomic.LoadUint32(&c.level))
}

func (c *Config) getSize() uint64 {
	return atomic.LoadUint64(&c.size)
}

func (c *Config) getCount() uint32 {
	return atomic.LoadUint32(&c.count)
}

func (c *Config) getCaller() bool {
	return atomic.LoadUint32(&c.caller) != 0
}

func (c *Config) getFlush() time.Duration {
	return time.Duration(atomic.LoadUint64(&c.flush))
}


func bool2uint32(b bool) uint32 {
	if b { return 1 }
	return 0
}

