package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Debug().Msg("This message appears only when log level set to debug")
	log.Info().Msg("This message appears when log level set to debug or info")
}
func main() {
	log.Print("Hello World")
	// 调用Msg()或Send()之后，日志会被输出
	log.Debug().
		Str("Scale", "833 cents").
		Float64("Interval", 833.09).
		Msg("Fibonacci is everywhere")
	log.Debug().
		Str("Name", "Tom").
		Send()
	// 记录的字段可以任意嵌套，这通过Dict()来实现
	log.Info().
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Msg("hello world")
	// 可以调用Enabled()方法来判断日志是否需要输出，需要时再调用相应方法输出，节省了添加字段和日志信息的开销
	if e := log.Debug(); e.Enabled() {
		e.Str("foo", "bar").Msg("some debug message")
	}
	// 不想输出日志级别（即level字段），这时可以使用log.Log()方法
	log.Log().Str("foo", "bar").Msg("")
}
