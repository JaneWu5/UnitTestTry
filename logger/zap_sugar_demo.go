package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"time"
)

var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger1()
	defer sugarLogger.Sync()
	simpleHttpGet1("www.google.com")
	simpleHttpGet1("http://www.baidu.com")
}

func InitLogger1() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",                       //结构化（json）输出：msg的key
		LevelKey:     "level",                     //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "ts",                        //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file",                      //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, //输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
	}
	// core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	////自定义日志级别：自定义Info级别
	//infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl <= zapcore.InfoLevel
	//})
	////自定义日志级别：自定义Warn级别
	//warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl >= zapcore.WarnLevel
	//})

	// 获取io.Writer的实现
	infoWriter := getLumberjackWriter("./med/go_logger_demo_info.log")
	warnWriter := getLumberjackWriter("./med/go_logger_demo_warn.log")

	// 实现多个输出
	//core := zapcore.NewTee(
	//	zapcore.NewCore(getMyEncoder(), getLogWriterByAddSync(), zapcore.DebugLevel),
	//	// 将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
	//	zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), infoWriter, infoLevel),
	//	// warn及以上写入errPath
	//	zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), warnWriter, warnLevel),
	//	// 同时将日志输出到控制台，NewJSONEncoder 是结构化输出
	//	zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zapcore.InfoLevel),
	//)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(
		getLogWriterByAddSync(),
		zapcore.AddSync(os.Stdout),
		infoWriter, warnWriter,
	), zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getMyEncoder() zapcore.Encoder {
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriterByAddSync() zapcore.WriteSyncer {
	file, _ := os.Create("./med/go_logger_demo.log")
	return zapcore.AddSync(file)
}

// getLumberjackWriter Lumberjack assumes that only one process is writing to the output files.
// Using the same lumberjack configuration from multiple processes on the same machine will result in improper behavior.
func getLumberjackWriter(filename string) zapcore.WriteSyncer {
	l := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1,  // megabytes，超过则切割
		MaxBackups: 5,  // 最大文件保留数，超过就删除最老的日志文件
		MaxAge:     30, // 保存30天
		LocalTime:  true,
		Compress:   false, // 是否压缩（using gzip. The default is not to perform compression.）
	}
	return zapcore.AddSync(l)
}

// rotatelogs 这个项目已经不维护了
//func getRotateLogsWriter(filename string) io.Writer {
//	// 生成rotatelogs的Logger 实际生成的文件名 filename.YYmmddHH
//	// filename是指向最新日志的链接
//	hook, err := rotatelogs.New(
//		filename+".%Y%m%d%H",
//		rotatelogs.WithLinkName(filename),
//		rotatelogs.WithMaxAge(time.Hour*24*30),    // 保存30天
//		rotatelogs.WithRotationTime(time.Hour*24), // 切割频率 24小时
//	)
//	if err != nil {
//		sugarLogger.Error("日志启动异常")
//		panic(err)
//	}
//	return hook
//}

func simpleHttpGet1(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}
