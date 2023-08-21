package log

import (
	"context"
	"fmt"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/log/rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"time"
)

// 修改来自： https://github.com/gohp/logger

const (
	TimeDivision = "time"
	SizeDivision = "size"

	_defaultEncoding = "console"
	_defaultDivision = "size"
	_defaultUnit     = Hour
)

var (
	Logger                    *Log
	_encoderNameToConstructor = map[string]func(zapcore.EncoderConfig) zapcore.Encoder{
		"console": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewConsoleEncoder(encoderConfig)
		},
		"json": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewJSONEncoder(encoderConfig)
		},
	}
	MinLevel = zapcore.DebugLevel
	isInit   = false
)

type Log struct {
	L *zap.Logger
}

type Options struct {
	Encoding      string
	InfoFilename  string
	ErrorFilename string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	Division      string
	LevelSeparate bool
	TimeUnit      TimeUnit
	Stacktrace    bool
	EncodeTime    string
	closeDisplay  int
	caller        bool
}

// 初始化日志工具
func initLogger() {

	if isInit {
		return
	}

	conf := config.GetLogConfig()

	c := &Options{
		Division:      _defaultDivision,
		LevelSeparate: false,
		TimeUnit:      _defaultUnit,
		Encoding:      _defaultEncoding,
		caller:        false,
	}
	c.SetDivision("time") // 设置归档方式，"time"时间归档 "size" 文件大小归档
	c.SetTimeUnit(Day)    // 时间归档 可以设置切割单位
	c.SetEncoding("json") // 输出格式 "json" 或者 "console"

	if !conf.Console {
		c.CloseConsoleDisplay()
	}

	c.SetInfoFile(conf.LogFile) // 设置info级别日志文件
	if conf.LogErrorFile != "" {
		c.SetErrorFile(conf.LogErrorFile) // 设置warn级别日志文件
	}

	// 设置最低记录级别
	c.SetMinLevel(ParseLevel(conf.Level))

	c.InitLogger()

	isInit = true
}

func infoLevel() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		if lvl < MinLevel {
			return false
		}
		return lvl < zapcore.WarnLevel
	}
}

func warnLevel() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		if lvl < MinLevel {
			return false
		}
		return lvl >= zapcore.WarnLevel
	}
}

func (c *Options) SetDivision(division string) {
	c.Division = division
}

func (c *Options) SetEncodeTime(format string) {
	c.EncodeTime = format
}

func (c *Options) CloseConsoleDisplay() {
	c.closeDisplay = 1
}

func (c *Options) SetCaller(b bool) {
	c.caller = b
}

func (c *Options) SetTimeUnit(t TimeUnit) {
	c.TimeUnit = t
}

func (c *Options) SetErrorFile(path string) {
	c.LevelSeparate = true
	c.ErrorFilename = path
}

func (c *Options) SetInfoFile(path string) {
	c.InfoFilename = path
}

func (c *Options) SetEncoding(encoding string) {
	c.Encoding = encoding
}

func (c *Options) SetMinLevel(level zapcore.Level) {
	MinLevel = level
}

// isOutput whether set output file
func (c *Options) isOutput() bool {
	return c.InfoFilename != ""
}

func (c *Options) InitLogger() *Log {
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
		c.EncodeTime = RFC3339
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
		case TimeDivision:
			infoHook = c.timeDivisionWriter(c.InfoFilename)
			if c.LevelSeparate {
				warnHook = c.timeDivisionWriter(c.ErrorFilename)
			}
		case SizeDivision:
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

	Logger = &Log{logger}
	return Logger
}

func (c *Options) sizeDivisionWriter(filename string) io.Writer {
	hook := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxSize,
		Compress:   c.Compress,
	}
	return hook
}

func (c *Options) timeDivisionWriter(filename string) io.Writer {

	s := filename
	i := strings.LastIndex(s, ".")
	if i >= 0 {
		s = s[:i] + c.TimeUnit.Format() + "." + s[i+1:]
	}
	filename = s

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
	initLogger()
	Logger.L.Info(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	initLogger()
	Logger.L.Error(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	initLogger()
	Logger.L.Warn(msg, args...)
}

func Debug(msg string, args ...zap.Field) {
	initLogger()
	Logger.L.Debug(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	initLogger()
	Logger.L.Fatal(msg, args...)
}

func Infof(format string, args ...any) {
	initLogger()
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Info(logMsg)
}

func Errorf(format string, args ...any) {
	initLogger()
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Error(logMsg)
}

func Warnf(format string, args ...any) {
	initLogger()
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Warn(logMsg)
}

func Debugf(format string, args ...any) {
	initLogger()
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Debug(logMsg)
}

func Fatalf(format string, args ...any) {
	initLogger()
	logMsg := fmt.Sprintf(format, args...)
	Logger.L.Fatal(logMsg)
}

func With(k string, v any) zap.Field {
	return zap.Any(k, v)
}

func WithError(err error) zap.Field {
	return zap.NamedError("error", err)
}

func AddContext(ctx context.Context, fields ...zap.Field) context.Context {
	l := ctx.Value("_logger_ctx_val")
	logArgs, ok := l.([]zap.Field)
	if ok || logArgs == nil {
		logArgs = append(logArgs, fields...)
		ctx = context.WithValue(ctx, "_logger_ctx_val", logArgs)
	}
	return ctx
}

func withContext(ctx context.Context) *Log {
	if ctx == nil {
		return nil
	}
	l := ctx.Value("_logger_ctx_val")
	logArgs, _ := l.([]zap.Field)

	ctxLogger := &Log{
		L: Logger.L,
	}
	if len(logArgs) > 0 {
		ctxLogger.L = ctxLogger.L.With(logArgs...)
	}
	return ctxLogger
}

// Ctx new func
func Ctx(ctx context.Context) *Log {
	return withContext(ctx)
}

// ParseLevel 将字符串转换成枚举
func ParseLevel(text string) zapcore.Level {
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

func (l *Log) Info(msg string, args ...zap.Field) {
	l.L.Info(msg, args...)
}

func (l *Log) Error(msg string, args ...zap.Field) {
	l.L.Error(msg, args...)
}

func (l *Log) Warn(msg string, args ...zap.Field) {
	l.L.Warn(msg, args...)
}

func (l *Log) Debug(msg string, args ...zap.Field) {
	l.L.Debug(msg, args...)
}

func (l *Log) Fatal(msg string, args ...zap.Field) {
	l.L.Fatal(msg, args...)
}

func (l *Log) Infof(format string, args ...any) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Info(logMsg)
}

func (l *Log) Errorf(format string, args ...any) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Error(logMsg)
}

func (l *Log) Warnf(format string, args ...any) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Warn(logMsg)
}

func (l *Log) Debugf(format string, args ...any) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Debug(logMsg)
}

func (l *Log) Fatalf(format string, args ...any) {
	logMsg := fmt.Sprintf(format, args...)
	l.L.Fatal(logMsg)
}
