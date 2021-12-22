package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"net/http"
	"os"
	"strings"
	"time"
)

var zLogger zerolog.Logger

type AddFieldHook struct {
}

func (AddFieldHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.DebugLevel {
		e.Str("hookField", "hook - msg: "+msg)
	}
}

func init() {
	// 打印堆栈信息
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	stdOut := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	stdOut.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %s |", i))
	}
	stdOut.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	stdOut.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	stdOut.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	//zLogger = zerolog.New(os.Stderr)
	//zLogger = zLogger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zLogger = log.Output(stdOut).With().Timestamp().Caller().Logger() // 根据全局log创建
	//zLogger = log.With().Timestamp().Caller().Logger()
	zLogger.Info().Str("foo111", "bar111").Msg("hello world111")
}

func main() {
	zLogger.Info().Str("foo1", "bar1").Msg("hello world1")
	//subLogger := zLogger.With().
	//	Str("foo", "bar").
	//	Logger()
	//subLogger.Info().Msg("hello world")
	// 日志采样功能
	sampled := zLogger.Sample(&zerolog.BasicSampler{N: 10})
	for i := 0; i < 20; i++ {
		sampled.Info().Msgf("will be logged every 10 message, idx: %d", i)
	}

	// 只采样Debug日志，在 1s 内最多输出 5 条日志，超过 5条 时，每隔 100条 输出一条
	sampled1 := log.Sample(&zerolog.LevelSampler{
		DebugSampler: &zerolog.BurstSampler{
			Burst:       5,
			Period:      time.Second,
			NextSampler: &zerolog.BasicSampler{N: 100},
		},
	})
	for i := 0; i < 1000; i++ {
		sampled1.Debug().Msgf("hello world, idx: %d", i)
	}

	hooked := zLogger.Hook(AddFieldHook{})
	hooked.Debug().Msg("Ha Ha Ha")

	simpleHttpGet2("www.google.com")
	simpleHttpGet2("http://www.baidu.com")
}

func simpleHttpGet2(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Str("url", url).Err(err).Msg("Error fetching url..")
	} else {
		log.Info().Str("statusCode", resp.Status).Str("url", url).Msg("Success..")
		resp.Body.Close()
	}
}
