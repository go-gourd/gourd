package log

import (
	"fmt"
	"github.com/go-gourd/gourd/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

//logger := zap.NewExample()
//logger, _ = zap.NewDevelopment()
//logger, _ = zap.NewProduction()

var logger *zap.Logger

func newLogger() {

	logConf := config.GetLogConfig()

	level, _ := zap.ParseAtomicLevel(logConf.Level)

	//输出路径，路径可以是文件路径和stdout
	var outPut []string

	//是否输出到控制台显示
	if logConf.Console {
		outPut = append(outPut, "stdout")
	}
	if logConf.LogFile != "" {

		paths, _ := filepath.Split(logConf.LogFile)

		//检查并创建目录
		err := checkAndMkdir(paths)
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

	//跳过调用栈记录
	if len(fields) == 1 && fields[0].Type == zapcore.SkipType {
		return fields
	}

	//调用栈信息
	pc, _, line, _ := runtime.Caller(2)
	caller := runtime.FuncForPC(pc).Name()

	fields = append(fields, zap.String("caller", getRelativePath(caller)))
	fields = append(fields, zap.Int("line", line))

	return fields
}

func Info(msg string, fields ...zapcore.Field) {
	fields = AddCommonField(fields)
	getLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	fields = AddCommonField(fields)
	getLogger().Debug(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	fields = AddCommonField(fields)
	getLogger().Warn(msg, fields...)
}

func Error(err string, fields ...zapcore.Field) {
	fields = AddCommonField(fields)
	getLogger().Error(err, fields...)
}

func checkAndMkdir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return nil
	}
	err = os.Mkdir(path, os.ModePerm)
	if err == nil {
		return err
	}
	return nil
}

func getRelativePath(path string) string {
	if filepath.IsAbs(path) {
		p, err := filepath.Rel(getRootPath(), path)
		if err != nil {
			return path
		}
		return p
	} else {
		return path
	}
}

func getRootPath() string {
	str, _ := os.Getwd()
	return str
}
