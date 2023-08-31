package log

import (
	"fmt"
	"github.com/go-gourd/gourd/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"time"
)

const (
	_defaultEncoding = "console"
	_defaultDivision = "time"
	_defaultUnit     = "day"
)

var (
	Logger                    *zap.Logger
	_encoderNameToConstructor = map[string]func(zapcore.EncoderConfig) zapcore.Encoder{
		"console": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewConsoleEncoder(encoderConfig)
		},
		"json": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewJSONEncoder(encoderConfig)
		},
	}
	minLevel = zapcore.DebugLevel
)

type options struct {
	Encoding      string
	InfoFilename  string
	ErrorFilename string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	Division      string
	LevelSeparate bool
	TimeUnit      timeUnit
	Stacktrace    bool
	EncodeTime    string
	closeDisplay  int
	caller        bool
}

// 初始化日志工具
func getLogger() *zap.Logger {

	if Logger != nil {
		return Logger
	}

	c := &options{
		Division:      _defaultDivision,
		LevelSeparate: false,
		TimeUnit:      _defaultUnit,
		Encoding:      _defaultEncoding,
		caller:        false,
	}

	conf := config.GetLogConfig()

	// 输出格式 "json" 或者 "console"
	if conf.Encoding != "" {
		c.Encoding = conf.Encoding
	}

	// 时间归档 可以设置切割单位
	if conf.TimeUnit != "" {
		c.TimeUnit = timeUnit(conf.TimeUnit)
	}

	// 设置归档方式，"time"时间归档 "size" 文件大小归档
	if conf.Division != "" {
		c.Division = conf.Division
	}

	// 设置是否开启控制台输出
	if !conf.Console {
		c.closeDisplay = 1
	}

	// 设置info级别日志文件
	if conf.LogFile != "" {
		c.InfoFilename = conf.LogFile
	}

	if conf.ErrorFile != "" {
		// 设置warn级别日志文件
		c.LevelSeparate = true
		c.ErrorFilename = conf.ErrorFile
	}

	// 设置最低记录级别
	if conf.Level == "" {
		conf.Level = "debug"
	}
	minLevel = parseLevel(conf.Level)

	c.InitLogger()

	return Logger
}

func infoLevel() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		if lvl < minLevel {
			return false
		}
		return lvl < zapcore.WarnLevel
	}
}

func warnLevel() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		if lvl < minLevel {
			return false
		}
		return lvl >= zapcore.WarnLevel
	}
}

// isOutput whether set output file
func (c *options) isOutput() bool {
	return c.InfoFilename != ""
}

func (c *options) InitLogger() *zap.Logger {
	var (
		logger             *zap.Logger
		infoHook, warnHook io.Writer
		wsInfo             []zapcore.WriteSyncer
		wsWarn             []zapcore.WriteSyncer
	)

	if c.Encoding == "" {
		c.Encoding = _defaultEncoding
	}
	if c.EncodeTime == "" {
		c.EncodeTime = "2006-01-02T15:04:05Z07:00"
	}
	encoder := _encoderNameToConstructor[c.Encoding]

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(c.EncodeTime),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	if c.closeDisplay == 0 {
		wsInfo = append(wsInfo, zapcore.AddSync(os.Stdout))
		wsWarn = append(wsWarn, zapcore.AddSync(os.Stdout))
	}

	// zap-core WriteSyncer setting
	if c.isOutput() {
		switch c.Division {
		case "time":
			infoHook = c.timeDivisionWriter(c.InfoFilename)
			if c.LevelSeparate {
				warnHook = c.timeDivisionWriter(c.ErrorFilename)
			}
		case "size":
			infoHook = c.sizeDivisionWriter(c.InfoFilename)
			if c.LevelSeparate {
				warnHook = c.sizeDivisionWriter(c.ErrorFilename)
			}
		}
		wsInfo = append(wsInfo, zapcore.AddSync(infoHook))
	}

	if c.ErrorFilename != "" {
		wsWarn = append(wsWarn, zapcore.AddSync(warnHook))
	}

	opts := make([]zap.Option, 0)
	cos := make([]zapcore.Core, 0)

	if c.LevelSeparate {
		cos = append(
			cos,
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsInfo...), infoLevel()),
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsWarn...), warnLevel()),
		)
	} else {
		cos = append(
			cos,
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsInfo...), zap.InfoLevel),
		)
	}

	opts = append(opts, zap.Development())

	if c.Stacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.WarnLevel))
	}

	if c.caller {
		opts = append(opts, zap.AddCaller())
	}

	logger = zap.New(zapcore.NewTee(cos...), opts...)

	Logger = logger
	return Logger
}

func (c *options) sizeDivisionWriter(filename string) io.Writer {
	hook := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxSize,
		Compress:   c.Compress,
	}
	return hook
}

func (c *options) timeDivisionWriter(filename string) io.Writer {

	s := filename
	i := strings.LastIndex(s, ".")
	if i >= 0 {
		s = s[:i] + c.TimeUnit.Format() + "." + s[i+1:]
	}
	filename = s

	// 用于日志切割归档
	hook, err := rotatelogs.New(
		filename,
		rotatelogs.WithMaxAge(time.Duration(int64(24*time.Hour)*int64(c.MaxAge))),
		rotatelogs.WithRotationTime(c.TimeUnit.RotationGap()),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func Skip() zap.Field {
	return zap.Skip()
}

func Info(msg string, args ...zap.Field) {
	getLogger().Info(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	getLogger().Error(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	getLogger().Warn(msg, args...)
}

func Debug(msg string, args ...zap.Field) {
	getLogger().Debug(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	getLogger().Fatal(msg, args...)
}

func Infof(format string, args ...any) {

	logMsg := fmt.Sprintf(format, args...)
	getLogger().Info(logMsg)
}

func Errorf(format string, args ...any) {

	logMsg := fmt.Sprintf(format, args...)
	getLogger().Error(logMsg)
}

func Warnf(format string, args ...any) {
	logMsg := fmt.Sprintf(format, args...)
	getLogger().Warn(logMsg)
}

func Debugf(format string, args ...any) {
	logMsg := fmt.Sprintf(format, args...)
	getLogger().Debug(logMsg)
}

func Fatalf(format string, args ...any) {
	logMsg := fmt.Sprintf(format, args...)
	getLogger().Fatal(logMsg)
}

// parseLevel 将字符串转换成枚举
func parseLevel(text string) zapcore.Level {
	switch text {
	case "debug", "DEBUG":
		return zapcore.DebugLevel
	case "info", "INFO", "": // make the zero value useful
		return zapcore.InfoLevel
	case "warn", "WARN":
		return zapcore.WarnLevel
	case "error", "ERROR":
		return zapcore.ErrorLevel
	case "fatal", "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
