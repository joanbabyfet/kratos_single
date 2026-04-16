package logger

import (
	"io"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/natefinch/lumberjack"
)

// NewLogger 提供全局日志
func NewLogger() log.Logger {

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    100, // MB
		MaxBackups: 10,
		MaxAge:     7,   // 天
		Compress:   true,
	}

	//同时输出：终端 + 文件
	writer := io.MultiWriter(os.Stdout, lumberjackLogger)

	logger := log.NewStdLogger(writer)

	return log.With(logger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)
}

var ProviderSet = wire.NewSet(NewLogger)