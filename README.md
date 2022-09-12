# pure-log

**A pure logging library for go, simple & fast.**

[中文](doc/README_cn.md)



### Overview

pure-log is a simple asynchronous logging library for golang.

although small, it is not a toy and is fully production-ready.



### Features

1. asynchronous logging
2. dynamic change configuration
3. rotate file and size/count control support



### Usage

simple:

```go
//do not forget close the default logging to flush log data
defer purelog.DefaultLogger.Close()

//same to fmt.Print
purelog.Debug("debug message.")
purelog.Info("info message.")
purelog.Warn("warn message.")
purelog.Error("error message.")

//same to fmt.Printf
purelog.Debugf("formatted message: %d or %v", 1234, 5678)

//set default logging output level
purelog.DefaultConfig.SetLevel(purelog.LevelWarn)

purelog.Info("this message can't be output!")
purelog.Warn("this message can be output.")

//set default logging output to file
purelog.DefaultConfig.SetFile("test.log")

//request flush default logging data
purelog.Flush()
```

outputs:

```
//Format: Date Time.Micro PID File:Line Level | Message\n
2022-09-12 23:49:51.542323 16180 purelog_demo/main.go:10 DBG | debug message.
2022-09-12 23:49:51.553323 16180 purelog_demo/main.go:11 INF | info message.
2022-09-12 23:49:51.553323 16180 purelog_demo/main.go:12 WAR | warn message.
2022-09-12 23:49:51.553323 16180 purelog_demo/main.go:13 ERR | error message.
...
```



recommend:

```go
//new log config
config := purelog.NewConfig().
	SetStdout(true).                  //enable log to stdout
	SetCaller(true).                  //enable output caller
	SetFlush(100 * time.Millisecond)  //set flush interval

//new custom logger
logger := purelog.New(config)
//don't forget close the logger
defer logger.Close()

purelog.Debug("debug message.")
purelog.Info("info message.")
purelog.Warn("warn message.")
purelog.Error("error message.")

//set custom logger's level
config.SetLevel(purelog.LevelWarn)

//flush log data
purelog.Flush()
```



rotate:

```go
config := purelog.NewConfig().
	SetFile("test.log").  //basic name (more file name be test_Y-M-D_H-M-S_NS.log)
	SetSize(50 * 1024 * 1024).  //set single file size (50MB)
	SetCount(10).               //set max file count
	SetFlush(time.Second)       //set flush interval

//new logger
logger := purelog.New(config)

//don't forget close the logger
defer logger.Close()

//change single file size
config.SetSize(100 * 1024 * 1024) //100MB

//change max file count
config.SetCount(20)

//change file name
config.SetFile("test2.log")

logger.Info("enjoy yourself!")
```



### Licence

MIT Licence

Copyright (c) 2022 pure-project team.