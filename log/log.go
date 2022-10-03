package log

import (
	"fmt"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"runtime"
	"time"
)

//logger := zap.NewExample()
//logger, _ = zap.NewDevelopment()
//logger, _ = zap.NewProduction()

var logger *zap.Logger

func newLogger() {

	logConf := config.LogConfig{}

	//获取日志配置信息
	err := config.ParseConfig("log", &logConf)
	if err != nil {
		panic(err)
	}

	level, _ := zap.ParseAtomicLevel(logConf.Level)

	//输出路径，路径可以是文件路径和stdout
	var outPut []string
	if logConf.Console {
		outPut = append(outPut, "stdout")
	}
	if logConf.LogFile != "" {

		paths, _ := filepath.Split(logConf.LogFile)

		//检查并创建目录
		err := utils.CheckAndMkdir(paths)
		if err != nil {
			fmt.Println("Err:" + err.Error())
			panic(err)
		}

		outPut = append(outPut, logConf.LogFile)
	}

	cfg := zap.Config{
		Level:       level,
		Encoding:    "json",
		OutputPaths: outPut,

		EncoderConfig: zapcore.EncoderConfig{
			LevelKey: "level",
			TimeKey:  "ts",
			NameKey:  "logger",
			//CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			//EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger = zap.Must(cfg.Build())
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getLogger() *zap.Logger {
	if logger == nil {
		newLogger()
	}
	return logger
}

// AddCommonField 添加公共字段
func AddCommonField(fields []zapcore.Field) []zapcore.Field {

	//调用栈信息
	pc, _, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	fields = append(fields, zap.String("func", utils.GetRelativePath(funcName)))
	fields = append(fields, zap.Int("line", line))

	return fields
}

func Info(msg string, fields ...zapcore.Field) {

	//添加公共字段
	fields = AddCommonField(fields)

	logger := getLogger()
	logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {

	//添加公共字段
	fields = AddCommonField(fields)

	logger := getLogger()
	logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {

	//添加公共字段
	fields = AddCommonField(fields)

	logger := getLogger()
	logger.Warn(msg, fields...)
}

func Error(err string, fields ...zapcore.Field) {

	//添加公共字段
	fields = AddCommonField(fields)

	logger := getLogger()
	logger.Error(err, fields...)
}
