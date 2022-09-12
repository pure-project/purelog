# pure-log

**一个纯粹的go日志库。简单，快捷，够用。**





### 前言

pure-log是一个简单的golang异步日志库，生产可用。



### 功能

1. 异步日志
2. 动态配置
3. 大小数量控制滚动



### 用法

简单用法：

```go
//由于是异步日志，记得关闭
defer purelog.DefaultLogger.Close()

//与fmt.Print行为相同
purelog.Debug("debug message.")
purelog.Info("info message.")
purelog.Warn("warn message.")
purelog.Error("error message.")

//与fmt.Printf行为相同
purelog.Debugf("formatted message: %d or %v", 1234, 5678)

//设置默认日志等级
purelog.DefaultConfig.SetLevel(purelog.LevelWarn)

purelog.Info("this message can't be output!")
purelog.Warn("this message can be output.")

//设置日志输出文件
purelog.DefaultConfig.SetFile("test.log")

//主动将日志数据存盘(异步)
purelog.Flush()
```

输出：

```
//格式: 日期 时间.微秒 PID 文件:行号 等级 | 数据\n
2022-09-12 23:49:51.542323 16180 purelog_demo/main.go:10 DBG | debug message.
2022-09-12 23:49:51.553323 16180 purelog_demo/main.go:11 INF | info message.
2022-09-12 23:49:51.553323 16180 purelog_demo/main.go:12 WAR | warn message.
2022-09-12 23:49:51.553323 16180 purelog_demo/main.go:13 ERR | error message.
...
```



推荐用法：

```go
//创建日志配置
config := purelog.NewConfig().
	SetStdout(true).                  //enable log to stdout
	SetCaller(true).                  //enable output caller
	SetFlush(100 * time.Millisecond)  //set flush interval

//根据配置创建日志实例
logger := purelog.New(config)
//记得关闭日志实例
defer logger.Close()

purelog.Debug("debug message.")
purelog.Info("info message.")
purelog.Warn("warn message.")
purelog.Error("error message.")

//可动态设置日志配置
config.SetLevel(purelog.LevelWarn)

//将日志数据存盘(异步)
purelog.Flush()
```



滚动文件：

```go
config := purelog.NewConfig().
	SetFile("test.log").  //基本文件名 (自动产生文件名为 test_Y-M-D_H-M-S_NS.log)
	SetSize(50 * 1024 * 1024).  //单个日志文件大小 (50MB)
	SetCount(10).               //日志文件数量
	SetFlush(time.Second)       //存盘时间

//创建实例
logger := purelog.New(config)

//记得关闭实例
defer logger.Close()

//修改单文件大小
config.SetSize(100 * 1024 * 1024) //100MB

//修改文件数量
config.SetCount(20)

//修改文件名
config.SetFile("test2.log")

logger.Info("enjoy yourself!")
```





### 许可

MIT许可证

版权所有 (c) 2022 pure-project团队。