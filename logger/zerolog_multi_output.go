package main

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func main() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout, getLumberjackWriter1("./med/go_logger_demo_zerolog.log"))
	l := zerolog.New(multi).With().Timestamp().Caller().Logger()
	l.Info().Msg("Hello World!")
}

func getLumberjackWriter1(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1,  // megabytes，超过则切割
		MaxBackups: 5,  // 最大文件保留数，超过就删除最老的日志文件
		MaxAge:     30, // 保存30天
		LocalTime:  true,
		Compress:   false, // 是否压缩（using gzip. The default is not to perform compression.）
	}
}
